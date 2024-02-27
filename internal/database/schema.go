package database

//SQL statement to initialise the database tables
const (
	CreateTables = `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username VARCHAR(64) NULL UNIQUE,
		firstname VARCHAR(64) NOT NULL,
		surname VARCHAR(64) NOT NULL,
		gender VARCHAR(64) NULL,
		email VARCHAR(64) NOT NULL UNIQUE,
		dob TEXT NULL,
		password VARCHAR(255) NOT NULL
	);

	CREATE TABLE IF NOT EXISTS sessions (
		session_uuid VARCHAR(255) NOT NULL UNIQUE,
		user_id INTEGER NOT NULL,
		FOREIGN KEY(user_id) REFERENCES users(id)
	);

	CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		category VARCHAR(64),
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		date TEXT NOT NULL,
		likes INTEGER NOT NULL,
		dislikes INTEGER NOT NULL,
		FOREIGN KEY(user_id) REFERENCES users(id)
	);

	CREATE TABLE IF NOT EXISTS comments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		post_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		content TEXT NOT NULL,
		date TEXT NOT NULL,
		FOREIGN KEY(post_id) REFERENCES posts(id),
		FOREIGN KEY(user_id) REFERENCES users(id)
	);

	CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		sender_id INTEGER NOT NULL,
		receiver_id INTEGER NOT NULL,
		content TEXT NOT NULL,
		date TEXT NOT NULL,
		FOREIGN KEY(sender_id) REFERENCES users(id),
		FOREIGN KEY(receiver_id) REFERENCES users(id)
	);

	CREATE TABLE IF NOT EXISTS chats (
		id_one INTEGER NOT NULL,
		id_two INTEGER NOT NULL,
		time INTEGER NOT NULL,
		FOREIGN KEY(id_one) REFERENCES users(id),
		FOREIGN KEY(id_two) REFERENCES users(id)
	);

	CREATE TABLE IF NOT EXISTS liked_posts (
		post_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		FOREIGN KEY(post_id) REFERENCES posts(id),
		FOREIGN KEY(user_id) REFERENCES users(id)
	);

	CREATE TABLE IF NOT EXISTS disliked_posts (
		post_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		FOREIGN KEY(post_id) REFERENCES posts(id),
		FOREIGN KEY(user_id) REFERENCES users(id)
	);
	`
)