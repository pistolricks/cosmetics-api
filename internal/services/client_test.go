package services

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

type captureTransport struct {
	t           *testing.T
	lastReq     *http.Request
	expectedURL string
}

func (ct *captureTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	ct.lastReq = req
	if ct.expectedURL != "" && req.URL.String() != ct.expectedURL {
		ct.t.Fatalf("unexpected URL: got %s want %s", req.URL.String(), ct.expectedURL)
	}
	if req.Method != http.MethodPost {
		ct.t.Fatalf("expected POST, got %s", req.Method)
	}
	body := io.NopCloser(bytes.NewBufferString(`{"data":{}}`))
	return &http.Response{StatusCode: 200, Body: body, Header: http.Header{"Content-Type": []string{"application/json"}}}, nil
}

func TestAuthHeaderWithToken(t *testing.T) {
	ct := &captureTransport{t: t, expectedURL: "https://example.myshopify.com/admin/api/2025-07/graphql.json"}
	c := NewClient("example", WithToken("ABC"), WithVersion("2025-07"), WithTransport(ct))
	var out struct{}
	if err := c.QueryString(context.Background(), "query { shop { name } }", nil, &out); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if ct.lastReq == nil {
		t.Fatal("no request captured")
	}
	if got := ct.lastReq.Header.Get(shopifyAccessTokenHeader); got != "ABC" {
		t.Fatalf("missing or wrong token header: %q", got)
	}
	if auth := ct.lastReq.Header.Get("Authorization"); auth != "" {
		t.Fatalf("unexpected Authorization header when using token: %q", auth)
	}
}

func TestAuthHeaderWithPrivateApp(t *testing.T) {
	ct := &captureTransport{t: t, expectedURL: "https://example.myshopify.com/admin/api/2025-07/graphql.json"}
	c := NewClient("example", WithPrivateAppAuth("API_KEY", "ACCESS"), WithVersion("2025-07"), WithTransport(ct))
	var out struct{}
	if err := c.QueryString(context.Background(), "query { shop { name } }", nil, &out); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if ct.lastReq == nil {
		t.Fatal("no request captured")
	}
	if got := ct.lastReq.Header.Get(shopifyAccessTokenHeader); got != "" {
		t.Fatalf("unexpected X-Shopify-Access-Token header: %q", got)
	}
	if auth := ct.lastReq.Header.Get("Authorization"); !strings.HasPrefix(auth, "Basic ") {
		t.Fatalf("expected Basic Authorization header, got: %q", auth)
	}
}
