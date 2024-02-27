package database

// Insert statements to add data to the database
const (
	AddUser    = `INSERT INTO users(username, firstname, surname, gender, email, dob, password) values(?, ?, ?, ?, ?, ?, ?)`
	AddPost    = `INSERT INTO posts(user_id, category, title, content, date, likes, dislikes) values(?, ?, ?, ?, ?, 0, 0)`
	AddComment = `INSERT INTO comments(post_id, user_id, content, date) values(?, ?, ?, ?)`
	AddMessage = `INSERT INTO messages(sender_id, receiver_id, content, date) values(?, ?, ?, ?)`
	AddLike    = `INSERT INTO liked_posts(post_id, user_id) values(?, ?)`
	AddDislike = `INSERT INTO disliked_posts(post_id, user_id) values(?, ?)`
	AddSession = `INSERT INTO sessions(session_uuid, user_id) values(?, ?)`
	AddChat    = `INSERT INTO chats(id_one, id_two, time) values(? ,?, ?)`
)

// Query statements to filter data from the database
const (
	GetUserById          = `SELECT * FROM users WHERE id = ?`
	GetUserByUsername    = `SELECT * FROM users WHERE username = ?`
	GetUserByEmail       = `SELECT * FROM users WHERE email = ?`
	GetAllUser           = `SELECT * FROM users ORDER BY username ASC`
	GetPostById          = `SELECT * FROM posts WHERE id = ? ORDER BY id DESC`
	GetAllPost           = `SELECT * FROM posts ORDER BY id DESC`
	GetAllPostByCategory = `SELECT * FROM posts WHERE category = ? ORDER BY id DESC`
	GetAllPostByUser     = `SELECT * FROM posts WHERE user_id = ? ORDER BY id DESC`
	GetCommentById       = `SELECT * FROM comments WHERE id = ?`
	GetAllPostComment    = `SELECT * FROM comments WHERE post_id = ?`
	GetAllUserComment    = `SELECT * FROM comments WHERE user_id = ?`
	GetMessage           = `SELECT * FROM messages WHERE id = ?`
	GetAllChatMessage    = `SELECT * FROM messages WHERE ((sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)) AND ( id <= ? ) ORDER BY id DESC LIMIT 10`
	GetLastMessage       = `SELECT * FROM messages WHERE ((sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)) ORDER BY id DESC LIMIT 1`
	GetPostLikes         = `SELECT users.* FROM liked_posts INNER JOIN users ON liked_posts.user_id = users.id WHERE liked_posts.post_id = ?`
	GetUserLikes         = `SELECT posts.* FROM liked_posts INNER JOIN posts ON liked_posts.post_id = posts.id WHERE liked_posts.user_id = ? ORDER BY id DESC`
	GetPostDislikes      = `SELECT users.* FROM disliked_posts INNER JOIN users ON disliked_posts.user_id = users.id WHERE disliked_posts.post_id = ?`
	GetUserDislikes      = `SELECT posts.* FROM disliked_posts INNER JOIN posts ON disliked_posts.post_id = posts.id WHERE disliked_posts.user_id = ? ORDER BY id DESC`
	GetSessionUser       = `SELECT users.* FROM sessions INNER JOIN users ON sessions.user_id = users.id WHERE sessions.session_uuid = ?`
	GetUserChats         = `SELECT * FROM chats WHERE id_one = ? OR id_two = ? ORDER BY time DESC`
	GetChatBetween       = `SELECT * FROM chats WHERE id_one = ? AND id_two = ? OR id_one = ? AND id_two = ?`
)

// Query statements to remove data from database
const (
	RemoveCookie  = `DELETE FROM sessions WHERE user_id = ?`
	RemoveLike    = `DELETE FROM liked_posts WHERE post_id = ? AND user_id = ?`
	RemoveDislike = `DELETE FROM disliked_posts WHERE post_id = ? AND user_id = ?`
)

// Query statements to update data in database
const (
	UpdateLike    = `UPDATE posts SET likes = ? WHERE id = ?`
	UpdateDislike = `UPDATE posts SET dislikes = ? WHERE id = ?`
	UpdateChat    = `UPDATE chats SET time = ? WHERE id_one = ? AND id_two = ?`
)
