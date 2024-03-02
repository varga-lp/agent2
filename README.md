# agent2

The simplest agent library.

Properties

- indicators (bbs, rsis)
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

 18 type Kline struct {
 19     OpenTime       int64
 20     Open           float64
 21     High           float64
 22     Low            float64
 23     Close          float64
 24     Volume         float64
 25     TakerBuyVolume float64
 26     NumberOfTrades int64
 27     CloseTime      int64
 28     IsFinal        bool
 29 }

Monitors
- Close1
- Close2
- CloseR --priority
- (High - Low)1
- (High - Low)2
- (High - Low)R
- Volume1
- Volume2
- VolumeR --priority
- (TakerBuyVol / Volume)1
- (TakerBuyVol / Volume)2
- (TakerBuyVol / Volume)R --priority
- NumberOfTrades1
- NumberOfTrades2
- NumberOfTradesR --priority

Indicators
- gets klns1, klns2
  - validates length == period
- returns active true, false
