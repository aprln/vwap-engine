package processor

import (
	"log"

	"github.com/aprln/vwap-engine/config"
	"github.com/aprln/vwap-engine/model"
	"github.com/shopspring/decimal"
)

type VWAPCalculator interface {
	VWAP() decimal.Decimal
	AddDataPoint(price, size decimal.Decimal) error
}

func SetUp(vwapCfg config.VWAP) (Processor, error) {
	c, err := NewVWAPCalc(vwapCfg.WindowSize)
	if err != nil {
		return Processor{}, err
	}

	return New(vwapCfg, c), nil
}

func New(vwapCfg config.VWAP, calc VWAPCalculator) Processor {
	return Processor{
		vwapCfg: vwapCfg,
		calc:    calc,
	}
}

type Processor struct {
	calc    VWAPCalculator
	vwapCfg config.VWAP
}

func (p Processor) GoProcess(in <-chan model.Trade) chan model.VWAP {
	out := make(chan model.VWAP, 1)

	go p.processForever(in, out)

	return out
}

func (p Processor) processForever(in <-chan model.Trade, out chan<- model.VWAP) {
	defer close(out)

	for {
		trade, more := <-in
		if !more {
			log.Println("no more to read from the trade channel")

			break
		}
		err := p.calc.AddDataPoint(trade.Price, trade.Size)
		if err != nil {
			log.Printf("calculator error %v", err)

			break
		}
		out <- model.VWAP{
			TradingPair: trade.TradingPair,
			LastTradeAt: trade.Time,
			VWAP:        p.calc.VWAP(),
		}
	}
}
