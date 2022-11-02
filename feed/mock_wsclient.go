package feed

import (
	"errors"
	"time"

	"github.com/aprln/vwap-engine/internal/wsclient"
	"github.com/shopspring/decimal"
)

func NewMockWSClient() *MockWSClient {
	return &MockWSClient{closed: make(chan bool, 1)}
}

type MockWSClient struct {
	closed chan bool
}

func (m *MockWSClient) Connect() error {
	return nil
}

func (m *MockWSClient) SubscribeToMatchesChannel(tradingPair string) error {
	return nil
}

func (m *MockWSClient) ReadTrade() (wsclient.TradeResponse, bool, error) {
	select {
	case <-m.closed:
		return wsclient.TradeResponse{}, false, errors.New("connection closed")
	default:
		return m.GetTradeResponse(), true, nil
	}
}

func (m *MockWSClient) GetTradeResponse() wsclient.TradeResponse {
	return wsclient.TradeResponse{
		TradingPair: "BTC-USD",
		Size:        decimal.NewFromFloat(1.1),
		Price:       decimal.NewFromFloat(2.2),
		Time:        time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC),
	}
}

func (m *MockWSClient) Close() error {
	m.closed <- true

	return nil
}
