package answerdb

import (
	"collect/business/core/answer"
	"time"
)

// dbAnswer represents an individual answer mapping
type dbAnswer struct {
	QuestionID     int64     `db:"question_id"`
	ResponseID     int64     `db:"response_id"`
	AnswerText     string    `db:"answer_text"`
	AnswerOptionID int       `db:"answer_option_id"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
	FormID         int64     `db:"-"`
}

func toDBAnswer(prd answer.Answer) dbAnswer {
	prdDB := dbAnswer{
		ResponseID:     prd.ResponseID,
		QuestionID:     prd.QuestionID,
		AnswerText:     *prd.AnswerText,
		AnswerOptionID: *prd.AnswerOptionID,
		CreatedAt:      prd.CreatedAt.UTC(),
		UpdatedAt:      prd.UpdatedAt.UTC(),
	}

	return prdDB
}

func toCoreAnswer(dbItm dbAnswer) answer.Answer {
	prd := answer.Answer{
		ResponseID:     dbItm.ResponseID,
		QuestionID:     dbItm.QuestionID,
		AnswerText:     &dbItm.AnswerText,
		AnswerOptionID: &dbItm.AnswerOptionID,
		CreatedAt:      dbItm.CreatedAt.In(time.Local),
		UpdatedAt:      dbItm.UpdatedAt.In(time.Local),
	}

	return prd
}

func toCoreAnswerSlice(dbItems []dbAnswer) []answer.Answer {
	prds := make([]answer.Answer, len(dbItems))
	for i, dbPrd := range dbItems {
		prds[i] = toCoreAnswer(dbPrd)
	}
	return prds
}
