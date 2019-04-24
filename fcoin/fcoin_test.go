package fcoin

import (
	"testing"
	"net/url"
	"encoding/json"
	"fmt"
)

func TestAuthorization(t *testing.T) {
	u := `https://api.fcoin.com`
	path := `/v2/orders`
	values := url.Values{}
	values.Add("type", "limit")
	values.Add("side", "buy")
	values.Add("amount", "100.0")
	values.Add("price", "100.0")
	values.Add("symbol", "btcusdt")
	fs,err := NewFcoinService(u, "", "")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fs.authorization("post", path, nil, values))
}

func TestAuthorization2(t *testing.T) {
	path := `/v2/public/server-time`
	fs,err := NewFcoinService("https://api.fcoin.com", "", "")
	if err != nil {
		t.Fatal(err)
	}
	data, err := fs.authorization("get", path, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(data))

	var ts int64
	json.Unmarshal(data, &ts)
	t.Log(ts)
}

func TestGetServerTime(t *testing.T) {
	fs,err := NewFcoinService("https://api.fcoin.com", "", "")
	if err != nil {
		t.Fatal(err)
	}
	ti, err := fs.GetServerTime()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ti.Format("2006-01-02 15:04:05.000000"))
}

func TestGetCurrencies(t *testing.T) {
	fs,err := NewFcoinService("https://api.fcoin.com", "", "")
	if err != nil {
		t.Fatal(err)
	}
	cs, err := fs.GetCurrencies()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(cs)
}

func TestGetSymbols(t *testing.T) {
	fs,err := NewFcoinService("https://api.fcoin.com", "", "")
	if err != nil {
		t.Fatal(err)
	}
	cs, err := fs.GetSymbols()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(cs)
}

func TestGetMarketTicker(t *testing.T) {
	fs,err := NewFcoinService("https://api.fcoin.com", "", "")
	if err != nil {
		t.Fatal(err)
	}
	cs, err := fs.GetMarketTicker("btcusdt")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(cs)
}

func TestGetMarketDepth(t *testing.T) {
	fs,err := NewFcoinService("https://api.fcoin.com", "", "")
	if err != nil {
		t.Fatal(err)
	}
	cs, err := fs.GetMarketDepth("L20", "gtceth")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(cs.Type)
	fmt.Println(cs.Bids)
	fmt.Println(cs.Asks)
}

func TestGetMarketTrades(t *testing.T) {
	fs,err := NewFcoinService("https://api.fcoin.com", "", "")
	if err != nil {
		t.Fatal(err)
	}
	cs, err := fs.GetMarketTrades("btcusdt", "", 20)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(cs)
}

func TestGetMarketCandle(t *testing.T) {
	fs,err := NewFcoinService("https://api.fcoin.com", "", "")
	if err != nil {
		t.Fatal(err)
	}
	cs, err := fs.GetMarketCandle("M3", "btcusdt", "", 20)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(cs)
}

func TestFcoinService_GetMarketTicker(t *testing.T) {
	fs,err := NewFcoinService("https://api.fcoin.com", "", "")
	if err != nil {
		t.Fatal(err)
	}
	cs, err := fs.GetMarketTicker("fteth")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(cs)
}

func TestFcoinService_GetMarketTrades(t *testing.T) {
	fs,err := NewFcoinService("https://api.fcoin.com", "", "")
	if err != nil {
		t.Fatal(err)
	}
	cs, err := fs.GetMarketTrades("fteth", "",100)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range cs {
		t.Log(v)
	}
}

func TestFcoinService_GetAccountBalance(t *testing.T) {
	fs,err := NewFcoinService("https://api.fcoin.com", "", "")
	if err != nil {
		t.Fatal(err)
	}
	cs, err := fs.GetAccountBalance()
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range cs {
		t.Log(v)
	}
}

func TestFcoinService_GetOrders(t *testing.T) {
	fs,err := NewFcoinService("https://api.fcoin.com", "", "")
	if err != nil {
		t.Fatal(err)
	}
	cs, err := fs.GetOrders("gtcft", "partial_canceled","","0","20")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(cs)
}

func TestFcoinService_CreateOrder(t *testing.T) {
	fs,err := NewFcoinService("https://api.fcoin.com", "", "")
	if err != nil {
		t.Fatal(err)
	}
	orderID, err := fs.CreateOrder("gtcft", "sell", "limit","0.230087","10")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(orderID)
}

func TestFcoinService_CreateOrderData(t *testing.T) {
	fs,err := NewFcoinService("https://api.fcoin.com", "", "")
	if err != nil {
		t.Fatal(err)
	}

	values := url.Values{}
	values.Add("symbol", "zipeth")  // 交易对
	values.Add("side", "sell")      // 交易方向
	values.Add("type", "limit") // 订单类型
	values.Add("price", "0.00000797")    // 价格
	values.Add("amount", "2000")  // 下单量

	path := `/v2/orders`
	data, err := fs.authorization("POST", path, nil, values)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(data)
}

func TestFcoinService_CreateOrder2(t *testing.T) {
	fs,err := NewFcoinService("https://api.fcoin.com", "", "")
	if err != nil {
		t.Fatal(err)
	}
	orderID, err := fs.CreateOrder("fteth", "sell", "limit","0.00145115","5.09")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(orderID)
}

func TestFcoinService_GetOrderByID(t *testing.T) {
	fs,err := NewFcoinService("https://api.fcoin.com", "", "")
	if err != nil {
		t.Fatal(err)
	}
	info, err := fs.GetOrderByID("=")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(info)
	t.Log(info.State)
	t.Log(info.ExecutedValue)
	t.Log(info.FilledAmount)
	t.Log(info.FillFees)
}


func TestFcoinService_GetOrderByID2(t *testing.T) {
	fs,err := NewFcoinService("https://api.fcoin.com", "", "")
	if err != nil {
		t.Fatal(err)
	}
	info, err := fs.GetOrderByID("=")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(info)
	t.Log(info.State)
}

func TestFcoinService_CancelOrder(t *testing.T) {
	fs,err := NewFcoinService("https://api.fcoin.com", "", "")
	if err != nil {
		t.Fatal(err)
	}
	success, err := fs.CancelOrder("=")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(success)
}

func TestFcoinService_HasEnoughAssets(t *testing.T) {
	fs,err := NewFcoinService("https://api.fcoin.com", "", "")
	if err != nil {
		t.Fatal(err)
	}
	success, err := fs.HasEnoughAssets("eth", 9.4)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(success)
}

func TestFcoinService_GetCurrentMarketPrice(t *testing.T) {
	fs,err := NewFcoinService("https://api.fcoin.com", "", "")
	if err != nil {
		t.Fatal(err)
	}
	success, err := fs.GetCurrentMarketPrice("gtcft", 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(success)
}

func TestFcoinService_GetAvailableAmount(t *testing.T) {
	fs,err := NewFcoinService("https://api.fcoin.com", "", "")
	if err != nil {
		t.Fatal(err)
	}
	f, err := fs.GetAvailableAmount("ft")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(f)
	price, err := fs.GetCurrentMarketPrice("gtcft", 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(f / price * 0.95)
}

func TestUrlValuesToJSON(t *testing.T) {
	values := url.Values{}
	values.Add("a", "A")
	values.Add("b", "B")
	t.Log(urlValuesToJSON(values))
}