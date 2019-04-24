package gateio

import "testing"

func TestGetMarketPrice(t *testing.T) {
	t.Log(GetMarketPrice("gtc_usdt"))
	t.Log(GetMarketPrice("gtc_eth"))
}
