package models

import (
	"github.com/hackallcode/db-homework/internal/pkg/verifier"
)

/* USER */

// easyjson:json
type User struct {
	Id       int64  `json:"-"`
	Nickname string `json:"nickname" example:"j.sparrow"`
	Email    string `json:"email" example:"captaina@blackpearl.sea"`
	FullName string `json:"fullname" example:"Captain Jack Sparrow"`
	About    string `json:"about" example:"This is the day you will always remember as the day that you almost caught Captain Jack Sparrow!"`
}

func (data User) Validate() bool {
	if !verifier.IsNickname(data.Nickname, false) {
		return false
	}
	if !verifier.IsEmail(data.Email, false) {
		return false
	}
	if !verifier.IsFullName(data.FullName, false) {
		return false
	}
	return true
}

/* USER UPDATE */

// easyjson:json
type UserUpdate struct {
	Nickname string `json:"-"`
	Email    string `json:"email" example:"captaina@blackpearl.sea"`
	FullName string `json:"fullname" example:"Captain Jack Sparrow"`
	About    string `json:"about" example:"This is the day you will always remember as the day that you almost caught Captain Jack Sparrow!"`
}

func (data UserUpdate) Validate() bool {
	if !verifier.IsNickname(data.Nickname, false) {
		return false
	}
	if !verifier.IsEmail(data.Email, true) {
		return false
	}
	if !verifier.IsFullName(data.FullName, true) {
		return false
	}
	return true
}

/* USERS */

// easyjson:json
type Users []User
