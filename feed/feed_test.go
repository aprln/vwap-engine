package feed

import (
	"testing"

	"github.com/aprln/vwap-engine/config"
	"github.com/aprln/vwap-engine/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFeed_GoFeed(t *testing.T) {
	mockWSClient := NewMockWSClient()
	mockResp := mockWSClient.GetTradeResponse()
	fd, err := New(config.Feed{}, config.VWAP{}, mockWSClient, "BTC-USD")
	require.NoError(t, err)

	out := fd.GoFeed()
	got := <-out

	err = mockWSClient.Close()
	require.NoError(t, err)

	assert.Equal(
		t,
		model.Trade{
			TradingPair: mockResp.TradingPair,
			Size:        mockResp.Size,
			Price:       mockResp.Price,
			Time:        mockResp.Time,
		},
		got,
	)
}
