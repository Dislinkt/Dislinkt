package domain

type Post struct {
	UserId     string
	PostText   string
	ImagePaths []string
	DatePosted string
	Reactions  []Reaction
	Comments   []Comment
	Links      Links
}

type Comment struct {
	Username    string
	CommentText string
}

type Reaction struct {
	Username string
	Reaction int
}

type Links struct {
	Comment string
	Like    string
	Dislike string
}

type ConnectionRequest struct {
	UserId    string
	Name      string
	Surname   string
	Biography string
}
