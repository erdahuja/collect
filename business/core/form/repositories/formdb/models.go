package formdb

import (
	"time"

	"collect/business/core/form"
)

// dbForm represents an individual form.
type dbForm struct {
	ID          int64     `db:"form_id"`          // Unique identifier.
	Title       string    `db:"form_title"`       // Display name of the form.
	Description string    `db:"form_description"` // Descp for one form.
	DateCreated time.Time `db:"created_at"`       // When the form was added.
	DateUpdated time.Time `db:"updated_at"`       // When the form record was last modified.
}

func toDBForm(frm form.Form) dbForm {
	prdDB := dbForm{
		Title:       frm.Title,
		Description: frm.Description,
		DateCreated: frm.DateCreated.UTC(),
		DateUpdated: frm.DateUpdated.UTC(),
	}

	return prdDB
}

func toCoreForm(dbFrm dbForm) form.Form {
	itm := form.Form{
		ID:          dbFrm.ID,
		Title:       dbFrm.Title,
		Description: dbFrm.Description,
		DateCreated: dbFrm.DateCreated.In(time.Local),
		DateUpdated: dbFrm.DateUpdated.In(time.Local),
	}

	return itm
}

func toCoreFormSlice(dbFrms []dbForm) []form.Form {
	frms := make([]form.Form, len(dbFrms))
	for i, dbFrm := range dbFrms {
		frms[i] = toCoreForm(dbFrm)
	}
	return frms
}
