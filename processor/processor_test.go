package processor

import (
	"testing"

	"github.com/aprln/vwap-engine/config"
	"github.com/aprln/vwap-engine/model"
	"github.com/stretchr/testify/assert"
)

func TestProcessor_GoProcess(t *testing.T) {
	in := make(chan model.Trade, 1)

	in <- model.Trade{TradingPair: "BTC-USD"}

	out := New(config.VWAP{}, NewMockVWAPCalc()).GoProcess(in)

	vwap := <-out

	close(in)

	assert.Equal(t, "1.1", vwap.VWAP.String())
}
