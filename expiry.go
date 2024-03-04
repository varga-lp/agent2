package agent2

import "math/rand"

const (
	minExpiryMillis = 45 * 60 * 1_000     // 45 minutes
	maxExpiryMillis = 6 * 60 * 60 * 1_000 // 6 hours
	expiryStep      = 60 * 1_000          // 1 minute
)

func randExpiry() int64 {
	r := rand.Intn(maxExpiryMillis - minExpiryMillis)

	return (int64(r)/expiryStep)*expiryStep + minExpiryMillis
}
