package questiondb

import (
	"collect/business/core/question"
	"time"
)

// Question represents an individual ques row in db.
type dbQues struct {
	ID           int64     `db:"question_id"`
	FormID       int64     `db:"form_id"`
	QuestionType string    `db:"question_type"`
	QuestionText string    `db:"question_text"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

// NewQuestion is what we require from clients when adding a Question.
type NewQuestion struct {
	UserID int64
}

func toDBQuestion(que question.Question) dbQues {
	queDB := dbQues{
		ID:           que.ID,
		FormID:       que.FormID,
		QuestionType: que.QuestionType,
		QuestionText: que.QuestionText,
		CreatedAt:    que.CreatedAt.UTC(),
		UpdatedAt:    que.UpdatedAt.UTC(),
	}

	return queDB
}

func toCoreQuestion(db dbQues) question.Question {
	prdDB := question.Question{
		ID:           db.ID,
		FormID:       db.FormID,
		QuestionType: db.QuestionType,
		QuestionText: db.QuestionText,
		CreatedAt:    db.CreatedAt.UTC(),
		UpdatedAt:    db.UpdatedAt.UTC(),
	}

	return prdDB
}

func toCoreQuestionSlice(dbQuestions []dbQues) []question.Question {
	frms := make([]question.Question, len(dbQuestions))
	for i, dbQuest := range dbQuestions {
		frms[i] = toCoreQuestion(dbQuest)
	}
	return frms
}
