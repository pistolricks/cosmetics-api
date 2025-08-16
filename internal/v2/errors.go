package v2

import (
	"strings"
)

func IsConnectionError(err error) bool {
	return err != nil && (strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "broken pipe"))
}
