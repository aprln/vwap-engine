package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func Test_main(t *testing.T) {
	// set up a fake WS server that mimics Coinbase WS server
	svr := httptest.NewServer(http.HandlerFunc(pushFakeWSResponse))
	defer svr.Close()
	connURL := strings.Replace(svr.URL, "http://", "ws://", 1)

	// set up configuration and point to the fake WS server above
	t.Setenv("FEED_NAME", "coinbase")
	t.Setenv("FEED_WS_CONNECTION_URL", connURL)
	t.Setenv("VWAP_TRADING_PAIRS", "BTC-USD|ETH-USD")
	t.Setenv("VWAP_WINDOW_SIZE", "3")

	// pipe stdout to a channel
	w, out := pipeStdoutToChan()

	// set expectation
	wantMsgsBTCUSD := []string{
		getVWAPMsg("BTC-USD", "20433.31"),
		getVWAPMsg("BTC-USD", "19427.7342170096256965"),
		getVWAPMsg("BTC-USD", "19786.8530027760498389"),
		getVWAPMsg("BTC-USD", "20016.3010161061277987"),
	}
	wantMsgsETHUSD := []string{
		getVWAPMsg("ETH-USD", "20433.31"),
		getVWAPMsg("ETH-USD", "19427.7342170096256965"),
		getVWAPMsg("ETH-USD", "19786.8530027760498389"),
		getVWAPMsg("ETH-USD", "20016.3010161061277987"),
	}

	// now run the app
	main()

	// read output
	w.Close()
	got := <-out
	gotMsgs := strings.Split(got, "\n")

	// assertion
	assert.Equal(t, wantMsgsBTCUSD, filterMsgsContain(gotMsgs, "BTC-USD"))
	assert.Equal(t, wantMsgsETHUSD, filterMsgsContain(gotMsgs, "ETH-USD"))
}

func pushFakeWSResponse(w http.ResponseWriter, r *http.Request) {
	upd := websocket.Upgrader{}
	conn, err := upd.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	_, reqJSON, err := conn.ReadMessage()
	if err != nil {
		return
	}

	// reqJSON example {"type":"subscribe","product_ids":["BTC-USD"],"channels":["matches"]}
	regData := struct {
		ProductIDs []string `json:"product_ids"`
	}{}
	err = json.Unmarshal(reqJSON, &regData)
	if err != nil {
		return
	}

	tradingPair := regData.ProductIDs[0]

	msgs := []string{
		getMatchResponse(tradingPair, "last_match", "20433.31", "0.0043007"),
		getMatchResponse(tradingPair, "match", "19405.75", "0.19671748"),
		getMatchResponse(tradingPair, "match", "20405.35", "0.11671747"),
		getMatchResponse(tradingPair, "match", "20605.78", "0.1267174"),
	}

	for _, msg := range msgs {
		err = conn.WriteMessage(1, []byte(msg))
		if err != nil {
			break
		}
	}
}

func getMatchResponse(tradingPair, msgType, price, size string) string {
	return fmt.Sprintf(
		`
			{
				"type": "%s",
				"trade_id": 443907480,
				"maker_order_id": "746a0f12-e2b3-4b0e-9538-1d3d5015b7e6",
				"taker_order_id": "b6a0e535-be60-4403-b84c-d2f1b4913e3b",
				"side": "sell",
				"size": "%s",
				"price": "%s",
				"product_id": "%s",
				"sequence": 49509894759,
				"time": "2022-11-02T14:27:48.932205Z"
			}
		`,
		msgType, size, price, tradingPair,
	)
}

func getVWAPMsg(tradingPair, vwap string) string {
	return fmt.Sprintf(
		`{"trading_pair":"%s","last_trade_at":"2022-11-02T14:27:48.932205Z","vwap":"%s"}`,
		tradingPair, vwap,
	)
}

func pipeStdoutToChan() (*os.File, chan string) {
	out := make(chan string)
	r, w, _ := os.Pipe()
	os.Stdout = w

	go func() {
		var buf bytes.Buffer
		_, err := io.Copy(&buf, r)
		if err != nil {
			panic(err)
		}
		out <- buf.String()
	}()

	return w, out
}

func filterMsgsContain(msgs []string, containingStr string) []string {
	var res []string
	for _, msg := range msgs {
		if strings.Contains(msg, containingStr) {
			res = append(res, msg)
		}
	}

	return res
}
