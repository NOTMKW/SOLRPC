package main

import (
	"log"
	"os"
	"time"

	"github.com/NOTMKW/RPC/internal/service"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}
	priceChannel := make(chan []int)

	apikey := os.Getenv("d3871320f46376d61cede7af78d3d37b653641b09314b51423ab0c4b8de4bf41")
	if apikey == "" {
		log.Fatal("API key not found in environment variables")

		go func() {
			ticker := time.NewTicker(30 * time.Second)
			defer ticker.Stop()

			for {
				price := service.FetchSolanaPrice(apikey)
				priceChannel <- price
				<-ticker.C
			}
		}()
		for price := range priceChannel {
			log.Printf("Current Solana Prices: USD: %d, EUR: %d, BTC: %d", price[0], price[1], price[2])
		}
	}
}
