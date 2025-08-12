package riman

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Extended struct {
	Session SessionModel
	Clients ClientModel
}

func NewExtended(db *sql.DB) Extended {
	return Extended{
		Session: SessionModel{DB: db},
		Clients: ClientModel{DB: db},
	}
}
