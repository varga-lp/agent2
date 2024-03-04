package agent2

import (
	"fmt"

	"github.com/varga-lp/data/klines"
)

const (
	commission = 0.000450
)

const (
	defaultAllocation = float64(1000.0)
)

type Position struct {
	Long  klines.Kline
	Short klines.Kline
}

var (
	ErrLongShortOpenTimeNotEqual = fmt.Errorf("long, short kline open times not equal")
)

func NewPosition(long klines.Kline, short klines.Kline) (*Position, error) {
	if long.OpenTime != short.OpenTime {
		return nil, ErrLongShortOpenTimeNotEqual
	}

	return &Position{
		Long:  long,
		Short: short,
	}, nil
}

func (p *Position) GrossProfit(long klines.Kline, short klines.Kline) float64 {
	return (long.Close/p.Long.Close + p.Short.Close/short.Close - 2.0) * defaultAllocation / 2.0
}

func (p *Position) NetProfit(long klines.Kline, short klines.Kline) float64 {
	gp := p.GrossProfit(long, short)

	return gp - (defaultAllocation*2+gp)*commission
}

func (p *Position) ExpiredAt(expiryMillis int64, at int64) bool {
	return at > (p.Long.CloseTime + expiryMillis)
}
