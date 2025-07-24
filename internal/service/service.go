package service

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	Price "github.com/NOTMKW/RPC/internal/price"
)

func FetchSolanaPrice(apikey string) []int {
	var price []int
	var body []byte
	for i := 0; i < 10; i++ {
		resp, err := http.Get("https://min-api.cryptocompare.com/data/price?fsym=SOL&tsyms=USD,EUR,BTC&api_key=" + apikey)
		if err != nil {
			log.Println("Error fetching Solana price:", err)
			continue
		}
		body, err = io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			log.Println("Error reading response body:", err)
			continue
		}
		time.Sleep(30 * time.Second)
	}

	var priceData Price.PriceResponse
	if err := json.Unmarshal(body, &priceData); err != nil {
		log.Fatal("Error parsing JSON:", err)
	}

	price = append(price, int(priceData.USD))
	price = append(price, int(priceData.EUR))
	price = append(price, int(priceData.BTC*100000))

	return price
}
