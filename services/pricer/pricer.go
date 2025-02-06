package pricer

import (
	"encoding/json"
	"fmt"
	"log"
	"project/helpers"
	"strconv"
)

type BinanceTicker struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

type BitGetTicker struct {
	Code        string `json:"code"`
	Msg         string `json:"msg"`
	RequestTime int64  `json:"requestTime"`
	Data        []struct {
		LastPr string `json:"lastPr"`
	} `json:"data"`
}

func GetLastPriceBySymbol(symbol string, hostname string) float64 {
	if hostname == "binance" {
		return getLastTickerFromBinance(symbol)
	} else if hostname == "bitget" {
		return getLastTickerFromBitget(symbol)
	} else {
		return 0
	}
}

func getLastTickerFromBinance(symbol string) float64 {
	params := map[string]string{
		"symbol": symbol,
	}
	response, err := helpers.SendRequest(params, "GET", "https://api.binance.com/api/v3/ticker/price", "binance")
	if err != nil {
		fmt.Printf("Ошибка https://api.binance.com/api/v3/ticker/price: %v\n", err)
		return 0
	}

	var ticker BinanceTicker
	err = json.Unmarshal([]byte(response), &ticker)
	if err != nil {
		log.Fatalf("Ошибка парсинга binance ticker JSON: %v", err)
		return 0
	}

	price, err := strconv.ParseFloat(ticker.Price, 64)
	if err == nil {
		return price
	} else {
		return 0
	}
}

func getLastTickerFromBitget(symbol string) float64 {
	params := map[string]string{
		"symbol": symbol,
	}
	response, err := helpers.SendRequest(params, "GET", "https://api.bitget.com/api/v2/spot/market/tickers", "bitget")
	if err != nil {
		fmt.Printf("Ошибка https://api.bitget.com/api/v2/spot/market/tickers: %v\n", err)
		return 0
	}

	var ticker BitGetTicker
	err = json.Unmarshal([]byte(response), &ticker)
	if err != nil {
		log.Fatalf("Ошибка парсинга bitget ticker JSON: %v", err)
		return 0
	}

	if len(ticker.Data) > 0 {
		price, err := strconv.ParseFloat(ticker.Data[0].LastPr, 64)
		if err == nil {
			return price
		}
	}

	return 0
}
