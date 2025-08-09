package main

import (
	"context"
	"net/http"

	"github.com/pistolricks/kbeauty-api/internal/data"
	"github.com/pistolricks/kbeauty-api/internal/riman"
)

type contextKey string

const clientContextKey = contextKey("client")
const userContextKey = contextKey("user")

func (app *application) contextSetClient(r *http.Request, client *riman.Client) *http.Request {
	ctx := context.WithValue(r.Context(), clientContextKey, client)
	return r.WithContext(ctx)
}

func (app *application) contextGetClient(r *http.Request) *riman.Client {
	client, ok := r.Context().Value(clientContextKey).(*riman.Client)
	if !ok {
		panic("missing client value in request context")
	}

	return client
}

func (app *application) contextSetUser(r *http.Request, user *data.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func (app *application) contextGetUser(r *http.Request) *data.User {
	user, ok := r.Context().Value(userContextKey).(*data.User)
	if !ok {
		panic("missing user value in request context")
	}

	return user
}
