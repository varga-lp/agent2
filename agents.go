package agent2

// agent has rsi, bb indicators
// indicators sorted by their period asc as an array
// both bb, rsi calc speed is similar, min period is
// ~5x faster than max period
// agent has tp, sl, pos. expiry millis, bakoff millis
// --dependent on position and trade implementation
