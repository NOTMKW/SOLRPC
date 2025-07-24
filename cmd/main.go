package main

import (
	"log"
	"time"

	"github.com/NOTMKW/RPC/internal/config"
	"github.com/NOTMKW/RPC/internal/service"
)

func main() {
	cfg := config.LoadConfig()
	priceChannel := make(chan []int)

	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		apikey := cfg.GetApiKey("API_KEY")
		for range ticker.C {
			if price := service.FetchSolanaPrice(apikey); price != nil {
				priceChannel <- []int{int(price.Usd), int(price.Eur), int(price.Btc)}
			} else {
				log.Println("Failed to fetch Solana price")
			}
		}
	}()

	for price := range priceChannel {
		if len(price) >= 3 {
			log.Printf("Current Solana Price: USD: %d, EUR: %d, BTC: %d",
				price[0], price[1], price[2])
		} else {
			log.Println("Received incomplete price data")
		}
	}
}
