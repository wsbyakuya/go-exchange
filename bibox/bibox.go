package bibox

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

//BiboxService service for call bibox api
type BiboxService struct {
	URL       string
	APIKey    string
	SecretKey string
}

//NewBiboxService  New A Bibox Service Object
func NewBiboxService(url, apiKey, secret string) (*BiboxService, error) {
	s := &BiboxService{
		URL:       url,
		APIKey:    apiKey,
		SecretKey: secret,
	}
	return s, nil
}

//GetAssets Get User Bibox Assets
func (bs *BiboxService) GetAssets() (*AssetsResult, error) {
	url := bs.URL + "v1/transfer"
	secret := bs.SecretKey
	params := new(Params)
	params.APIKey = bs.APIKey
	cmd := new(CMD)
	cmd.Cmd = "transfer/assets"
	cmd.Body = make(map[string]interface{})
	cmd.Body["select"] = 1
	cmds := make([]*CMD, 0)
	cmds = append(cmds, cmd)
	dataCmds, err := json.Marshal(cmds)
	if err != nil {
		return nil, err
	}
	params.Cmds = string(dataCmds)
	params.Sign = Hmac(secret, params.Cmds)
	dataParams, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(dataParams))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var results Results
	var assetsResult AssetsResult
	err = json.Unmarshal(body, &results)
	if err != nil {
		return nil, err
	}
	if results.Error != nil {
		return nil, errors.New(results.Error.Msg)
	}
	if len(results.Result) != 1 {
		return nil, errors.New("get assets result length invalid")
	}
	err = json.Unmarshal(results.Result[0], &assetsResult)
	if err != nil {
		return nil, err
	}
	return &assetsResult, nil
}

//GetDepth Get Market Depth Data
func (bs *BiboxService) GetDepth(pair string, size int) (*DepthResult, error) {
	url := bs.URL + "v1/mdata"
	secret := bs.SecretKey
	params := new(Params)
	params.APIKey = bs.APIKey
	cmd := new(CMD)
	cmd.Cmd = "api/depth"
	cmd.Body = make(map[string]interface{})
	cmd.Body["pair"] = pair
	cmd.Body["size"] = size
	cmds := make([]*CMD, 0)
	cmds = append(cmds, cmd)
	dataCmds, err := json.Marshal(cmds)
	if err != nil {
		return nil, err
	}
	params.Cmds = string(dataCmds)
	params.Sign = Hmac(secret, params.Cmds)
	dataParams, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(dataParams))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var results Results
	var depthResult DepthResult
	err = json.Unmarshal(body, &results)
	if err != nil {
		return nil, err
	}
	if results.Error != nil {
		return nil, errors.New(results.Error.Msg)
	}
	if len(results.Result) != 1 {
		return nil, errors.New("get depth result length invalid")
	}
	err = json.Unmarshal(results.Result[0], &depthResult)
	if err != nil {
		return nil, err
	}
	return &depthResult, nil
}

//GetBatchDepth Get Market Batch Depth Data
func (bs *BiboxService) GetBatchDepth(pairs []string, size int) ([]*DepthResult, error) {
	url := bs.URL + "v1/mdata"
	secret := bs.SecretKey
	params := new(Params)
	params.APIKey = bs.APIKey
	cmds := make([]*CMD, 0)
	for _, pair := range pairs {
		cmd := new(CMD)
		cmd.Cmd = "api/depth"
		cmd.Body = make(map[string]interface{})
		cmd.Body["pair"] = pair
		cmd.Body["size"] = size
		cmds = append(cmds, cmd)
	}
	dataCmds, err := json.Marshal(cmds)
	if err != nil {
		return nil, err
	}
	params.Cmds = string(dataCmds)
	params.Sign = Hmac(secret, params.Cmds)
	dataParams, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(dataParams))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New(err.Error() + ":" + string(body))
	}
	var results Results
	var depthResults []*DepthResult
	err = json.Unmarshal(body, &results)
	if err != nil {
		return nil, errors.New(err.Error() + ":" + string(body))
	}
	if results.Error != nil {
		return nil, errors.New(results.Error.Msg)
	}
	if len(results.Result) <= 0 {
		return nil, errors.New("get depth result length invalid")
	}
	for _, result := range results.Result {
		var depthResult DepthResult
		err = json.Unmarshal(result, &depthResult)
		if err != nil {
			return nil, err
		}
		depthResults = append(depthResults, &depthResult)
	}
	if len(depthResults) != len(pairs) {
		return nil, errors.New("batch get depth length invalid")
	}
	return depthResults, nil
}

//Trade Trade in Bibox
func (bs *BiboxService) Trade(tradeBody *TradeBody) (*TradeResult, error) {
	url := bs.URL + "v1/orderpending"
	secret := bs.SecretKey
	params := new(Params)
	params.APIKey = bs.APIKey
	cmd := new(CMD)
	cmd.Cmd = "orderpending/trade"
	cmd.Index = 1
	cmd.Body = make(map[string]interface{})
	cmd.Body["pair"] = tradeBody.Pair
	cmd.Body["account_type"] = tradeBody.AccountType
	cmd.Body["order_type"] = tradeBody.OrderType
	cmd.Body["order_side"] = tradeBody.OrderSide
	cmd.Body["pay_bix"] = tradeBody.PayBix
	cmd.Body["price"] = tradeBody.Price
	cmd.Body["amount"] = tradeBody.Amount
	cmd.Body["money"] = tradeBody.Money
	cmds := make([]*CMD, 0)
	cmds = append(cmds, cmd)
	dataCmds, err := json.Marshal(cmds)
	if err != nil {
		return nil, err
	}
	params.Cmds = string(dataCmds)
	params.Sign = Hmac(secret, params.Cmds)
	dataParams, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(dataParams))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var results Results
	var tradeResult TradeResult
	err = json.Unmarshal(body, &results)
	if err != nil {
		return nil, err
	}
	if results.Error != nil {
		return nil, errors.New(results.Error.Msg)
	}
	if len(results.Result) != 1 {
		return nil, errors.New("get trade result length invalid")
	}
	err = json.Unmarshal(results.Result[0], &tradeResult)
	if err != nil {
		return nil, err
	}
	return &tradeResult, nil
}

//BatchTrade Batch Trade
func (bs *BiboxService) BatchTrade(trades []*TradeBody) ([]*TradeResult, error) {
	url := bs.URL + "v1/orderpending"
	secret := bs.SecretKey
	params := new(Params)
	params.APIKey = bs.APIKey
	cmds := make([]*CMD, 0)
	for index, tradeBody := range trades {
		cmd := new(CMD)
		cmd.Cmd = "orderpending/trade"
		cmd.Index = index + 1
		cmd.Body = make(map[string]interface{})
		cmd.Body["pair"] = tradeBody.Pair
		cmd.Body["account_type"] = tradeBody.AccountType
		cmd.Body["order_type"] = tradeBody.OrderType
		cmd.Body["order_side"] = tradeBody.OrderSide
		cmd.Body["pay_bix"] = tradeBody.PayBix
		cmd.Body["price"] = tradeBody.Price
		cmd.Body["amount"] = tradeBody.Amount
		cmd.Body["money"] = tradeBody.Money
		cmds = append(cmds, cmd)
	}
	dataCmds, err := json.Marshal(cmds)
	if err != nil {
		return nil, err
	}
	params.Cmds = string(dataCmds)
	params.Sign = Hmac(secret, params.Cmds)
	dataParams, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(dataParams))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var results Results
	var tradeResults []*TradeResult
	err = json.Unmarshal(body, &results)
	if err != nil {
		return nil, err
	}
	if len(results.Result) != len(trades) {
		return nil, errors.New("get trade result length invalid")
	}
	for _, result := range results.Result {
		var oneResult TradeResult
		err = json.Unmarshal(result, &oneResult)
		if err != nil {
			return nil, err
		}
		tradeResults = append(tradeResults, &oneResult)
	}
	if results.Error != nil {
		return tradeResults, errors.New(results.Error.Msg)
	}
	return tradeResults, nil
}

//CancelTrade Cancel Pending Trade
func (bs *BiboxService) CancelTrade(id uint64) (*CancelTradeResult, error) {
	url := bs.URL + "v1/orderpending"
	secret := bs.SecretKey
	params := new(Params)
	params.APIKey = bs.APIKey
	cmd := new(CMD)
	cmd.Cmd = "orderpending/cancelTrade"
	cmd.Index = 1
	cmd.Body = make(map[string]interface{})
	cmd.Body["orders_id"] = id
	cmds := make([]*CMD, 0)
	cmds = append(cmds, cmd)
	dataCmds, err := json.Marshal(cmds)
	if err != nil {
		return nil, err
	}
	params.Cmds = string(dataCmds)
	params.Sign = Hmac(secret, params.Cmds)
	dataParams, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(dataParams))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var results Results
	var cancelResult CancelTradeResult
	err = json.Unmarshal(body, &results)
	if err != nil {
		return nil, err
	}
	if results.Error != nil {
		return nil, errors.New(results.Error.Msg)
	}
	if len(results.Result) != 1 {
		return nil, errors.New("get cancel trade result length invalid")
	}
	err = json.Unmarshal(results.Result[0], &cancelResult)
	if err != nil {
		return nil, err
	}
	return &cancelResult, nil
}

//BatchCancelTrade Batch Cancel Pending Trade
func (bs *BiboxService) BatchCancelTrade(ids []uint64) ([]*CancelTradeResult, error) {
	url := bs.URL + "v1/orderpending"
	secret := bs.SecretKey
	params := new(Params)
	cmds := make([]*CMD, 0)
	params.APIKey = bs.APIKey
	for index, id := range ids {
		cmd := new(CMD)
		cmd.Cmd = "orderpending/cancelTrade"
		cmd.Index = index + 1
		cmd.Body = make(map[string]interface{})
		cmd.Body["orders_id"] = id
		cmds = append(cmds, cmd)
	}
	dataCmds, err := json.Marshal(cmds)
	if err != nil {
		return nil, err
	}
	params.Cmds = string(dataCmds)
	params.Sign = Hmac(secret, params.Cmds)
	dataParams, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(dataParams))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var results Results
	err = json.Unmarshal(body, &results)
	if err != nil {
		return nil, err
	}
	if len(results.Result) != len(ids) {
		return nil, errors.New("get batch cancel trade result length invalid")
	}
	var returnResults []*CancelTradeResult
	for _, result := range results.Result {
		var oneResult CancelTradeResult
		err = json.Unmarshal(result, &oneResult)
		if err != nil {
			return nil, err
		}
		returnResults = append(returnResults, &oneResult)
	}
	if results.Error != nil {
		return returnResults, errors.New(results.Error.Msg)
	}
	return returnResults, nil
}

//CurrentPending Get Current Pending Orders
func (bs *BiboxService) CurrentPending(pendingBody *PendingBody) (*PendingResult, error) {
	url := bs.URL + "v1/orderpending"
	secret := bs.SecretKey
	params := new(Params)
	params.APIKey = bs.APIKey
	cmd := new(CMD)
	cmd.Cmd = "orderpending/orderPendingList"
	cmd.Index = 1
	cmd.Body = make(map[string]interface{})
	if pendingBody.Pair != "" {
		cmd.Body["pair"] = pendingBody.Pair
	}
	if pendingBody.AccountType != -1 {
		cmd.Body["account_type"] = pendingBody.AccountType
	}
	cmd.Body["page"] = pendingBody.Page
	cmd.Body["size"] = pendingBody.Size
	if pendingBody.CoinSymbol != "" {
		cmd.Body["coin_symbol"] = pendingBody.CoinSymbol
	}
	if pendingBody.CurrencySymbol != "" {
		cmd.Body["currency_symbol"] = pendingBody.CurrencySymbol
	}
	if pendingBody.OrderSide != 0 {
		cmd.Body["order_side"] = pendingBody.OrderSide
	}
	cmds := make([]*CMD, 0)
	cmds = append(cmds, cmd)
	dataCmds, err := json.Marshal(cmds)
	if err != nil {
		return nil, err
	}
	params.Cmds = string(dataCmds)
	params.Sign = Hmac(secret, params.Cmds)
	dataParams, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(dataParams))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var results Results
	err = json.Unmarshal(body, &results)
	if err != nil {
		return nil, err
	}
	if results.Error != nil {
		return nil, errors.New(results.Error.Msg)
	}
	if len(results.Result) != 1 {
		return nil, errors.New("pending result length invalid")
	}
	var pendingResult PendingResult
	err = json.Unmarshal(results.Result[0], &pendingResult)
	if err != nil {
		return nil, errors.New("parse pending result error")
	}
	return &pendingResult, nil
}

//HistoryPending Get History Pending Orders
func (bs *BiboxService) HistoryPending(pendingBody *PendingBody) (*PendingResult, error) {
	url := bs.URL + "v1/orderpending"
	secret := bs.SecretKey
	params := new(Params)
	params.APIKey = bs.APIKey
	cmd := new(CMD)
	cmd.Cmd = "orderpending/pendingHistoryList"
	cmd.Index = 1
	cmd.Body = make(map[string]interface{})
	if pendingBody.Pair != "" {
		cmd.Body["pair"] = pendingBody.Pair
	}
	if pendingBody.AccountType != -1 {
		cmd.Body["account_type"] = pendingBody.AccountType
	}
	cmd.Body["page"] = pendingBody.Page
	cmd.Body["size"] = pendingBody.Size
	if pendingBody.CoinSymbol != "" {
		cmd.Body["coin_symbol"] = pendingBody.CoinSymbol
	}
	if pendingBody.CurrencySymbol != "" {
		cmd.Body["currency_symbol"] = pendingBody.CurrencySymbol
	}
	if pendingBody.OrderSide != 0 {
		cmd.Body["order_side"] = pendingBody.OrderSide
	}
	if pendingBody.HideCancel != -1 {
		cmd.Body["hide_cancel"] = pendingBody.HideCancel
	}
	cmds := make([]*CMD, 0)
	cmds = append(cmds, cmd)
	dataCmds, err := json.Marshal(cmds)
	if err != nil {
		return nil, err
	}
	params.Cmds = string(dataCmds)
	params.Sign = Hmac(secret, params.Cmds)
	dataParams, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(dataParams))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var results Results
	err = json.Unmarshal(body, &results)
	if err != nil {
		return nil, err
	}
	if results.Error != nil {
		return nil, errors.New(results.Error.Msg)
	}
	if len(results.Result) != 1 {
		return nil, errors.New("pending result length invalid")
	}
	var pendingResult PendingResult
	err = json.Unmarshal(results.Result[0], &pendingResult)
	if err != nil {
		return nil, errors.New("parse pending result error")
	}
	return &pendingResult, nil
}
