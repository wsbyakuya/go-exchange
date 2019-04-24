package fcoin

import (
	"testing"
	"net/url"
)

func TestSortedURL(t *testing.T) {
	params := url.Values{}
	params.Add("a", "A")
	params.Add("b", "B")
	params.Add("c", "C")
	t.Log(sortedURI("https://api.fcoin.com/v2/orders", params))
}

func TestSortedBody(t *testing.T) {
	values := url.Values{}
	values.Add("type", "limit")
	values.Add("side", "buy")
	values.Add("amount", "100.0")
	values.Add("price", "100.0")
	values.Add("symbol", "btcusdt")
	t.Log(sortedBody(values))
}
