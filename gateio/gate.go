package gateio

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	defaultHOST   = "data.gateio.io"
	defaultScheme = "https"
)

type Service struct {
	apiKey string
	secret string
}

func NewService(apiKey, secret string) *Service {
	return &Service{
		apiKey: apiKey,
		secret: secret,
	}
}

func (s *Service) requestJSON(method, path string, values url.Values, target interface{}) error {
	resp, err := s.doHTTP(method, path, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(target)
	return err
}

func (s *Service) requestBlob(method, path string, values url.Values) ([]byte, error) {
	resp, err := s.doHTTP(method, path, values)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func (s *Service) doHTTP(method, path string, values url.Values) (*http.Response, error) {
	params := values.Encode()
	u, err := url.Parse(path)
	if err != nil {
		return nil, err
	}
	u.Host = defaultHOST
	u.Scheme = defaultScheme
	if strings.ToLower(method) == "get" {
		u.RawQuery = values.Encode()
	}
	req, err := http.NewRequest(method, u.String(), strings.NewReader(params))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("key", s.apiKey)
	req.Header.Set("sign", s.getSign(params))

	return http.DefaultClient.Do(req)
}

func (s *Service) getSign(params string) string {
	key := []byte(s.secret)
	mac := hmac.New(sha512.New, key)
	mac.Write([]byte(params))
	return fmt.Sprintf("%x", mac.Sum(nil))
}

// GetPairs 获取所有交易对
func (s *Service) GetPairs() ([]string, error) {
	path := "/api2/1/pairs"
	ss := make([]string, 0)
	err := s.requestJSON("GET", path, nil, &ss)
	if err != nil {
		return ss, err
	}
	return ss, nil
}

// MarketInfo 所有市场订单参数 API
func (s *Service) MarketInfo() (*MarketInfoResult, error) {
	path := "/api2/1/marketinfo"
	res := new(MarketInfoResult)
	err := s.requestJSON("GET", path, nil, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type MarketInfoResult struct {
	Result string `json:"result"`
	Pairs  []map[string]struct {
		DecimalPlaces float64 `json:"decimal_places"`
		MinAmount     float64 `json:"min_amount"`
		MinAmountA    float64 `json:"min_amount_a"`
		MinAmountB    float64 `json:"min_amount_b"`
		Fee           float64 `json:"fee"`
		TradeDisabled int     `json:"trade_disabled"`
	} `json:"pairs"`
}

// MarketList 交易市场详细行情 API
func (s *Service) MarketList() (string, error) {
	path := "/api2/1/marketlist"
	bs, err := s.requestBlob("GET", path, nil)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

type Ticker struct {
	Result        string `json:"result"`
	Last          string `json:"last"`
	LowestAsk     string `json:"lowestAsk"`
	HighestBid    string `json:"highestBid"`
	PercentChange string `json:"percentChange"`
	BaseVolume    string `json:"baseVolume"`
	QuoteVolume   string `json:"quoteVolume"`
	High24Hr      string `json:"high24hr"`
	Low24Hr       string `json:"low24hr"`
	Elapsed       string `json:"elapsed"`
}

// Tickers 获取所有交易详情
func (s *Service) Tickers() (map[string]Ticker, error) {
	path := "/api2/1/tickers"
	ts := make(map[string]Ticker, 0)
	err := s.requestJSON("GET", path, nil, &ts)
	if err != nil {
		return nil, err
	}
	return ts, nil
}

// Ticker 获取单项交易详情 gtc_usdt
func (s *Service) Ticker(ticker string) (*Ticker, error) {
	path := "/api2/1/ticker/" + ticker
	fmt.Println(path)
	res := new(Ticker)
	err := s.requestJSON("GET", path, nil, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// OrderBook 交易深度
type OrderBook struct {
	Result  string `json:"result"`
	Elapsed string `json:"elapsed"`
	// 深度里面各种类型混着放 法克
	Asks [][]interface{} `json:"asks"`
	Bids [][]interface{} `json:"bids"`
}

// OrderBooks 返回系统支持的所有交易对的市场深度（委托挂单），其中 asks 是委卖单, bids 是委买单
func (s *Service) OrderBooks() (map[string]OrderBook, error) {
	path := "/api2/1/orderBooks"
	ts := make(map[string]OrderBook, 0)
	err := s.requestJSON("GET", path, nil, &ts)
	if err != nil {
		return nil, err
	}
	return ts, nil
}

// OrderBook 返回当前市场深度
func (s *Service) OrderBook(pair string) (*OrderBook, error) {
	path := "/api2/1/orderBook/" + pair
	fmt.Println(path)
	res := new(OrderBook)
	err := s.requestJSON("GET", path, nil, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type TradeHistoryResult struct {
	Result string `json:"result"`
	Data   []struct {
		TradeID   string `json:"tradeID"`
		Date      string `json:"date"`
		Timestamp string `json:"timestamp"`
		Type      string `json:"type"`
		Rate      string `json:"rate"`
		Amount    string `json:"amount"`
		Total     string `json:"total"`
	} `json:"data"`
	Elapsed string `json:"elapsed"`
}

func (s *Service) TradeHistory(pair string) (*TradeHistoryResult, error) {
	path := "/api2/1/tradeHistory/" + pair
	fmt.Println(path)
	res := new(TradeHistoryResult)
	err := s.requestJSON("GET", path, nil, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type BalanceResult struct {
	Result    string            `json:"result"`
	Available map[string]string `json:"available"`
	Locked    map[string]string `json:"locked"`
}

// Balances 获取帐号资金余额API
func (s *Service) Balances() (*BalanceResult, error) {
	path := "/api2/1/private/balances"
	fmt.Println(path)
	res := new(BalanceResult)
	err := s.requestJSON("POST", path, nil, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}



func (s *Service) NewOne() (string, error) {
	path := "/api2/1/orderBooks"
	bs, err := s.requestBlob("GET", path, nil)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}
