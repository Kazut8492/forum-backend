package src

// type Session struct {
// 	userId     int
// 	username   string
// 	sessionId  int
// 	expiration time.Time
// }

type User struct {
	ID       int
	Username string
	Email    string
	Pass     string
}

type Comment struct {
	ID      int
	PostId  int
	Title   string
	Content string
}

type Post struct {
	ID       int
	Title    string
	Content  string
	Comments []Comment
}
