package config

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewVWAP(t *testing.T) {
	tests := []struct {
		name    string
		envVars map[string]string
		want    VWAP
	}{
		{
			name: "no env vars",
			want: VWAP{
				TradingPairs: strings.Split(deftTradingPairs, "|"),
				WindowSize:   deftWindowSize,
			},
		},
		{
			name: "with env vars",
			envVars: map[string]string{
				"VWAP_WINDOW_SIZE":   "2",
				"VWAP_TRADING_PAIRS": "ABC|DEF",
			},
			want: VWAP{
				TradingPairs: strings.Split("ABC|DEF", "|"),
				WindowSize:   2,
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				for k, v := range tt.envVars {
					t.Setenv(k, v)
				}
				got := NewVWAP()
				assert.Equal(t, tt.want, got)
			},
		)
	}
}
