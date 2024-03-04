package agent2

import (
	"fmt"

	"github.com/varga-lp/data/klines"
)

type Trade struct {
	OpenTime     int64
	CloseTime    int64
	Reason       ClosingReason
	NetProfit    float64
	DurationSecs int64
}

var (
	ErrPositionCantBeNilForTrade = fmt.Errorf("position can't be nil when creating a trade")
)

func NewTrade(pos *Position, cr ClosingReason, longClose klines.Kline, shortClose klines.Kline) (*Trade, error) {
	if pos == nil {
		return nil, ErrPositionCantBeNilForTrade
	}

	return &Trade{
		OpenTime:     pos.Long.OpenTime,
		CloseTime:    longClose.CloseTime,
		Reason:       cr,
		NetProfit:    pos.NetProfit(longClose, shortClose),
		DurationSecs: (longClose.CloseTime - pos.Long.OpenTime) / 1_000,
	}, nil
}

func (tr *Trade) String() string {
	return fmt.Sprintf("[trd] OT=%d, DSecs=%d, Net=%.2f, R=%s",
		tr.OpenTime, tr.DurationSecs, tr.NetProfit, tr.Reason)
}
