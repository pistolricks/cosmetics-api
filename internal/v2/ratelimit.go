package v2

import (
	"math"
	"time"

	"github.com/spf13/cast"
)

// CalculateWaitTime returns a duration needed to wait in order to avoid reaching rate limit.
// respExt is the data of the “extensions“ field in Shopify GraphQL response:
//
//	"cost": {
//	  "requestedQueryCost": 101,
//	  "actualQueryCost": 46,
//	  "throttleStatus": {
//	    "maximumAvailable": 1000,
//	    "currentlyAvailable": 954,
//	    "restoreRate": 50
//	  }
//	}
func CalculateWaitTime(respExt map[string]any) time.Duration {
	v, ok := respExt["cost"]
	if !ok {
		return 0
	}
	costData, ok := v.(map[string]any)
	if !ok {
		return 0
	}
	v, ok = costData["throttleStatus"]
	if !ok {
		return 0
	}
	throttleStatus, ok := v.(map[string]any)
	if !ok {
		return 0
	}
	v, _ = costData["requestedQueryCost"]
	requestedQueryCost := cast.ToInt(v)
	v, _ = throttleStatus["currentlyAvailable"]
	currentlyAvailable := cast.ToInt(v)
	if currentlyAvailable >= requestedQueryCost {
		return 0
	}
	v, _ = throttleStatus["restoreRate"]
	restoreRate := cast.ToInt(v)
	lacking := requestedQueryCost - currentlyAvailable
	waitSec := math.Ceil(float64(lacking) / float64(restoreRate))
	return time.Duration(waitSec) * time.Second
}
