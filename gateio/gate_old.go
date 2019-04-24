package gateio

import (
	"crypto/hmac"
	"crypto/sha512"
	// "encoding/hex"
	// "encoding/json"
	"net/http"
	// "net/url"
	// "sort"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

const KEY = ""    // gate.io api key
const SECRET = "" // gate.io api secret

// get deposit address
func depositAddress(currency string) string {
	var method string = "POST"
	var url string = "https://api.gateio.io/api2/1/private/depositAddress"
	var param string = "currency=" + currency
	var ret string = httpDo(method, url, param)
	return ret
}

// get deposit withdrawal history
func depositsWithdrawals(start string, end string) string {
	var method string = "POST"
	var url string = "https://api.gateio.io/api2/1/private/depositsWithdrawals"
	var param string = "start=" + start + "&end=" + end
	var ret string = httpDo(method, url, param)
	return ret
}

// Place order buy
func buy(currencyPair string, rate string, amount string) string {
	var method string = "POST"
	var url string = "https://api.gateio.io/api2/1/private/buy"
	var param string = "currencyPair=" + currencyPair + "&rate=" + rate + "&amount=" + amount
	var ret string = httpDo(method, url, param)
	return ret
}

// Place order sell
func sell(currencyPair string, rate string, amount string) string {
	var method string = "POST"
	var url string = "https://api.gateio.io/api2/1/private/sell"
	var param string = "currencyPair=" + currencyPair + "&rate=" + rate + "&amount=" + amount
	var ret string = httpDo(method, url, param)
	return ret
}

// Cancel order
func cancelOrder(orderNumber string, currencyPair string) string {
	var method string = "POST"
	var url string = "https://api.gateio.io/api2/1/private/cancelOrder"
	var param string = "orderNumber=" + orderNumber + "&currencyPair=" + currencyPair
	var ret string = httpDo(method, url, param)
	return ret
}

// Cancel all orders
func cancelAllOrders(types string, currencyPair string) string {
	var method string = "POST"
	var url string = "https://api.gateio.io/api2/1/private/cancelAllOrders"
	var param string = "type=" + types + "&currencyPair=" + currencyPair
	var ret string = httpDo(method, url, param)
	return ret
}

// Get order status
func getOrder(orderNumber string, currencyPair string) string {
	var method string = "POST"
	var url string = "https://api.gateio.io/api2/1/private/getOrder"
	var param string = "orderNumber=" + orderNumber + "&currencyPair=" + currencyPair
	var ret string = httpDo(method, url, param)
	return ret
}

// Get my open order list
func openOrders() string {
	var method string = "POST"
	var url string = "https://api.gateio.io/api2/1/private/openOrders"
	var param string = ""
	var ret string = httpDo(method, url, param)
	return ret
}

// 获取我的24小时内成交记录
func myTradeHistory(currencyPair string, orderNumber string) string {
	var method string = "POST"
	var url string = "https://api.gateio.io/api2/1/private/tradeHistory"
	var param string = "orderNumber=" + orderNumber + "&currencyPair=" + currencyPair
	var ret string = httpDo(method, url, param)
	return ret
}

// Get my last 24h trades
func withdraw(currency string, amount string, address string) string {
	var method string = "POST"
	var url string = "https://api.gateio.io/api2/1/private/withdraw"
	var param string = "currency=" + currency + "&amount=" + amount + "address=" + address
	var ret string = httpDo(method, url, param)
	return ret
}

func getSign(params string) string {
	key := []byte(SECRET)
	mac := hmac.New(sha512.New, key)
	mac.Write([]byte(params))
	return fmt.Sprintf("%x", mac.Sum(nil))
}

/**
*  http request
 */
func httpDo(method string, url string, param string) string {
	client := &http.Client{}

	req, err := http.NewRequest(method, url, strings.NewReader(param))
	if err != nil {
		// handle error
	}
	var sign string = getSign(param)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("key", KEY)
	req.Header.Set("sign", sign)

	resp, err := client.Do(req)
	if err != nil {
		return ""
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	return string(body)
}

type DepthResult struct {
	Elapsed string     `json:"elapsed"`
	Asks    [][]string `json:"asks"`
	Bids    [][]string `json:"bids"`
	Result  string     `json:"result"`
}

func GetMarketPrice(symbol string) float64 {
	var method string = "GET"
	var url string = "http://data.gateio.io/api2/1/orderBook/" + symbol
	var param string = ""
	var ret string = httpDo(method, url, param)
	res := new(DepthResult)
	json.Unmarshal([]byte(ret), res)
	if len(res.Bids) > 0 && len(res.Asks) > 0 {
		ask, _ := strconv.ParseFloat(res.Asks[len(res.Asks)-1][0], 64)
		bid, _ := strconv.ParseFloat(res.Bids[0][0], 64)
		return (ask + bid) / 2
	}
	return 0.0
}
