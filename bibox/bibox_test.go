package bibox

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAssets(t *testing.T) {
	url := "https://api.bibox.com/"
	s, err := NewBiboxService(url, "", "")
	if err != nil {
		t.Error(err)
	}
	assetsResult, err := s.GetAssets()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(assetsResult)
}

func TestDepth(t *testing.T) {
	url := "https://api.bibox.com/"
	s, err := NewBiboxService(url, "", "")
	if err != nil {
		t.Error(err)
	}
	for {
		// ts1 := time.Now()
		results, err := s.GetBatchDepth([]string{"PAI_ETH", "PAI_BTC", "ETH_BTC"}, 1)
		if err != nil {
			t.Error(err)
		}
		// ts2 := time.Now()
		// fmt.Println(ts2.UnixNano() - ts1.UnixNano())
		var price float64
		price = 1
		for _, result := range results {
			ask := result.Result.Asks[0]
			bid := result.Result.Bids[0]
			if result.Result.Pair == "PAI_ETH" {
				// bidFloat, _ := strconv.ParseFloat(bid.Price, 64)
				askFloat, _ := strconv.ParseFloat(ask.Price, 64)
				price = price / askFloat
			}
			if result.Result.Pair == "PAI_BTC" {
				bidFloat, _ := strconv.ParseFloat(bid.Price, 64)
				// askFloat, _ := strconv.ParseFloat(ask.Price, 64)
				price = price * bidFloat
			}
			if result.Result.Pair == "ETH_BTC" {
				// bidFloat, _ := strconv.ParseFloat(bid.Price, 64)
				askFloat, _ := strconv.ParseFloat(ask.Price, 64)
				price = price / askFloat
			}
		}
		if price > 1 {
			fmt.Println(price)
		}
		time.Sleep(100 * time.Millisecond)
	}
	// for _, result := range results {
	// 	fmt.Println(result)
	// }
	// _, err = s.GetDepth("BIX_BTC", 2)
	// if err != nil {
	// 	t.Error(err)
	// }
	// _, err = s.GetDepth("ETH_BTC", 2)
	// if err != nil {
	// 	t.Error(err)
	// }

}

func TestUser(t *testing.T) {
	url := "https://api.bibox.com/v1/user/"
	url = "https://api.bibox.com/v1/transfer"
	secret := ""
	params := new(Params)
	params.APIKey = ""
	cmd := new(CMD)
	cmd.Cmd = "transfer/assets"
	cmd.Body = make(map[string]interface{})
	cmds := make([]*CMD, 0)
	cmds = append(cmds, cmd)
	dataCmds, err := json.Marshal(cmds)
	if err != nil {
		t.Error(err)
	}
	params.Cmds = string(dataCmds)
	params.Sign = Hmac(secret, params.Cmds)
	dataParams, err := json.Marshal(params)
	if err != nil {
		t.Error(err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(dataParams))
	if err != nil {
		t.Error(err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(body))
}

func TestSimpleRequest(t *testing.T) {
	url := "https://api.bibox.com/v1/mdata?cmd=pairList"
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}

	resp, err := client.Get(url)
	if err != nil {
		t.Error(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(body))
}

func TestTrade(t *testing.T) {
	url := "https://api.bibox.com/"
	s, err := NewBiboxService(url, "", "")
	if err != nil {
		t.Error(err)
	}
	body := &TradeBody{
		Pair:        "BIX_ETH",
		AccountType: 0,
		OrderType:   2,
		OrderSide:   1,
		PayBix:      0,
		Price:       0.0000008647,
		Amount:      1,
		Money:       0.00008647,
	}
	result, err := s.Trade(body)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(result.Result)
}

func TestBatchTrade(t *testing.T) {
	url := "https://api.bibox.com/"
	s, err := NewBiboxService(url, "", "")
	if err != nil {
		t.Error(err)
	}
	body0 := &TradeBody{
		Pair:        "BIX_ETH",
		AccountType: 0,
		OrderType:   2,
		OrderSide:   1,
		PayBix:      0,
		Price:       0.0000008647,
		Amount:      1,
		Money:       0.0000008647,
	}
	body1 := &TradeBody{
		Pair:        "BIX_BTC",
		AccountType: 0,
		OrderType:   2,
		OrderSide:   1,
		PayBix:      0,
		Price:       0.0000008647,
		Amount:      1,
		Money:       0.0000008647,
	}
	trades := make([]*TradeBody, 0)
	trades = append(trades, body0, body1)
	tradeResults, err := s.BatchTrade(trades)
	if err != nil {
		fmt.Println(err)
	}
	for _, trade := range tradeResults {
		if trade.Error != nil {
			fmt.Println(trade.Error.Msg)
		}
	}
}

func TestCancelTrade(t *testing.T) {
	url := "https://api.bibox.com/"
	s, err := NewBiboxService(url, "", "")
	if err != nil {
		t.Error(err)
	}
	_, err = s.CancelTrade(612216386)
	if err != nil {
		t.Error(err)
	}
}

func TestBatchCancelTrade(t *testing.T) {
	url := "https://api.bibox.com/"
	s, err := NewBiboxService(url, "", "")
	if err != nil {
		t.Error(err)
	}
	results, err := s.BatchCancelTrade([]uint64{612216386, 612212608})
	for _, result := range results {
		fmt.Println(result.Error.Msg)
	}
	if err != nil {
		t.Error(err)
	}
}

func TestBuyPai(t *testing.T) {
	url := "https://api.bibox.com/"
	s, err := NewBiboxService(url, "", "")
	if err != nil {
		t.Error(err)
	}
	result, err := s.GetDepth("PAI_ETH", 1)
	if err != nil {
		t.Error(err)
	}
	// askPrice, err := strconv.ParseFloat(result.Result.Asks[0].Price, 64)
	// if err != nil {
	// 	t.Error(err)
	// }
	bidPrice, err := strconv.ParseFloat(result.Result.Bids[0].Price, 64)
	if err != nil {
		t.Error(err)
	}
	askVolume, err := strconv.ParseFloat(result.Result.Asks[0].Volume, 64)
	if err != nil {
		t.Error(err)
	}
	if askVolume < 1 {
		t.Error("pai volume < 1")
	}
	// fmt.Println(askVolume)
	//buy
	body := &TradeBody{
		Pair:        "PAI_ETH",
		AccountType: 0,
		OrderType:   2,
		OrderSide:   2,
		PayBix:      0,
		Price:       bidPrice,
		Amount:      1.00011,
		Money:       0,
	}
	_, err = s.Trade(body)
	if err != nil {
		t.Error(err)
	}
}

func TestBatchBuyPai(t *testing.T) {
	url := "https://api.bibox.com/"
	s, err := NewBiboxService(url, "", "")
	if err != nil {
		t.Error(err)
	}
	coins := []string{"PAI_ETH", "PAI_BTC", "ETH_BTC"}
	body1 := &TradeBody{
		Pair:        coins[0],
		AccountType: 0,
		OrderType:   2,
		OrderSide:   1,
		PayBix:      0,
		Price:       0.00029543,
		Amount:      25.3904,
		Money:       0,
	}
	body2 := &TradeBody{
		Pair:        coins[1],
		AccountType: 0,
		OrderType:   2,
		OrderSide:   2,
		PayBix:      0,
		Price:       0.00002336,
		Amount:      25.3904,
		Money:       0,
	}
	body3 := &TradeBody{
		Pair:        coins[2],
		AccountType: 0,
		OrderType:   2,
		OrderSide:   1,
		PayBix:      0,
		Price:       0.07897439,
		Amount:      0.0075102795222603165,
		Money:       0,
	}
	trades := make([]*TradeBody, 0)
	trades = append(trades, body1, body2, body3)
	tradeResults, err := s.BatchTrade(trades)
	if err != nil {
		fmt.Println(err)
	}
	for _, trade := range tradeResults {
		if trade.Error != nil {
			fmt.Println(trade.Error.Msg)
		}
	}
}

func TestPending(t *testing.T) {
	url := "https://api.bibox.com/"
	s, err := NewBiboxService(url, "", "")
	if err != nil {
		t.Error(err)
	}
	body := &PendingBody{
		Page: 1,
		Size: 10,
	}
	result, err := s.CurrentPending(body)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(result.Result.Items[0].Money)
}

func TestPendingHistory(t *testing.T) {
	url := "https://api.bibox.com/"
	s, err := NewBiboxService(url, "", "")
	if err != nil {
		t.Error(err)
	}
	body := &PendingBody{
		Page: 1,
		Size: 10,
	}
	result, err := s.HistoryPending(body)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(len(result.Result.Items))
}

func TestSign(t *testing.T) {
	secret := ""
	cmds := `[{"cmd":"user/userInfo","body":{}}]`
	data := Hmac(secret, cmds)
	assert.Equal(t, data, "", "md5 hmac not equal")
}
