package feed

import (
	"fmt"
	"log"

	"github.com/aprln/vwap-engine/config"
	"github.com/aprln/vwap-engine/internal/wsclient"
	"github.com/aprln/vwap-engine/model"
)

type WSClient interface {
	Connect() error
	SubscribeToMatchesChannel(tradingPair string) error
	ReadTrade() (wsclient.TradeResponse, bool, error)
	Close() error
}

func SetUp(feedCfg config.Feed, vwapCfg config.VWAP, tradingPair string) (Feed, error) {
	var ws WSClient
	switch feedCfg.Name {
	case config.FeedNameCoinbase:
		ws = wsclient.NewCoinbase(feedCfg.WSConnectionURL)

	default:
		return Feed{}, fmt.Errorf(`feed "%s" is unsupported`, feedCfg.Name)
	}

	return New(feedCfg, vwapCfg, ws, tradingPair)
}

func New(
	feedCfg config.Feed,
	vwapCfg config.VWAP,
	wsClient WSClient,
	tradingPair string,
) (Feed, error) {
	if err := wsClient.Connect(); err != nil {
		return Feed{}, err
	}

	if err := wsClient.SubscribeToMatchesChannel(tradingPair); err != nil {
		return Feed{}, err
	}

	return Feed{
		wsClient:    wsClient,
		feedCfg:     feedCfg,
		vwapCfg:     vwapCfg,
		tradingPair: tradingPair,
	}, nil
}

type Feed struct {
	wsClient    WSClient
	feedCfg     config.Feed
	vwapCfg     config.VWAP
	tradingPair string
}

func (f Feed) GoFeed() chan model.Trade {
	out := make(chan model.Trade, 1)

	go f.feedForever(out)

	return out
}

func (f Feed) feedForever(out chan<- model.Trade) {
	defer close(out)

	for {
		resp, isTradeMsg, err := f.wsClient.ReadTrade()
		if err != nil {
			log.Printf(`failed to read from ws client: "%s"`, err)
			break
		}

		if !isTradeMsg {
			continue
		}

		out <- model.Trade{
			TradingPair: resp.TradingPair,
			Price:       resp.Price,
			Size:        resp.Size,
			Time:        resp.Time,
		}
	}
}
