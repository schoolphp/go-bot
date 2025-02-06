package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"project/config"
	"project/helpers"
	"project/services/fastex"
	"project/services/pricer"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	if os.Getenv("IS_BOT_ENABLED") == "0" {
		log.Fatal("Бот отключен")
	}

	go runLiquidityBot()
	select {}
}

func runLiquidityBot() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		fastex.GetAllOrders()

		for _, symbol := range config.Symbols {
			go fullFitOrderBookBySymbol(symbol)
		}
	}
}

func fullFitOrderBookBySymbol(symbol config.Symbol) {
	lastPrice := pricer.GetLastPriceBySymbol(symbol.HostSymbolName, symbol.HostName)
	fmt.Printf("Symbol %s Host %s Price: %f\r\n", symbol.FastexName, symbol.HostName, lastPrice)

	if lastPrice == 0 {
		return
	}

	lastPriceRounded := helpers.RoundToDecimals(lastPrice, symbol.Precision)

	fastex.Trade(symbol, lastPriceRounded)
	fastex.FillOrderbook(symbol, lastPriceRounded)
}
