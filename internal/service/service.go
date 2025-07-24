package service

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	Price "github.com/NOTMKW/RPC/internal/price"
)

func FetchSolanaPrice(apikey string) (*Price.SolanaPrice, error) {
	var body []byte
	for i := 0; i < 10; i++ {
		resp, err := http.Get("https://min-api.cryptocompare.com/data/price?fsym=SOL&tsyms=USD,EUR,BTC&api_key=" + apikey)
		if err != nil {
			log.Println("Error fetching Solana price:", err)
			continue
		}
		defer resp.Body.Close()
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error reading response body:", err)
			return nil, err
		}
	}

	var priceData Price.PriceResponse
	if err := json.Unmarshal(body, &priceData); err != nil {
		log.Println("Error parsing JSON:", err)
		return nil, err
	}
	return &Price.SolanaPrice{
		Usd: priceData.USD,
		Eur: priceData.EUR,
		Btc: priceData.BTC,
	}, nil
}
