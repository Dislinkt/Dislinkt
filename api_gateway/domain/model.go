package domain

type Post struct {
	UserId     string
	PostText   string
	ImagePaths []string
	DatePosted string
	Reactions  []Reaction
	Comments   []Comment
}

type Comment struct {
	Username    string
	CommentText string
}

type Reaction struct {
	Username string
	Reaction int
}
