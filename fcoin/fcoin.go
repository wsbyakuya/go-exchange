package fcoin

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// FcoinService service for call fcoin api
type FcoinService struct {
	URL       string
	APIKey    string
	SecretKey string
}

// NewFcoinService  New A fcoin Service Object
func NewFcoinService(url, apiKey, secret string) (*FcoinService, error) {
	s := &FcoinService{
		URL:       url,
		APIKey:    apiKey,
		SecretKey: secret,
	}
	return s, nil
}

// authorization 授权请求
func (fs *FcoinService) authorization(method, path string, params, body url.Values) (json.RawMessage, error) {
	method = strings.ToUpper(method)
	sURI := sortedURI(fs.URL+path, params)
	ts := strconv.Itoa(int(time.Now().Unix()) * 1e3)
	sBody := sortedBody(body)

	src := []byte(method + sURI + ts + sBody)
	buf := make([]byte, base64.StdEncoding.EncodedLen(len(src)))
	base64.StdEncoding.Encode(buf, src)
	mac := hmac.New(sha1.New, []byte(fs.SecretKey))
	mac.Write(buf)
	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	var reader io.Reader
	if body != nil {
		reader = strings.NewReader(urlValuesToJSON(body))
	}

	req, err := http.NewRequest(method, sURI, reader)
	if err != nil {
		return nil, err
	}
	if method == "POST" || method == "PUT" {
		req.Header.Add("Content-Type", "application/json")
	}
	req.Header.Add("FC-ACCESS-KEY", fs.APIKey)
	req.Header.Add("FC-ACCESS-SIGNATURE", signature)
	req.Header.Add("FC-ACCESS-TIMESTAMP", ts)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	res := new(Result)
	err = json.NewDecoder(resp.Body).Decode(res)
	if err != nil {
		return nil, err
	}

	if res.Status != 0 {
		if res.Status == 2000 {
			return nil, ErrAccountError
		} else {
			return nil, fmt.Errorf("response status is %d msg %s", res.Status, res.Msg)
		}
	}

	return res.Data, nil
}

func (fs *FcoinService) public(method, path string, params, body url.Values) (json.RawMessage, error) {
	method = strings.ToUpper(method)
	sURI := sortedURI(fs.URL+path, params)

	var reader io.Reader
	if body != nil {
		reader = strings.NewReader(urlValuesToJSON(body))
	}

	req, err := http.NewRequest(method, sURI, reader)
	if err != nil {
		return nil, err
	}
	if method == "POST" || method == "PUT" {
		req.Header.Add("Content-Type", "application/json")
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	res := new(Result)
	err = json.NewDecoder(resp.Body).Decode(res)
	if err != nil {
		return nil, err
	}

	if res.Status != 0 {
		if res.Status == 2000 {
			return nil, ErrAccountError
		} else {
			return nil, fmt.Errorf("response status is %d msg %s", res.Status, res.Msg)
		}
	}

	return res.Data, nil
}

func urlValuesToJSON(values url.Values) string {
	m := make(map[string]string)
	for key := range values {
		m[key] = values.Get(key)
	}
	bs, _ := json.Marshal(&m)
	return string(bs)
}

// GetServerTime 查询服务器时间
func (fs *FcoinService) GetServerTime() (time.Time, error) {
	path := `/v2/public/server-time`
	data, err := fs.public("GET", path, nil, nil)
	if err != nil {
		return time.Time{}, err
	}
	var ts int64
	err = json.Unmarshal(data, &ts)
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(int64(ts/1e3), int64(ts%1e3*1e6)), nil
}

// GetCurrencies 查询可用币种
func (fs *FcoinService) GetCurrencies() ([]string, error) {
	path := `/v2/public/currencies`
	data, err := fs.public("GET", path, nil, nil)
	if err != nil {
		return nil, err
	}

	res := make([]string, 0)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return res, err
	}

	return res, nil
}

// GetSymbols 查询可用交易对
func (fs *FcoinService) GetSymbols() ([]Symbol, error) {
	path := `/v2/public/symbols`
	data, err := fs.public("GET", path, nil, nil)
	if err != nil {
		return nil, err
	}

	res := make([]Symbol, 0)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return res, err
	}

	return res, nil
}

// GetMarketTicker 获取 ticker 数据
/**
  "最新成交价",
  "最近一笔成交的成交量",
  "最大买一价",
  "最大买一量",
  "最小卖一价",
  "最小卖一量",
  "24小时前成交价",
  "24小时内最高价",
  "24小时内最低价",
  "24小时内基准货币成交量, 如 btcusdt 中 btc 的量",
  "24小时内计价货币成交量, 如 btcusdt 中 usdt 的量"
*/
func (fs *FcoinService) GetMarketTicker(symbol string) (*MarketTicker, error) {
	path := `/v2/market/ticker/` + symbol
	data, err := fs.public("GET", path, nil, nil)
	if err != nil {
		return nil, err
	}

	res := new(MarketTicker)
	err = json.Unmarshal(data, res)
	if err != nil {
		return res, err
	}

	return res, nil
}

// GetMarketDepth 获取最新的深度明细
/**
level 包括
L20	20 档行情深度.
L100	100 档行情深度.
full	全量的行情深度, 不做时间保证和推送保证.
*/
func (fs *FcoinService) GetMarketDepth(level, symbol string) (*MarketDepth, error) {
	path := `/v2/market/depth/` + level + `/` + symbol
	data, err := fs.authorization("GET", path, nil, nil)
	if err != nil {
		return nil, err
	}

	res := new(MarketDepth)
	err = json.Unmarshal(data, res)
	if err != nil {
		return res, err
	}

	return res, nil
}

// GetMarketTrades 获取最新的成交明细
// before		查询某个 id 之前的 Trade
// limit		默认为 20 条
func (fs *FcoinService) GetMarketTrades(symbol, before string, limit int) ([]MarketTrade, error) {
	values := url.Values{}
	values.Add("before", before)
	if limit == 0 {
		limit = 20
	}
	values.Add("limit", strconv.Itoa(limit))

	path := `/v2/market/trades/` + symbol
	data, err := fs.public("GET", path, values, nil)
	if err != nil {
		return nil, err
	}

	res := make([]MarketTrade, 0)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return res, err
	}

	return res, nil
}

// GetMarketTrades 获取 Candle 信息
// before		查询某个 id 之前的 Trade
// limit		默认为 20 条
/**
resolution 包含的种类

类型	说明
M1	1 分钟
M3	3 分钟
M5	5 分钟
M15	15 分钟
M30	30 分钟
H1	1 小时
H4	4 小时
H6	6 小时
D1	1 日
W1	1 周
MN	1 月
*/
func (fs *FcoinService) GetMarketCandle(resolution, symbol, before string, limit int) ([]MarketCandle, error) {
	values := url.Values{}
	values.Add("before", before)
	if limit == 0 {
		limit = 20
	}
	values.Add("limit", strconv.Itoa(limit))

	path := `/v2/market/candles/` + resolution + `/` + symbol
	data, err := fs.public("GET", path, values, nil)
	if err != nil {
		return nil, err
	}

	res := make([]MarketCandle, 0)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return res, err
	}

	return res, nil
}

// GetAccountBalance 查询账户资产
func (fs *FcoinService) GetAccountBalance() ([]AccountBalance, error) {
	path := `/v2/accounts/balance`
	data, err := fs.authorization("GET", path, nil, nil)
	if err != nil {
		return nil, err
	}

	res := make([]AccountBalance, 0)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return res, err
	}

	return res, nil
}

/**
订单模型由以下属性构成：

属性	类型	含义解释
id	String	订单 ID
symbol	String	交易对
side	String	交易方向（buy, sell）
type	String	订单类型（limit，market）
price	String	下单价格
amount	String	下单数量
state	String	订单状态
executed_value	String	已成交
filled_amount	String	成交量
fill_fees	String	手续费
created_at	Long	创建时间
source	String	来源
*/

/**
订单状态说明：

属性	含义解释
submitted	已提交
partial_filled	部分成交
partial_canceled	部分成交已撤销
filled	完全成交
canceled	已撤销
pending_cancel	撤销已提交
*/

// CreateOrder 创建新的订单 返回订单ID
func (fs *FcoinService) CreateOrder(symbol, side, orderType, price, amount string) (string, error) {
	values := url.Values{}
	values.Add("symbol", symbol)  // 交易对
	values.Add("side", side)      // 交易方向
	values.Add("type", orderType) // 订单类型
	values.Add("price", price)    // 价格
	values.Add("amount", amount)  // 下单量

	path := `/v2/orders`
	data, err := fs.authorization("POST", path, nil, values)
	if err != nil {
		return "", err
	}

	ordID := ""
	err = json.Unmarshal(data, &ordID)
	if err != nil {
		return ordID, err
	}

	return ordID, nil
}

// GetOrders 查询订单列表
func (fs *FcoinService) GetOrders(symbol, states, before, after, limit string) ([]OrderInformation, error) {
	values := url.Values{}
	values.Add("symbol", symbol) // 交易对
	values.Add("states", states) // 订单状态
	values.Add("before", before) // 查询某个页码之前的订单
	values.Add("after", after)   // 查询某个页码之后的订单
	values.Add("limit", limit)   // 每页的订单数量，默认为 20 条

	path := `/v2/orders`
	data, err := fs.authorization("GET", path, values, nil)
	if err != nil {
		return nil, err
	}

	res := make([]OrderInformation, 0)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return res, err
	}

	return res, nil
}

// GetOrderByID 返回指定的订单详情
func (fs *FcoinService) GetOrderByID(orderID string) (*OrderInformation, error) {
	path := `/v2/orders/` + orderID
	data, err := fs.authorization("GET", path, nil, nil)
	if err != nil {
		return nil, err
	}

	res := new(OrderInformation)
	err = json.Unmarshal(data, res)
	if err != nil {
		return res, err
	}

	return res, nil
}

// CancelOrder 申请撤销订单
func (fs *FcoinService) CancelOrder(orderID string) (bool, error) {
	path := `/v2/orders/` + orderID + `/submit-cancel`
	data, err := fs.authorization("POST", path, nil, nil)
	if err != nil {
		return false, err
	}
	fmt.Println(string(data))
	res := false
	err = json.Unmarshal(data, &res)
	if err != nil {
		return res, err
	}

	return res, nil
}

// OrderMatchResult 查询指定订单的成交记录
func (fs *FcoinService) OrderMatchResult(orderID string) ([]OrderMatchResult, error) {
	path := `/v2/orders/` + orderID + `/match-results`
	data, err := fs.authorization("POST", path, nil, nil)
	if err != nil {
		return nil, err
	}

	res := make([]OrderMatchResult, 0)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (fs *FcoinService) GetCurrentMarketPrice(symbol string, priceDecimal int) (float64, error) {
	depth, err := fs.GetMarketDepth("L20", symbol)
	if err != nil {
		return 0.0, err
	}
	if len(depth.Bids) <= 1 || len(depth.Asks) <= 1 { // bids和asks成对出现 长度小于1越界
		return 0.0, fmt.Errorf("return wrong error")
	}

	bid := depth.Bids[0] // 卖单最小价 当前可买入的最低价格
	ask := depth.Asks[0] // 买单最大价 当前可卖出的最高价格
	//return (bid + ask) / 2, nil
	price := bid + (ask-bid)*1/2
	if priceDecimal > 0 {
		if priceStr := strconv.FormatFloat(price, 'f', priceDecimal, 64); priceStr == strconv.FormatFloat(bid, 'f', priceDecimal, 64) ||
			priceStr == strconv.FormatFloat(ask, 'f', priceDecimal, 64) {
			return 0.0, errors.New("the price range is too small")
		}
	}

	return price, nil
}

func (fs *FcoinService) GetAvailableAmount(coin string) (float64, error) {
	accountBalance, err := fs.GetAccountBalance()
	if err != nil {
		return 0.0, fmt.Errorf("get account balance error %v", err)
	}
	amount := 0.0
	for _, v := range accountBalance {
		if v.Currency == coin {
			if f, err := strconv.ParseFloat(v.Available, 64); err == nil {
				amount = f
			}
		}
	}
	return amount, nil
}

func (fs *FcoinService) HasEnoughAssets(coin string, amount float64) (bool, error) {
	var enoughFlag bool
	accountBalance, err := fs.GetAccountBalance()
	if err != nil {
		return false, fmt.Errorf("get account balance error %v", err)
	}
	for _, v := range accountBalance {
		if v.Currency == coin {
			if f, err := strconv.ParseFloat(v.Available, 64); err != nil || f < amount {
				return false, nil
			} else if f >= amount {
				enoughFlag = true
			}
		}
	}
	return enoughFlag, nil
}

func (fs *FcoinService) IsOrderFinished(orderID string) (bool, error) {
	order, err := fs.GetOrderByID(orderID)
	if err != nil {
		return false, err
	}

	if order.State == "filled" || order.State == "canceled" || order.State == "partial_canceled" { // 卖出完全交易
		return true, nil
	}
	return false, nil
}

func GetCurrentCoinType(symbol string) (string, string) {
	switch symbol {
	case "ethusdt":
		return "eth", "usdt"
	case "fteth":
		return "ft", "eth"
	case "zipeth":
		return "zip", "eth"
	case "zileth":
		return "zil", "eth"
	}
	return "", ""
}
