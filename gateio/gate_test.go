package gateio

import (
	"io/ioutil"
	"testing"
)

const testKey = ""                                // gate.io api key
const testSecret = "" // gate.io api secret


func TestServiceDoHTTP(t *testing.T) {
	s := NewService("", "")
	resp, err := s.doHTTP("GET", "http://data.gateio.io/api2/1/orderBooks", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(bs))
}

func TestService_GetPairs(t *testing.T) {
	s := NewService("", "")
	res, err := s.GetPairs()
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range res {
		t.Log(v)
	}
}

func TestService_MarketInfo(t *testing.T) {
	s := NewService("", "")
	res, err := s.MarketInfo()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(res)
}
func TestService_MarketList(t *testing.T) {
	s := NewService("", "")
	res, err := s.OrderBooks()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(res)
}

func TestService_Ticker(t *testing.T) {
	s := NewService("", "")
	res, err := s.TradeHistory("stx_usdt")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(res)
}

func TestService_NewOne(t *testing.T) {
	s := NewService(testKey, testSecret)
	res, err := s.NewOne()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(res)
}

func TestService_Balances(t *testing.T) {
	s := NewService(testKey, testSecret)
	res, err := s.Balances()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(res)
}