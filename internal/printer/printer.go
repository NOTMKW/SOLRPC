package printer

import (
	"log"

	Price "github.com/NOTMKW/RPC/internal/price"
)

type NewPriceServicestruct struct {
	priceChannel chan *Price.SolanaPrice
}

func NewPriceService(ch *Price.SolanaPrice) *NewPriceServicestruct {
	return &NewPriceServicestruct{priceChannel: make(chan *Price.SolanaPrice)}
}

func (ps *NewPriceServicestruct) Start() {
	ps.priceChannel = make(chan *Price.SolanaPrice)
}

func (ps *NewPriceServicestruct) GetPriceChannel() chan<- *Price.SolanaPrice {
	return ps.priceChannel
}

func (ps *NewPriceServicestruct) StartFetchingPrices() {
	go func() {
		log.Println("Starting to fetch Solana prices...")
		for price := range ps.priceChannel {
			ps.printPrice(price)
		}
	}()
}

func (ps *NewPriceServicestruct) printPrice(price *Price.SolanaPrice) {
	log.Printf("Current Solana Price: USD: %.2f, EUR: %.2f, BTC: %.8f",
		price.Usd, price.Eur, price.Btc)
}
