package config

import (
	"strings"

	"github.com/aprln/vwap-engine/internal/env"
)

const (
	deftTradingPairs = "BTC-USD|ETH-USD|ETH-BTC"
	deftWindowSize   = 200
)

func NewVWAP() VWAP {
	return VWAP{
		TradingPairs: env.LoadEnvStringSlice("VWAP_TRADING_PAIRS", strings.Split(deftTradingPairs, "|")),
		WindowSize:   env.MustLoadEnvPositiveInt("VWAP_WINDOW_SIZE", deftWindowSize),
	}
}

type VWAP struct {
	TradingPairs []string
	WindowSize   int
}
