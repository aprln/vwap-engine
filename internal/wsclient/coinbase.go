package wsclient

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/shopspring/decimal"
)

type CoinbaseChannelName string
type CoinbaseProductID string
type CoinbaseResponseType string
type CoinbaseRequestType string

const (
	CoinbaseChannelNameMatches CoinbaseChannelName = "matches"
)

const (
	CoinbaseResponseTypeMatch     CoinbaseResponseType = "match"
	CoinbaseResponseTypeLastMatch CoinbaseResponseType = "last_match"
)

const (
	CoinbaseRequestTypeSubscribe CoinbaseRequestType = "subscribe"
)

type CoinbaseRequest struct {
	Type       CoinbaseRequestType   `json:"type"`
	ProductIDs []CoinbaseProductID   `json:"product_ids,omitempty"`
	Channels   []CoinbaseChannelName `json:"channels,omitempty"`
}

type CoinbaseMatchesResponse struct {
	Type      CoinbaseResponseType `json:"type"`
	ProductID CoinbaseProductID    `json:"product_id"`
	Size      decimal.Decimal      `json:"size"`
	Price     decimal.Decimal      `json:"price"`
	Time      time.Time            `json:"time"`
}

func NewCoinbase(connURL string) *Coinbase {
	return &Coinbase{connURL: connURL}
}

type Coinbase struct {
	connURL string
	conn    *websocket.Conn
}

func (c *Coinbase) Connect() error {
	// The "Sec-WebSocket-Extensions" header allows for message compression
	// which can increase total throughput and potentially reduce message delivery latency.
	// Ref: https://docs.cloud.coinbase.com/exchange/docs/websocket-overview#websocket-compression-extension
	conn, _, err := websocket.DefaultDialer.DialContext(
		context.Background(),
		c.connURL,
		http.Header{"Sec-WebSocket-Extensions": {"permessage-deflate"}},
	)
	if err != nil {
		return fmt.Errorf("failed to connect with URL %s, %v", c.connURL, err)
	}

	//log.Println("coinbase WS connection established")
	c.conn = conn

	return nil
}

func (c *Coinbase) SubscribeToMatchesChannel(tradingPair string) error {
	if err := c.conn.WriteJSON(
		CoinbaseRequest{
			Type:       CoinbaseRequestTypeSubscribe,
			ProductIDs: []CoinbaseProductID{CoinbaseProductID(tradingPair)},
			Channels:   []CoinbaseChannelName{CoinbaseChannelNameMatches},
		},
	); err != nil {
		return fmt.Errorf(
			`failed to subscribe to trading pair "%s" on channel "%s"`,
			tradingPair,
			CoinbaseChannelNameMatches,
		)
	}

	log.Printf(`subscribed to product "%s" on channel "%s"`, tradingPair, CoinbaseChannelNameMatches)

	return nil
}

func (c *Coinbase) ReadTrade() (TradeResponse, bool, error) {
	resp := &CoinbaseMatchesResponse{}
	err := c.conn.ReadJSON(resp)
	if err != nil {
		return TradeResponse{}, false, err
	}

	if resp.Type != CoinbaseResponseTypeMatch && resp.Type != CoinbaseResponseTypeLastMatch {
		return TradeResponse{}, false, nil
	}

	return TradeResponse{
		TradingPair: string(resp.ProductID),
		Size:        resp.Size,
		Price:       resp.Price,
		Time:        resp.Time,
	}, true, nil
}

func (c *Coinbase) Close() error {
	return c.conn.Close()
}
