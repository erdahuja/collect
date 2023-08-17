package sheetsmock

import (
	"fmt"
	"time"
)

// Response represents an individual form in the response.
type Response struct {
	ResponseID   int64     `json:"response_id,omitempty"`
	FormID       int64     `json:"form_id,omitempty"`
	RespondentID string    `json:"respondent_id,omitempty"` // some unique identifier so that one respondent can have one response in one form, this can be some id also
	CreatedAt    time.Time `json:"created_at,omitempty"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
}

type Form struct {
	ID          int64     `json:"form_id"`
	Title       string    `json:"form_title"`
	Description string    `json:"form_description"`
	DateCreated time.Time `json:"created_at"`
	DateUpdated time.Time `json:"updated_at"`
}

type AnswerEvent struct {
	Response Response
	Answer   Answer
}

// Answer represents an individual answer
type Answer struct {
	AnswerID       int64
	QuestionID     int64
	ResponseID     int64
	AnswerText     *string
	AnswerOptionID *int
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// Question represents an individual user in db
type Question struct {
	ID           int64     `json:"question_id"`
	FormID       int64     `json:"form_id"`
	QuestionType string    `json:"question_type"`
	QuestionText string    `json:"question_text"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
}

func (q Question) Column() string {
	return fmt.Sprintf("%d:{%s}:%s", q.ID, q.QuestionText, q.QuestionType)
}

func (q Question) QueID() string {
	return fmt.Sprintf("%d:{%s}:%s", q.ID, q.QuestionText, q.QuestionType)
}

