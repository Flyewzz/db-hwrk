package models

import (
	"github.com/hackallcode/db-homework/internal/pkg/verifier"
)

// easyjson:json
type Post struct {
	Id       int64  `json:"id" example:"314"`
	Forum    string `json:"forum" example:"pirate-stories"`
	Thread   string `json:"-" example:"pirate-stories"`
	ThreadId int64  `json:"thread" example:"311"`
	Author   string `json:"author" example:"j.sparrow"`
	Created  string `json:"created" example:"2017-01-01T00:00:00.000Z"`
	Message  string `json:"message" example:"We should be afraid of the Kraken."`
	IsEdited bool   `json:"isEdited" example:"true"`
	Parent   int64  `json:"parent" example:"0"`
}

func (data Post) Validate() bool {
	if !verifier.IsNickname(data.Author, false) {
		return false
	}
	return true
}

// easyjson:json
type PostUpdate struct {
	Id      int64  `json:"-"`
	Message string `json:"message" example:"We should be afraid of the Kraken."`
}

func (data PostUpdate) Validate() bool {
	return true
}

// easyjson:json
type Posts []Post

// easyjson:json
type PostFull struct {
	Author *User   `json:"author"`
	Forum  *Forum  `json:"forum"`
	Post   *Post   `json:"post"`
	Thread *Thread `json:"thread"`
}

func (data PostFull) Validate() bool {
	return true
}
