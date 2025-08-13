package vendors

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Vendors struct {
	Session SessionModel
	Clients ClientModel
}

func NewVendors(db *sql.DB) Vendors {
	return Vendors{
		Session: SessionModel{DB: db},
		Clients: ClientModel{DB: db},
	}
}
