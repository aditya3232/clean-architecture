package entity

type PublishMessage struct {
	UserId int64  `json:"user_id"`
	Email  string `json:"email"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}
