package agent2

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/varga-lp/data/klines"
)

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
	NoReason
)

func (cr ClosingReason) String() string {
	switch cr {
	case Expiry:
		return "expiry"
	case StopLoss:
		return "stopLoss"
	case TakeProfit:
		return "takeProfit"
	case NoReason:
		return "noReason"
	}
	return ""
}

const (
	maxBBCount  = 5
	maxRSICount = 5
)

type Agent struct {
	Tpsl         *TPSL    `json:"tpsl"`
	Backoff      *Backoff `json:"backoff"`
	ExpiryMillis int64    `json:"expiry_mls"`
	Bbs          []*BB    `json:"bbs"`
	Rsis         []*RSI   `json:"rsis"`
}

func (ag *Agent) Marshal() ([]byte, error) {
	pload, err := json.Marshal(ag)
	if err != nil {
		return nil, err
	}
	return pload, nil
}

func UnmarshalAgent(pload []byte) (*Agent, error) {
	var agent Agent

	if err := json.Unmarshal(pload, &agent); err != nil {
		return nil, err
	}
	return &agent, nil
}

func RandomAgent() *Agent {
	bbMons, rsiMons := make(map[Monitor]struct{}), make(map[Monitor]struct{})

	ag := &Agent{
		Tpsl:         RandomTPSL(),
		Backoff:      RandomBackoff(),
		ExpiryMillis: randExpiry(),
		Bbs:          make([]*BB, 0, maxBBCount),
		Rsis:         make([]*RSI, 0, maxRSICount),
	}

	for i := 0; i < maxBBCount; i++ {
		bb := RandomBB()

		if _, ok := bbMons[bb.Mon]; !ok {
			bbMons[bb.Mon] = struct{}{}

			ag.Bbs = append(ag.Bbs, bb)
		}
	}
	for i := 0; i < maxRSICount; i++ {
		rsi := RandomRSI()

		if _, ok := rsiMons[rsi.Mon]; !ok {
			rsiMons[rsi.Mon] = struct{}{}

			ag.Rsis = append(ag.Rsis, rsi)
		}
	}
	sort.Slice(ag.Bbs[:], func(i, j int) bool {
		return ag.Bbs[i].Period < ag.Bbs[j].Period
	})
	sort.Slice(ag.Rsis[:], func(i, j int) bool {
		return ag.Rsis[i].Period < ag.Rsis[j].Period
	})

	return ag
}

const (
	minActivationKlineLength = 250
)

var (
	ErrKlinesAreBelowMinActivationKlineLength = fmt.Errorf("kline length is below min activation kline length 250")
	ErrPositionCantBeNilForClose              = fmt.Errorf("position can't be nil for close pos")
)

func (ag *Agent) OpenPos(klns1 []klines.Kline, klns2 []klines.Kline, lastTrade *Trade) (bool, error) {
	// check backoff
	if !ag.Backoff.TradeAllowed(lastTrade) {
		return false, nil
	}
	klns1Len, klns2Len := len(klns1), len(klns2)

	if klns1Len < minActivationKlineLength || klns2Len < minActivationKlineLength {
		return false, ErrKlinesAreBelowMinActivationKlineLength
	}

	// check rsi indicators first as its faster than bb
	for _, rsi := range ag.Rsis {
		active, err := rsi.Active(klns1[klns1Len-rsi.Period:], klns2[klns2Len-rsi.Period:])
		if err != nil {
			return false, err
		}
		if !active {
			return false, nil
		}
	}
	// check bb indicators
	for _, bb := range ag.Bbs {
		active, err := bb.Active(klns1[klns1Len-bb.Period:], klns2[klns2Len-bb.Period:])
		if err != nil {
			return false, err
		}
		if !active {
			return false, nil
		}
	}
	return true, nil
}

func (ag *Agent) ClosePos(pos *Position, closeLong klines.Kline, closeShort klines.Kline) (bool, ClosingReason, error) {
	if pos == nil {
		return false, NoReason, ErrPositionCantBeNilForClose
	}

	// check sl
	if clos, err := ag.Tpsl.SLNetClose(pos, closeLong, closeShort); err != nil {
		return false, NoReason, err
	} else if clos {
		return true, StopLoss, nil
	}
	// check tp
	if clos, err := ag.Tpsl.TPNetClose(pos, closeLong, closeShort); err != nil {
		return false, NoReason, err
	} else if clos {
		return true, TakeProfit, nil
	}
	// check expiry
	if pos.ExpiredAt(ag.ExpiryMillis, closeLong.CloseTime) {
		return true, Expiry, nil
	}
	return false, NoReason, nil
}
