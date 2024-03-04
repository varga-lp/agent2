package agent2

import (
	"math/rand"
	"time"
)

type Backoff struct {
	DurationMillis int64
}

const (
	minBackoffMillis = 1 * 60 * 1_000  // 1 minute
	maxBackoffMillis = 30 * 60 * 1_000 // 30 minutes
	backoffStep      = 10 * 1_000      // 10 seconds
)

func randBackoff() int64 {
	r := rand.Intn(maxBackoffMillis - minBackoffMillis)

	return (int64(r)/backoffStep)*backoffStep + minBackoffMillis
}

func RandomBackoff() *Backoff {
	return &Backoff{
		DurationMillis: randBackoff(),
	}
}

func (bo *Backoff) TradeAllowed(lastTrade *Trade) bool {
	if lastTrade == nil {
		return true
	}

	return (time.Now().UnixMilli() - lastTrade.CloseTime) > bo.DurationMillis
}
