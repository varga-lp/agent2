# agent2

The simplest agent library.

Properties

- indicators
  - activation period
- take profit
- stop loss
- pos. expiry
- backoff

- signalOpen(s1, s2 []klines.Kline, lastTrade Trade)
  - data validation
  - backoff validation
  - len s1, s2 should be == activation period
- signalClose(lastS1 klines.Kline, lastS2 klines.Kline, position)
