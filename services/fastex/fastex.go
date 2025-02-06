package fastex

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"project/config"
	"project/helpers"
	"strconv"
)

type Order struct {
	ID    int     `json:"id"`
	Price float64 `json:"price"`
	Type  string  `json:"type"` // buy или sell
}

var OrderList = make(map[string][]Order)

type APIResponse struct {
	Errors   bool `json:"errors"`
	Response struct {
		Entities []struct {
			OrderID       int    `json:"order_id"`
			Price         string `json:"price"`
			Type          string `json:"type"`
			Volume        string `json:"volume"`
			TradingSymbol string `json:"trading_symbol"`
			Created       int64  `json:"created"`
		} `json:"entities"`
	} `json:"response"`
	Pagination struct {
		ItemsPerPage int `json:"items_per_page"`
		TotalItems   int `json:"total_items"`
		CurrentPage  int `json:"current_page"`
		TotalPages   int `json:"total_pages"`
	} `json:"pagination"`
}

func GetAllOrders() {
	for _, symbol := range config.Symbols {
		params := map[string]string{
			"symbol": symbol.FastexName,
		}
		response, err := helpers.SendRequest(params, "POST", os.Getenv("FASTEX_DOMAIN_URL")+"/api/v1/order/list", "fastex")
		if err != nil {
			fmt.Printf("Ошибка: %v\n", err)
			return
		}

		var apiResponse APIResponse
		parseJsonErr := json.Unmarshal([]byte(response), &apiResponse)
		if parseJsonErr != nil {
			fmt.Printf("Ошибка при парсинге JSON: %v\n", err)
			return
		}

		OrderList[symbol.FastexName] = nil

		for _, entity := range apiResponse.Response.Entities {
			price, err := strconv.ParseFloat(entity.Price, 64)
			if err != nil {
				fmt.Printf("Ошибка при парсинге цены ордера %d: %v\n", entity.OrderID, err)
				continue
			}
			order := Order{
				ID:    entity.OrderID,
				Price: price,
				Type:  entity.Type,
			}
			OrderList[symbol.FastexName] = append(OrderList[symbol.FastexName], order)
		}
	}
}

func CloseOrder(orderId int) string {
	params := map[string]string{
		"order_id": fmt.Sprintf("%d", orderId),
	}
	response, err := helpers.SendRequest(params, "POST", os.Getenv("FASTEX_DOMAIN_URL")+"/api/v1/order/cancel", "fastex")
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	}
	return response
}

func CreateOrder(symbol string, amount float64, price float64, orderType string, tradeType string) int {
	params := map[string]string{
		"symbol":     symbol,
		"amount":     fmt.Sprintf("%f", amount),
		"price":      fmt.Sprintf("%f", price),
		"type":       orderType,
		"type_trade": tradeType,
	}
	response, err := helpers.SendRequest(params, "POST", os.Getenv("FASTEX_DOMAIN_URL")+"/api/v1/order/new", "fastex")
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return 0
	}

	var result struct {
		Response struct {
			Entity struct {
				OrderID int `json:"order_id"`
			} `json:"entity"`
		} `json:"response"`
	}

	if err := json.Unmarshal([]byte(response), &result); err != nil {
		fmt.Printf("Create Order Error Request: %s\n", params)
		fmt.Printf("Create Order Error Response: %s\n", response)
		fmt.Println("Ошибка парсинга Create Order JSON:", err)
		return 0
	}

	return result.Response.Entity.OrderID
}

func Trade(symbol config.Symbol, price float64) {
	volume := helpers.RoundToDecimals((1+rand.Float64()*10)*symbol.Volume, symbol.VolumePrecision)
	var firstOrderId = CreateOrder(symbol.FastexName, volume, price, "buy", "limit")
	var secondOrderId = CreateOrder(symbol.FastexName, volume, price, "sell", "limit")
	CloseOrder(firstOrderId)
	CloseOrder(secondOrderId)
}

func FillOrderbook(symbol config.Symbol, lastPriceRounded float64) {
	var ordersTypeBuyCount = 0
	var ordersTypeSellCount = 0
	for _, order := range OrderList[symbol.FastexName] {
		/*
			if order.Type == "buy" {
				fmt.Printf("BUY %f >= %f || %f < %f\r\n", order.Price, lastPriceRounded, order.Price, lastPriceRounded*0.92)
			} else {
				fmt.Printf("SELL %f <= %f || %f > %f\r\n", order.Price, lastPriceRounded, order.Price, lastPriceRounded*1.08)
			}
		*/

		if order.Type == "buy" && (order.Price >= lastPriceRounded || order.Price < lastPriceRounded*0.92) {
			CloseOrder(order.ID)
			continue
		} else if order.Type == "sell" && (order.Price <= lastPriceRounded || order.Price > lastPriceRounded*1.08) {
			CloseOrder(order.ID)
			continue
		}

		if order.Type == "buy" {
			ordersTypeBuyCount++
		} else if order.Type == "sell" {
			ordersTypeSellCount++
		}
	}

	var volume = 0.0
	for ordersTypeBuyCount < 5 {
		volume = helpers.RoundToDecimals((1+rand.Float64()*10)*symbol.Volume, symbol.VolumePrecision)
		var randomizedPrice = helpers.RoundToDecimals(lastPriceRounded*(1-symbol.Step*float64(ordersTypeBuyCount+1)), symbol.Precision)
		CreateOrder(symbol.FastexName, volume, randomizedPrice, "buy", "limit")
		ordersTypeBuyCount++
	}

	for ordersTypeSellCount < 5 {
		volume = helpers.RoundToDecimals((1+rand.Float64()*10)*symbol.Volume, symbol.VolumePrecision)
		var randomizedPrice = helpers.RoundToDecimals(lastPriceRounded*(1+symbol.Step*float64(ordersTypeSellCount+1)), symbol.Precision)
		CreateOrder(symbol.FastexName, volume, randomizedPrice, "sell", "limit")
		ordersTypeSellCount++
	}
}
