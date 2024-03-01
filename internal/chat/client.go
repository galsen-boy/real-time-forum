package chat

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"

	"real-time-forum/internal/config"
	"real-time-forum/internal/database"
	"real-time-forum/internal/structure"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub        *Hub
	conn       *websocket.Conn // The websocket connection.
	send       chan []byte     // Buffered channel of outbound messages.
	userID     int             // The user id of the client
	typing     bool            // Track the typing status
	typingLock chan bool
	isReceiver bool // Track if the client is the receiver
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		log.Printf("Received message from client: %s", string(message))

		var msg structure.Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			break
		}

		msg.Sender_id = c.userID

		if msg.Msg_type == "msg" {
			msg.Date = time.Now().Format("01-02-2006 15:04:05")

			err = database.NewMessage(config.Path, msg)
			if err != nil {
				log.Printf("Error storing new message: %v", err)
				break
			}

		} else if msg.Msg_type == "typing" {
			c.hub.UpdateTypingStatus(c.userID, msg.Receiver_id, msg.IsTyping)
		}

		sendMsg, err := json.Marshal(msg)
		if err != nil {
			log.Printf("Error marshaling message: %v", err)
			break
		}

		c.hub.broadcast <- sendMsg
	}

	// Stop typing when the readPump exits
	c.typingLock <- false
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
// writePump pumps messages from the hub to the websocket connection.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		case isTyping := <-c.typingLock:
			c.typing = isTyping

			// Send typing status only if the client is the receiver
			if c.isReceiver {
				typingStatus := struct {
					ReceiverID int    `json:"receiver_id"`
					IsTyping   bool   `json:"is_typing"`
					MsgType    string `json:"msg_type"`
				}{
					ReceiverID: c.userID,
					IsTyping:   c.typing,
					MsgType:    "typing",
				}
				sendTypingStatus, err := json.Marshal(typingStatus)
				if err != nil {
					log.Println("Error marshaling typing status:", err)
					continue
				}
				c.hub.broadcast <- sendTypingStatus
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	cookie, err := r.Cookie("session")
	if err != nil {
		return
	}

	foundVal := cookie.Value

	curr, err := database.CurrentUser(config.Path, foundVal)
	if err != nil {
		return
	}

	client := &Client{
		hub:        hub,
		conn:       conn,
		send:       make(chan []byte, 256),
		userID:     curr.Id,
		typing:     false,
		typingLock: make(chan bool),
		isReceiver: false, // true or false based on your logic to determine if the client is the receiver,
	}

	log.Println("Client isReceiver:", client.isReceiver)
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()

	// Start a goroutine to listen for typing events
	go func() {
		for{
			select {
			case <-client.typingLock:
				if client.typing {
					log.Println("User started typing")
				} else {
					log.Println("User stopped typing")
				}
			}
		}
	}()
}
