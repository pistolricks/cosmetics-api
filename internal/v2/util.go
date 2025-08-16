package v2

import (
	"fmt"
	"strings"
)

// shopFullName returns the full shop name, including .myshopify.com
func shopFullName(name string) string {
	name = strings.TrimSpace(name)
	name = strings.Trim(name, ".")
	if strings.Contains(name, shopifyBaseDomain) {
		return name
	}
	return name + "." + shopifyBaseDomain
}

// shopBaseURL returns the Shop's base URL.
func shopBaseURL(name string) string {
	name = shopFullName(name)
	return fmt.Sprintf("%s://%s", defaultAPIProtocol, name)
}

func buildAPIEndpoint(shopName string, apiPathPrefix string) string {
	baseURL := shopBaseURL(shopName)
	return fmt.Sprintf("%s/%s/%s", baseURL, apiPathPrefix, defaultAPIEndpoint)
}
