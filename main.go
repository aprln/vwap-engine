package main

import (
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/aprln/vwap-engine/config"
	"github.com/aprln/vwap-engine/feed"
	"github.com/aprln/vwap-engine/processor"
	"github.com/aprln/vwap-engine/publisher"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	go func() {
		<-interrupt
		os.Exit(0)
	}()

	var wg sync.WaitGroup

	start(&wg)

	wg.Wait()
}

func start(wg *sync.WaitGroup) {
	cfg := config.NewConfig()

	wg.Add(len(cfg.VWAP.TradingPairs))

	for _, tradingPair := range cfg.VWAP.TradingPairs {
		fd, proc, pub := setup(cfg, tradingPair)
		go func() {
			pub.GoPublish(proc.GoProcess(fd.GoFeed()), wg)
		}()
	}
}

func setup(cfg config.Config, tradingPair string) (feed.Feed, processor.Processor, publisher.Publisher) {
	fd, err := feed.SetUp(cfg.Feed, cfg.VWAP, tradingPair)
	if err != nil {
		log.Fatalf("failed to create a feed: %v", err)
	}

	proc, err := processor.SetUp(cfg.VWAP)
	if err != nil {
		log.Fatalf("failed to create a processor: %v", err)
	}

	pub := publisher.SetUp()

	return fd, proc, pub
}
