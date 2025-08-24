package services

import "testing"

func TestBuildAPIEndpoint(t *testing.T) {
	shop := "example"
	endpoint := buildAPIEndpoint(shop, "admin/api/2025-07")
	expected := "https://example.myshopify.com/admin/api/2025-07/graphql.json"
	if endpoint != expected {
		t.Fatalf("unexpected endpoint: got %s want %s", endpoint, expected)
	}
}

func TestShopFullNameNormalization(t *testing.T) {
	if got := shopFullName("example"); got != "example.myshopify.com" {
		t.Fatalf("unexpected: %s", got)
	}
	if got := shopFullName("example.myshopify.com"); got != "example.myshopify.com" {
		t.Fatalf("unexpected: %s", got)
	}
}
