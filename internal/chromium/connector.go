package chromium

import (
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Connector struct {
	// Client ClientConnector
}

func NewConnector() Connector {
	return Connector{
		// Client: ClientConnector{},
	}
}
