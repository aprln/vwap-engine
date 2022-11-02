package config

import "github.com/aprln/vwap-engine/internal/env"

type FeedName string

const (
	FeedNameCoinbase FeedName = "coinbase"
)

// set "wss://ws-feed.exchange.coinbase.com" by default for convenience only
// should not do this in real code
const deftFeedWSConnectionURL = "wss://ws-feed.exchange.coinbase.com"

func NewFeed() Feed {
	return Feed{
		Name:            FeedName(env.LoadEnvString("FEED_NAME", string(FeedNameCoinbase))),
		WSConnectionURL: env.LoadEnvString("FEED_WS_CONNECTION_URL", deftFeedWSConnectionURL),
	}
}

type Feed struct {
	Name            FeedName
	WSConnectionURL string
}
