package services

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// exampleTransport is a simple RoundTripper that fakes a successful JSON response.
// It lets us show how to use the client without making real network calls.
type exampleTransport struct {
	lastReq *http.Request
}

func (et *exampleTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	et.lastReq = req
	body := io.NopCloser(bytes.NewBufferString(`{"data":{}}`))
	return &http.Response{StatusCode: 200, Body: body, Header: http.Header{"Content-Type": []string{"application/json"}}}, nil
}

// ExampleClientApi_QueryString demonstrates a read-style GraphQL call (GET-like usage)
// using the ClientApi. Note: GraphQL requests are POST under the hood, but QueryString
// is typically used for read operations.
func ExampleClientApi_QueryString() {
	// Create a client with a token and API version, and a mock transport.
	et := &exampleTransport{}
	c := NewClient("example", WithToken("TOKEN"), WithVersion("2025-07"), WithTransport(et), WithTimeout(3*time.Second))

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var out struct{}
	err := c.QueryString(ctx, "query { shop { name } }", nil, &out)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	// Confirm headers were set as expected (token header should be present)
	fmt.Println(et.lastReq.Header.Get(shopifyAccessTokenHeader) != "")
	fmt.Println("GET example OK")
	// Output:
	// true
	// GET example OK
}

// ExampleClientApi_Mutate demonstrates a write-style GraphQL call (POST-like usage)
// using the ClientApi.Mutate method.
func ExampleClientApi_Mutate() {
	et := &exampleTransport{}
	c := NewClient("example", WithPrivateAppAuth("API_KEY", "ACCESS"), WithVersion("2025-07"), WithTransport(et), WithTimeout(3*time.Second))

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Define a dummy mutation payload
	type mutationEcho struct {
		Echo struct{} `json:"echo"`
	}
	m := &mutationEcho{}

	vars := map[string]interface{}{
		"input": map[string]any{"id": "gid://shopify/Order/123"},
	}

	if err := c.Mutate(ctx, m, vars); err != nil {
		fmt.Println("error:", err)
		return
	}

	// Confirm Authorization header exists for private app auth (Basic ...)
	fmt.Println(et.lastReq.Header.Get("Authorization") != "")
	fmt.Println("POST example OK")
	// Output:
	// true
	// POST example OK
}
