package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

const (
	lastMatchMsg = `{"type":"last_match","trade_id":443907480,"maker_order_id":"746a0f12-e2b3-4b0e-9538-1d3d5015b7e6","taker_order_id":"b6a0e535-be60-4403-b84c-d2f1b4913e3b","side":"sell","size":"0.0043007","price":"20433.31","product_id":"BTC-USD","sequence":49509894759,"time":"2022-11-02T14:20:45.505047Z"}`
	matchMsg     = `{"type":"match","trade_id":443910503,"maker_order_id":"2b455dd8-7c17-40dc-94fb-1cb658e2b378","taker_order_id":"0a6bdbae-da23-4119-9b51-892e8e7fec9f","side":"buy","size":"0.19671748","price":"20405.75","product_id":"BTC-USD","sequence":49510442399,"time":"2022-11-02T14:27:48.932205Z"}`

	vwapMsg1 = `{"trading_pair":"BTC-USD","last_trade_at":"2022-11-02T14:20:45.505047Z","vwap":"20433.31"}`
	vwapMsg2 = `{"trading_pair":"BTC-USD","last_trade_at":"2022-11-02T14:27:48.932205Z","vwap":"20406.3396346887629766"}`
	vwapMsg3 = `{"trading_pair":"BTC-USD","last_trade_at":"2022-11-02T14:27:48.932205Z","vwap":"20406.0480051926950679"}`
)

func Test_main(t *testing.T) {
	// set up a fake WS server that mimics Coinbase WS server
	svr := httptest.NewServer(http.HandlerFunc(pushFakeResponse))
	defer svr.Close()
	connURL := strings.Replace(svr.URL, "http://", "ws://", 1)

	// configure the app to point to the fake server
	t.Setenv("FEED_WS_CONNECTION_URL", connURL)

	// pipe stdout to a channel
	w, out := pipeStdoutToChan()

	// set expectation
	want := ""
	for i := 0; i < 3; i++ {
		want += vwapMsg1 + "\n" + vwapMsg2 + "\n" + vwapMsg3 + "\n"
	}

	// now run the app
	main()

	// read output
	w.Close()
	got := <-out

	// assertion
	assert.Equal(t, want, got)
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

func pushFakeResponse(w http.ResponseWriter, r *http.Request) {
	upd := websocket.Upgrader{}
	c, err := upd.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer c.Close()

	for i := 0; i < 3; i++ {
		msg := matchMsg
		if i == 0 {
			msg = lastMatchMsg
		}
		err = c.WriteMessage(1, []byte(msg))
		if err != nil {
			break
		}
	}
}
