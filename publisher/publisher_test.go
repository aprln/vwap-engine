package publisher

import (
	"encoding/json"
	"sync"
	"testing"
	"time"

	"github.com/aprln/vwap-engine/model"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPublisher_GoPublish(t *testing.T) {
	in := make(chan model.VWAP, 1)
	mockSender := NewMockSender()
	vwap := model.VWAP{
		TradingPair: "BTC-USD",
		LastTradeAt: time.Date(2020, 11, 1, 1, 1, 1, 1, time.UTC),
		VWAP:        decimal.NewFromFloat(1.1),
	}
	wantMsg, err := json.Marshal(vwap)
	require.NoError(t, err)

	in <- vwap

	var wg sync.WaitGroup
	wg.Add(1)
	New(mockSender).GoPublish(in, &wg)

	gotMsg := mockSender.Read()
	mockSender.Close()
	close(in)

	assert.Equal(t, gotMsg, wantMsg)
}
