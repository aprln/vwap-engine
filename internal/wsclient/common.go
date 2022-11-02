package wsclient

import (
	"time"

	"github.com/shopspring/decimal"
)

type TradeResponse struct {
	TradingPair string
	Size        decimal.Decimal
	Price       decimal.Decimal
	Time        time.Time
}
