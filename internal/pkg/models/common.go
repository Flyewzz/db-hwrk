package models

import (
	"errors"
)

/********************/
/*      ERRORS      */
/********************/

var (
	NotFoundError      = errors.New("not found")
	AlreadyExistsError = errors.New("already exists")
	ConflictDataError = errors.New("conflict data")
)

/********************/
/*    API MODELS    */
/********************/

type InputModel interface {
	Validate() bool
}

type OutputModel interface {
	MarshalJSON() ([]byte, error)
}

/* MESSAGE */

// easyjson:json
type MessageAnswer struct {
	Message string `json:"message" example:"All right!"`
}

func NewMessageAnswer(message string) MessageAnswer {
	return MessageAnswer{Message: message}
}

/* ERROR */

// easyjson:json
type ErrorAnswer struct {
	Message string `json:"message" example:"Can't find user with id #42"`
}

func NewErrorAnswer(message string) ErrorAnswer {
	return ErrorAnswer{Message: message}
}

var (
	IncorrectJsonAnswer = NewErrorAnswer("Incorrect JSON!")
	IncorrectUrlAnswer  = NewErrorAnswer("Incorrect URL!")
	IncorrectDataAnswer = NewErrorAnswer("Incorrect data!")
	DuplicateDataAnswer = NewErrorAnswer("Duplicate data!")
	UserNotFoundAnswer  = NewErrorAnswer("Author not found!")
	ForumNotFoundAnswer  = NewErrorAnswer("Forum not found!")
	ThreadNotFoundAnswer  = NewErrorAnswer("Thread not found!")
	ParentIncorrectAnswer  = NewErrorAnswer("Parent not found!")
)

