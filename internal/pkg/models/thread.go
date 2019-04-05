package models

import (
	"github.com/hackallcode/db-homework/internal/pkg/verifier"
)

// easyjson:json
type Thread struct {
	Id      int64  `json:"id" example:"42"`
	Forum   string `json:"forum" example:"pirate-stories"`
	Author  string `json:"author" example:"j.sparrow"`
	Created string `json:"created" example:"2017-01-01T00:00:00.000Z"`
	Slug    string `json:"slug" example:"jones-cache"`
	Title   string `json:"title" example:"Davy Jones cache"`
	Message string `json:"message" example:"An urgent need to reveal the hiding place of Davy Jones. Who is willing to help in this matter?"`
	Votes   int64  `json:"votes" example:"1000"`
}

func (data Thread) Validate() bool {
	if !verifier.IsNickname(data.Author, false) {
		return false
	}
	return true
}

// easyjson:json
type ThreadUpdate struct {
	SlugOrId string `json:"-"`
	Message  string `json:"message" example:"An urgent need to reveal the hiding place of Davy Jones. Who is willing to help in this matter?"`
	Title    string `json:"title" example:"Davy Jones cache"`
}

func (data ThreadUpdate) Validate() bool {
	return true
}

// easyjson:json
type Vote struct {
	SlugOrId string `json:"-"`
	Nickname string `json:"nickname" example:"j.sparrow"`
	Vote     int64  `json:"voice" example:"-1"`
}

func (data Vote) Validate() bool {
	return true
}

// easyjson:json
type Threads []Thread
