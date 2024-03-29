package domain

type Post struct {
	Id             string
	UserId         string
	PostText       string
	ImagePaths     []string
	DatePosted     string
	LikesNumber    int
	DislikesNumber int
	CommentsNumber int
	Links          Links
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

type ConnectionUser struct {
	UserId    string
	Name      string
	Surname   string
	Biography string
	Username  string
}
