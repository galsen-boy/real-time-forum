package structure

type Post struct {
	Id       int    `json:"id"`
	User_id  int    `json:"user_id"`
	Category string `json:"category"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Date     string `json:"date"`
	Likes    int    `json:"likes"`
	Dislikes int    `json:"dislikes"`
}

type Comment struct {
	Id      int    `json:"id"`
	Post_id int    `json:"post_id"`
	User_id int    `json:"user_id"`
	Content string `json:"content"`
	Date    string `json:"date"`
}

type User struct {
	Id        int    `json:"id"`
	Username  string `json:"username"`
	Firstname string `json:"firstname"`
	Surname   string `json:"surname"`
	Gender    string `json:"gender"`
	Email     string `json:"email"`
	DOB       string `json:"dob"`
	Password  string `json:"password"`
}

type Message struct {
	Id          int    `json:"id"`
	Sender_id   int    `json:"sender_id"`
	Receiver_id int    `json:"receiver_id"`
	Content     string `json:"content"`
	Date        string `json:"date"`
	Msg_type    string `json:"msg_type"`
	UserID      int    `json:"user_id"`
	IsTyping    bool   `json:"is_typing"`
	ImageData   string `json:"image_data"`
}

type Login struct {
	Data     string `json:"emailUsername"`
	Password string `json:"password"`
}

type Chat struct {
	User_one int
	User_two int
	Time     int
}

type OnlineUsers struct {
	UserIds  []int  `json:"user_ids"`
	Msg_type string `json:"msg_type"`
}

type Resp struct {
	Msg string `json:"msg"`
}

type Session struct {
	Session_uuid string
	User_id      int
}

// typing status
type TypingStatus struct {
	UserID     int    `json:"user_id"`
	IsTyping   bool   `json:"is_typing"`
	Msgtype    string `json:"msg_type"`
	ReceiverID int    `json:"receiver_id"`
	SenderID   int    `json:"sender_id"`
}

// ...

// Check for typing status updates
var data struct {
	Typing bool `json:"typing"`
}
