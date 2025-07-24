package Price

type PriceResponse struct {
	USD float64 `json:"USD"`
	EUR float64 `json:"EUR"`
	BTC float64 `json:"BTC"`
}

type SolanaPrice struct {
	Usd float64
	Eur float64
	Btc float64
}
