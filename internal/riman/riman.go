package riman

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Riman struct {
	Session SessionModel
	Clients ClientModel
}

func NewRiman(db *sql.DB) Riman {
	return Riman{
		Session: SessionModel{DB: db},
		Clients: ClientModel{DB: db},
	}
}
