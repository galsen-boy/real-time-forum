package database

import (
	"database/sql"
	"errors"
	"strconv"

	"real-time-forum/internal/structure"
)

// Attempts to insert a new message into the database
func NewMessage(path string, m structure.Message) error {
	//Opens the database
	db, err := OpenDB(path)
	if err != nil {
		return err
	}

	defer db.Close()

	//Executes the insert statement
	_, err = db.Exec(AddMessage, m.Sender_id, m.Receiver_id, m.Content, m.Date)
	if err != nil {
		return err
	}

	err = UpdateChatTime(m.Sender_id, m.Receiver_id, db)
	if err != nil {
		return err
	}

	return nil
}

// Converts message table query results into an array of message structs
func ConvertRowToMessage(rows *sql.Rows) ([]structure.Message, error) {
	var messages []structure.Message

	//Loops through the rows provided
	for rows.Next() {
		var m structure.Message

		//Stores the row data in a temporary message struct
		err := rows.Scan(&m.Id, &m.Sender_id, &m.Receiver_id, &m.Content, &m.Date)
		if err != nil {
			break
		}

		//Appends the temporary struct to the array
		messages = append(messages, m)
	}

	return messages, nil
}

// Finds chat messages between users
func FindChatMessages(path, sender, receiver string, firstId int) ([]structure.Message, error) {
	//Opens the database
	db, err := OpenDB(path)
	if err != nil {
		return []structure.Message{}, errors.New("failed to open database")
	}

	defer db.Close()

	//Converts sender and receiver ids to integers
	s, err := strconv.Atoi(sender)
	if err != nil {
		return []structure.Message{}, errors.New("sender id must be an integer")
	}

	r, err := strconv.Atoi(receiver)
	if err != nil {
		return []structure.Message{}, errors.New("receiver id must be an integer")
	}

	//Searches database for all messages between the two users
	q, err := db.Query(GetAllChatMessage, s, r, r, s, firstId)
	//`SELECT * FROM messages WHERE (sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?) AND  id <= ?  ORDER BY id DESC LIMIT 10`
	if err != nil {
		return []structure.Message{}, errors.New("could not find chat messages")
	}

	//search database for last message between the two users
	// q, err := db.Query(GetLastMessage, r, r, r, s)
	// //`SELECT * FROM messages WHERE (sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?) ORDER BY id DESC LIMIT 1`
	// if err != nil {
	// 	return []structure.Message{}, errors.New("could not find chat messages")
	// }

	//Converts rows to an array of message structs
	messages, err := ConvertRowToMessage(q)
	if err != nil {
		return []structure.Message{}, errors.New("failed to convert")
	}

	return messages, nil
}

// find the last message between two users
func FindLastMessage(path, sender, receiver string) (structure.Message, error) {
	//Opens the database
	db, err := OpenDB(path)
	if err != nil {
		return structure.Message{}, errors.New("failed to open database")
	}

	defer db.Close()

	//Converts sender and receiver ids to integers
	s, err := strconv.Atoi(sender)
	if err != nil {
		return structure.Message{}, errors.New("sender id must be an integer")
	}

	r, err := strconv.Atoi(receiver)
	if err != nil {
		return structure.Message{}, errors.New("receiver id must be an integer")
	}

	//search database for last message between the two users
	q, err := db.Query(GetLastMessage, r, r, r, s)
	//`SELECT * FROM messages WHERE (sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?) ORDER BY id DESC LIMIT 1`
	if err != nil {
		return structure.Message{}, errors.New("could not find chat messages")
	}

	//Converts rows to an array of message structs
	messages, err := ConvertRowToMessage(q)
	if err != nil {
		return structure.Message{}, errors.New("failed to convert")
	}

	return messages[0], nil
}
