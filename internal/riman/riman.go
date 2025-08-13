package riman

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-rod/rod/lib/proto"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Riman struct {
	Session  SessionModel
	Clients  ClientModel
	Products ProductModel
	Orders   OrderModel
	Carts    CartModel
	Shipping ShipModel
}

func NewRiman(db *sql.DB) Riman {
	return Riman{
		Session:  SessionModel{DB: db},
		Clients:  ClientModel{DB: db},
		Products: ProductModel{DB: db},
		Orders:   OrderModel{DB: db},
		Carts:    CartModel{DB: db},
		Shipping: ShipModel{DB: db},
	}
}

func CookieStatus(s proto.NetworkCookieSameSite) (int, error) {
	switch s {
	case "Strict":
		return 3, nil
	case "Lax":
		return 2, nil
	case "None":
		return 4, nil
	// Return a zero value for Status and an error for invalid input.
	default:
		return 1, fmt.Errorf("unknown status: %q", s)
	}
}

func restyCookies(cookies []*proto.NetworkCookie) []*http.Cookie {

	var updatedCookies []*http.Cookie

	for _, cookie := range cookies {

		status, err := CookieStatus(cookie.SameSite)
		if err != nil {
			fmt.Println(err)
		}

		var epochSeconds = int64(cookie.Expires)

		t := time.Unix(epochSeconds, 0)

		updatedCookies = append(updatedCookies, &http.Cookie{
			Name:     cookie.Name,
			Value:    cookie.Value,
			Domain:   cookie.Domain,
			Path:     cookie.Path,
			Secure:   cookie.Secure,
			HttpOnly: cookie.HTTPOnly,
			SameSite: http.SameSite(status),
			Expires:  t,
		})
	}

	return updatedCookies
}
