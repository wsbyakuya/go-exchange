package bibox

import "encoding/json"

//Params make request params
type Params struct {
	Cmds   string `json:"cmds"`
	APIKey string `json:"apikey"`
	Sign   string `json:"sign"`
}

//CMD request commands
type CMD struct {
	Cmd   string                 `json:"cmd"`
	Index int                    `json:"index"`
	Body  map[string]interface{} `json:"body"`
}

//Error Error Response
type Error struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

//Results response results
type Results struct {
	Result []json.RawMessage `json:"result"`
	Error  *Error            `json:"error"`
}

//AssetsResult User Assets Result
type AssetsResult struct {
	// Error  Error `json:"error"`
	Result struct {
		TotalBTC   string        `json:"total_btc"`
		TotalCNY   string        `json:"total_cny"`
		TotalUSD   string        `json:"total_usd"`
		AssetsList []SingleAsset `json:"assets_list"`
	} `json:"result"`
	CMD string `json:"cmd"`
}

//DepthResult Market depth Result
type DepthResult struct {
	// Error  Error `json:"error"`
	Result struct {
		Pair       string  `json:"pair"`
		UpdateTime uint64  `json:"update_time"`
		Asks       []Order `json:"asks"`
		Bids       []Order `json:"bids"`
	} `json:"result"`
	CMD string `json:"cmd"`
}

//Order Order
type Order struct {
	Price  string `json:"price"`
	Volume string `json:"volume"`
}

//SingleAsset SingleAsset
type SingleAsset struct {
	CoinSymbol string `json:"coin_symbol"`
	Balance    string `json:"balance"`
	Freeze     string `json:"freeze"`
	BTCValue   string `json:"BTCValue"`
	CNYValue   string `json:"CNYValue"`
	USDValue   string `json:"USDValue"`
}

//CancelTradeResult CancelTradeResult
type CancelTradeResult struct {
	Error  *Error `json:"error"`
	Result string `json:"result"`
	CMD    string `json:"cmd"`
	Index  int    `json:"index"`
}

//TradeResult Single Trade Result
type TradeResult struct {
	Error  *Error `json:"error"`
	Result uint64 `json:"result"`
	CMD    string `json:"cmd"`
	Index  int    `json:"index"`
}

//TradeBody Trade Body
type TradeBody struct {
	Pair        string  `json:"pair"`
	AccountType int     `json:"account_type"`
	OrderType   int     `json:"order_type"`
	OrderSide   int     `json:"order_side"`
	PayBix      int     `json:"pay_bix"`
	Price       float64 `json:"price"`
	Amount      float64 `json:"amount"`
	Money       float64 `json:"money"`
}

//PendingBody Current Pending Body
type PendingBody struct {
	Pair           string `json:"pair"`
	AccountType    int    `json:"account_type"`
	Page           int    `json:"page"`
	Size           int    `json:"size"`
	CoinSymbol     string `json:"coin_symbol"`
	CurrencySymbol string `json:"currency_symbol"`
	OrderSide      int    `json:"order_side"`
	HideCancel     int    `json:"hide_cancel"`
}

//PendingResult PendingResult
type PendingResult struct {
	Result struct {
		Count int           `json:"count"`
		Page  int           `json:"page"`
		Items []PendingItem `json:"items"`
	} `json:"result"`
	CMD string `json:"cmd"`
}

//PendingItem PendingItem
type PendingItem struct {
	ID             int    `json:"id"`
	CreatedAt      uint64 `json:"createdAt"`
	AccountType    int    `json:"account_type"`
	CoinSymbol     string `json:"coin_symbol"`
	CurrencySymbol string `json:"currency_symbol"`
	OrderSide      int    `json:"order_side"`
	OrderType      int    `json:"order_type"`
	Price          string `json:"price"`
	Amount         string `json:"amount"`
	Money          string `json:"money"`
	DealAmount     string `json:"deal_amount"`
	DealPercent    string `json:"deal_percent"`
	UnExcecuted    string `json:"unexecuted"`
	Status         int    `json:"status"`
}
