package fcoin

import (
	"encoding/json"
	"errors"
)

var (
	ErrAccountError = errors.New("account error")
)

type Result struct {
	Status int             `json:"status"`
	Msg    string          `json:"msg"`
	Data   json.RawMessage `json:"data"`
}

type Symbol struct {
	Name          string `json:"name"`
	BaseCurrency  string `json:"base_currency"`
	QuoteCurrency string `json:"quote_currency"`
	PriceDecimal  int    `json:"price_decimal"`
	AmountDecimal int    `json:"amount_decimal"`
}

type MarketTicker struct {
	Type   string    `json:"type"`
	Seq    int       `json:"seq"`
	Ticker []float64 `json:"ticker"`
}

type MarketDepth struct {
	Type string    `json:"type"`
	Ts   int64     `json:"ts"`
	Seq  int       `json:"seq"`
	Bids []float64 `json:"bids"`
	Asks []float64 `json:"asks"`
}

type MarketTrade struct {
	Amount float64 `json:"amount"`
	Ts     int64   `json:"ts"`
	ID     int     `json:"id"`
	Side   string  `json:"side"`
	Price  float64 `json:"price"`
}

type MarketCandle struct {
	Type     string  `json:"type"`
	ID       int     `json:"id"`
	Seq      int     `json:"seq"`
	Open     float64 `json:"open"`
	Close    float64 `json:"close"`
	High     float64 `json:"high"`
	Low      float64 `json:"low"`
	Count    int     `json:"count"`
	BaseVol  float64 `json:"base_vol"`
	QuoteVol float64 `json:"quote_vol"`
}

type AccountBalance struct {
	Currency  string `json:"currency"`
	Available string `json:"available"`
	Frozen    string `json:"frozen"`
	Balance   string `json:"balance"`
}

type OrderInformation struct {
	ID            string `json:"id"`
	Symbol        string `json:"symbol"`
	Type          string `json:"type"`
	Side          string `json:"side"`
	Price         string `json:"price"`
	Amount        string `json:"amount"`
	State         string `json:"state"`
	ExecutedValue string `json:"executed_value"`
	FillFees      string `json:"fill_fees"`
	FilledAmount  string `json:"filled_amount"`
	CreatedAt     int    `json:"created_at"`
	Source        string `json:"source"`
}

type OrderMatchResult struct {
	Price        string `json:"price"`
	FillFees     string `json:"fill_fees"`
	FilledAmount string `json:"filled_amount"`
	Side         string `json:"side"`
	Type         string `json:"type"`
	CreatedAt    int    `json:"created_at"`
}
