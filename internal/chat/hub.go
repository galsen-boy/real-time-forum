package chat

import (
	"encoding/json"
	"fmt"

	"real-time-forum/internal/structure"
)

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {
	clients      map[int]*Client // Registered clients
	broadcast    chan []byte     // Inbound messages from the clients
	register     chan *Client    // Register requests from the clients
	unregister   chan *Client    // Unregister requests from clients
	typing       map[int]bool    // Map to store the typing status of clients
	typing2      map[string]bool // Map to store the typing status of clients
	typingStatus map[int]int     // Map to store the typing s
}

func NewHub() *Hub {
	return &Hub{
		broadcast:    make(chan []byte),     // Initialize the broadcast channel
		register:     make(chan *Client),    // Initialize the register channel
		unregister:   make(chan *Client),    // Initialize the unregister channel
		clients:      make(map[int]*Client), // Initialize the clients map
		typing:       make(map[int]bool),    // Initialize the typing map
		typing2:      make(map[string]bool), // Initialize the typing map
		typingStatus: make(map[int]int),     // Initialize the typing status map
	}
}

func (h *Hub) Run() { // Run the hub
	for {
		select {
		case client := <-h.register: // Register a client
			h.clients[client.userID] = client // Add the client to the clients map

			// Notify other clients that this client is online
			uids := make([]int, 0, len(h.clients)) // Create a slice of user IDs
			for id := range h.clients {            // Iterate over the clients map
				uids = append(uids, id) // Append the user ID to the slice
			}
			msg := structure.OnlineUsers{ // Create a message
				UserIds:  uids,     // Set the user IDs
				Msg_type: "online", // Set the message type
			}
			sendMsg, err := json.Marshal(msg)
			if err != nil {
				panic(err)
			}

			for _, c := range h.clients {
				select {
				case c.send <- sendMsg: // Send the message to the client
				default:
					close(c.send)               // Close the send channel
					delete(h.clients, c.userID) // Delete the client from the clients map
				}
			}
		case client := <-h.unregister: // Unregister a client
			if _, ok := h.clients[client.userID]; ok { // Check if the client is registered
				delete(h.clients, client.userID)

				// Notify other clients that this client is offline
				uids := make([]int, 0, len(h.clients))
				for id := range h.clients {
					uids = append(uids, id)
				}
				msg := structure.OnlineUsers{
					UserIds:  uids,
					Msg_type: "online",
				}
				sendMsg, err := json.Marshal(msg)
				if err != nil {
					panic(err)
				}

				for _, c := range h.clients {
					select {
					case c.send <- sendMsg:
					default:
						close(c.send)
						delete(h.clients, c.userID)
					}
				}

				close(client.send)
			}
		case message := <-h.broadcast:
			// Process the message
			var msg structure.Message
			if err := json.Unmarshal(message, &msg); err != nil {
				panic(err)
			}

			sendMsg, err := json.Marshal(msg)
			if err != nil {
				panic(err)
			}

			if msg.Msg_type == "msg" { // Check if the message is a chat message
				for _, client := range h.clients {
					if client.userID == msg.Receiver_id {
						select {
						case client.send <- sendMsg:
						default:
							close(client.send)
							delete(h.clients, client.userID)
						}
					}
				}
			} else { // Check if the message is a typing status update
				for _, client := range h.clients { // Iterate over the clients map
					if client.userID != msg.Sender_id { // Check if the client is not the sender
						select {
						case client.send <- sendMsg: // Send the message to the client
						default:
							close(client.send)
							delete(h.clients, client.userID)
						}
					}
				}
			}
		}
	}
}

// UpdateTypingStatus updates the typing status of a client in the hub.
func (h *Hub) UpdateTypingStatus(senderID, receiverID int, isTyping bool) {
	h.typing2[fmt.Sprintf("%d_%d", senderID, receiverID)] = isTyping // Update the typing status of the client

	// Notify the receiver client about the typing status change
	msg := structure.TypingStatus{
		UserID:     senderID,
		IsTyping:   isTyping,
		Msgtype:    "typing",
		ReceiverID: receiverID,
		SenderID:   senderID,
	}
	sendMsg, err := json.Marshal(msg) // Create a message
	if err != nil {
		panic(err)
	}

	for _, client := range h.clients {
		if client.userID == receiverID {
			select {
			case client.send <- sendMsg:
			default:
				close(client.send)
				delete(h.clients, client.userID)
			}
			break // Exit the loop after sending the message to the receiver client
		}
	}
}
