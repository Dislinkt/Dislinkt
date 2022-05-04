package domain

type User struct {
	Id       string `gorm:"index:idx_name,unique"`
	Username string `gorm:"unique"`
	Password string
}

type LoginRequest struct {
	Username string `bson:"password"`
	Password string `bson:"password"`
}
