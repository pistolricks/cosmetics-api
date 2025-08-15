package riman

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/pistolricks/cosmetics-api/internal/validator"
	"golang.org/x/crypto/bcrypt"
	"resty.dev/v3"
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

func (s *Session) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(s.Hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

type password struct {
	plaintext *string
	hash      []byte
}

type SessionModel struct {
	DB      *sql.DB
	CartKey string
	Token   string
}

func (m SessionModel) CheckForDuplicateToken(hash []byte) (bool, error) {
	query := `
		SELECT EXISTS(SELECT 1 FROM sessions WHERE hash = $1)`

	var exists bool

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, hash).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
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

	exists, err := m.CheckForDuplicateToken(token.Hash)
	if err != nil {
		return nil, err
	}
	if exists {
		// Refresh the existing session row (update expiry, scope, cart_key, data, client_id) and return the token
		if err := m.UpdateByHash(token); err != nil {
			return nil, err
		}
		return token, nil
	}

	err = m.Insert(token)
	return token, err

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

var AnonymousClient = &Client{}

func (c *Client) IsAnonymous() bool {
	return c == AnonymousClient
}

type ClientCredentials struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

type LoggedInResponse struct {
	SecurityRedirect bool   `json:"securityRedirect"`
	Status           string `json:"-"`
	LiToken          string `json:"liToken"`
	LiUser           string `json:"liUser"`
	Jwt              string `json:"jwt"`
}

const loginUrl = "https://security-api.riman.com/api/v2/CheckAttemptsAndLogin"

func (m SessionModel) Login(userName string, password string, token string) (*LoggedInResponse, error) {

	client := resty.New()
	defer client.Close()

	res, err := client.R().
		SetHeader("Accept", "application/json").
		SetAuthToken(token).
		SetBody(ClientCredentials{
			UserName: userName,
			Password: password,
		}).
		SetResult(&LoggedInResponse{}).
		SetError(&Errors{}).
		Post(loginUrl)

	return res.Result().(*LoggedInResponse), err
}

func (m SessionModel) UpdateByHash(session *Session) error {
	query := `
        UPDATE sessions 
        SET client_id = $2, expiry = $3, scope = $4, cart_key = $5, data = $6
        WHERE hash = $1`

	args := []any{session.Hash, session.ClientID, session.Expiry, session.Scope, session.CartKey, session.Data}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, args...)
	return err
}
