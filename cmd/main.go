package main

import (
	"log"
	"time"

	"github.com/NOTMKW/RPC/internal/config"
	Price "github.com/NOTMKW/RPC/internal/price"
	printer "github.com/NOTMKW/RPC/internal/printer"
	"github.com/NOTMKW/RPC/internal/service"
)

func main() {
	cfg := config.LoadConfig()
	pricechanstruct := Price.SolanaPrice{}
	pricechan := make(chan *Price.SolanaPrice, 10)
	printer.NewPriceService(&pricechanstruct)
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		apikey := cfg.GetApiKey("apikey")

		if price, err := service.FetchSolanaPrice(apikey); err == nil {
			log.Println("sending price to channel:")
			pricechan <- price
		} else {
			log.Println("Failed to fetch Solana price:", err)
		}
		for range ticker.C {
			if price, err := service.FetchSolanaPrice(apikey); err == nil {
				log.Println("sending price to channel:")
				pricechan <- price
			} else {
				log.Println("Failed to fetch Solana price:", err)
			}
		}
		log.Println("Application started successfully")
	}()
}
