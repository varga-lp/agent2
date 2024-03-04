package agent2

// agent has rsi, bb indicators
// indicators sorted by their period asc as an array
// both bb, rsi calc speed is similar, min period is
// ~5x faster than max period
// agent has tp, sl, pos. expiry millis, backoff millis
// --dependent on position and trade implementation

// agent needs to return open, close signals
// and when returning close signals, need to return a closingReason

type ClosingReason uint8

const (
	Expiry ClosingReason = iota
	StopLoss
	TakeProfit
)
