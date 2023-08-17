package responsedb

import (
	"collect/business/core/response"
	"time"
)

// dbResponse represents an individual response mapping
type dbResponse struct {
	ResponseID   int64     `db:"response_id"`   // Unique identifier.
	FormID       int64     `db:"form_id"`       // form id of this response.
	RespondentID string    `db:"respondent_id"` // unique identifir for end user.
	CreatedAt    time.Time `db:"created_at"`    // When the response was added.
	UpdatedAt    time.Time `db:"updated_at"`    // When the response record was last modified.
}

func toDBResponse(prd response.Response) dbResponse {
	prdDB := dbResponse{
		ResponseID:   prd.ResponseID,
		FormID:       prd.FormID,
		RespondentID: prd.RespondentID,
		CreatedAt:    prd.CreatedAt.UTC(),
		UpdatedAt:    prd.UpdatedAt.UTC(),
	}

	return prdDB
}

func toCoreResponse(dbItm dbResponse) response.Response {
	prd := response.Response{
		ResponseID:   dbItm.ResponseID,
		FormID:       dbItm.FormID,
		RespondentID: dbItm.RespondentID,
		CreatedAt:    dbItm.CreatedAt.In(time.Local),
		UpdatedAt:    dbItm.UpdatedAt.In(time.Local),
	}

	return prd
}

func toCoreResponseSlice(dbItems []dbResponse) []response.Response {
	prds := make([]response.Response, len(dbItems))
	for i, dbPrd := range dbItems {
		prds[i] = toCoreResponse(dbPrd)
	}
	return prds
}
