package v2

import (
	"net/http"
)

const shopifyAccessTokenHeader = "X-Shopify-Access-Token"

type transport struct {
	accessToken  string
	apiKey       string
	apiBasePath  string
	roundTripper http.RoundTripper
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	isAccessTokenSet := t.accessToken != ""
	areBasicAuthCredentialsSet := t.apiKey != "" && isAccessTokenSet

	if areBasicAuthCredentialsSet {
		req.SetBasicAuth(t.apiKey, t.accessToken)
	} else if isAccessTokenSet {
		req.Header.Set(shopifyAccessTokenHeader, t.accessToken)
	}

	return t.roundTripper.RoundTrip(req)
}
