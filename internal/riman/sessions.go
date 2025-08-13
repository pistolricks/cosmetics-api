package vendors

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/pistolricks/kbeauty-api/internal/validator"
)

const (
	ScopeAuthentication = "authentication"
)

var (
	ErrDuplicateToken = errors.New("duplicate Token")
)

type Session struct {
	Plaintext string    `json:"token"`
	Hash      []byte    `json:"-"`
	ClientID  int64     `json:"-"`
	Expiry    time.Time `json:"expiry"`
	Scope     string    `json:"-"`
	CartKey   string    `json:"cart_key"`
	Data      []byte    `json:"data"`
}

func ValidateSessionPlaintext(v *validator.Validator, tokenPlaintext string) {
	v.Check(tokenPlaintext != "", "token", "must be provided")
	v.Check(len(tokenPlaintext) == 26, "token", "must be 26 bytes long")
}

type SessionModel struct {
	DB *sql.DB
}

func generateRimanSession(clientID int64, ttl time.Duration, scope string, plainText string, cartKey string, data map[string]any) (*Session, error) {

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	token := &Session{
		ClientID:  clientID,
		Expiry:    time.Now().Add(ttl),
		Scope:     scope,
		Plaintext: plainText,
		CartKey:   cartKey,
		Data:      jsonData,
	}

	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]

	return token, nil
}

func (m SessionModel) NewRimanSession(clientID int64, ttl time.Duration, scope string, plainText string, cartKey string, data map[string]any) (*Session, error) {

	token, err := generateRimanSession(clientID, ttl, scope, plainText, cartKey, data)
	if err != nil {
		return nil, err
	}

	isSaved, session, err := m.SessionCheck(token)

	if isSaved {
		return session, nil
	} else {
		err = m.Insert(token)
		return token, err
	}
}

func (m SessionModel) SessionCheck(session *Session) (bool, *Session, error) {
	queryCheck := `
        SELECT hash, client_id, expiry, scope, cart_key, data
        FROM sessions
        WHERE hash = $1`

	var checkedSession Session

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, queryCheck, session.Hash).Scan(&checkedSession.Hash, &checkedSession.ClientID, &checkedSession.Expiry, &checkedSession.Scope, &checkedSession.CartKey, &checkedSession.Data)

	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return true, &checkedSession, ErrDuplicateToken
		default:
			return true, session, err
		}
	}

	return false, session, nil
}

func (m SessionModel) Insert(session *Session) error {
	query := `
        INSERT INTO sessions (hash, client_id, expiry, scope, cart_key, data)  
        VALUES ($1, $2, $3, $4, $5, $6)`

	args := []any{session.Hash, session.ClientID, session.Expiry, session.Scope, session.CartKey, session.Data}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, args...)
	return err
}

func (m SessionModel) DeleteAllForUser(scope string, clientID int64) error {
	query := `
        DELETE FROM sessions 
        WHERE scope = $1 AND client_id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, scope, clientID)
	return err
}
