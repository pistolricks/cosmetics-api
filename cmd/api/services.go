package main

import (
	"net/http"
	"os"

	"github.com/pistolricks/cosmetics-api/internal/services"
)

const shopifyAccessTokenHeader = "X-Shopify-Access-Token"

func (app *application) RoundTrip(req *http.Request) (*http.Response, error) {
	isAccessTokenSet := app.transport.accessToken != ""
	areBasicAuthCredentialsSet := app.transport.apiKey != "" && isAccessTokenSet

	if areBasicAuthCredentialsSet {
		req.SetBasicAuth(app.transport.apiKey, app.transport.accessToken)
	} else if isAccessTokenSet {
		req.Header.Set(shopifyAccessTokenHeader, app.transport.accessToken)
	}

	return app.transport.roundTripper.RoundTrip(req)
}

func (app *application) clientWithToken() *services.ClientApi {
	return services.NewClient(os.Getenv("STORE_NAME"), services.WithToken(os.Getenv("SHOPIFY_TOKEN")), services.WithVersion("2025-07"))
}
