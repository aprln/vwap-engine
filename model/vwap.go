package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type VWAP struct {
	TradingPair string          `json:"trading_pair"`
	LastTradeAt time.Time       `json:"last_trade_at"`
	VWAP        decimal.Decimal `json:"vwap"`
}
