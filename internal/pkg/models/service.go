package models

// easyjson:json
type Status struct {
	Forum  int64 `json:"forum" example:"100"`
	Post   int64 `json:"post" example:"1000000"`
	Thread int64 `json:"thread" example:"1000"`
	User   int64 `json:"user" example:"1000"`
}

