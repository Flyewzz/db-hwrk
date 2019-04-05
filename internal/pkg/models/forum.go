package models

import (
	"github.com/hackallcode/db-homework/internal/pkg/verifier"
)

// easyjson:json
type Forum struct {
	Id      int64  `json:"-"`
	Slug    string `json:"slug" example:"pirate-stories"`
	Title   string `json:"title" example:"Pirate stories"`
	User    string `json:"user" example:"j.sparrow"`
	Threads int64  `json:"threads" example:"200"`
	Posts   int64  `json:"posts" example:"200000"`
}

func (data Forum) Validate() bool {
	if !verifier.IsTitle(data.Title, false) {
		return false
	}
	if !verifier.IsSlug(data.Slug, false) {
		return false
	}
	if !verifier.IsNickname(data.User, false) {
		return false
	}
	return true
}
