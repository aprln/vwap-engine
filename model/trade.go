package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type Trade struct {
	TradingPair string
	Size        decimal.Decimal
	Price       decimal.Decimal
	Time        time.Time
}
