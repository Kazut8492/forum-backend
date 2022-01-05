package data

// type Session struct {
// 	userId     int
// 	username   string
// 	sessionId  int
// 	expiration time.Time
// }

type User struct {
	User_ID  int
	Username string
	Pass     string
}

type Comment struct {
	Id      int
	PostId  int
	Title   string
	Content string
}

type Post struct {
	Id       int
	Title    string
	Content  string
	Comments []Comment
}
