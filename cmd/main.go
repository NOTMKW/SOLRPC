package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/NOTMKW/RPC/internal/config"
	Price "github.com/NOTMKW/RPC/internal/price"
	printer "github.com/NOTMKW/RPC/internal/printer"
	"github.com/NOTMKW/RPC/internal/service"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var (
	clients    = make(map[*websocket.Conn]bool)
	clientsMux = sync.RWMutex{}
)

func setupRoutes() {
	http.HandleFunc("/", homepage)
	http.HandleFunc("/ws", wsEndpoint)
}

func homepage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the Solana Price WebSocket Server!\n")
}

func addClient(conn *websocket.Conn) {
	clientsMux.Lock()
	clients[conn] = true
	clientsMux.Unlock()
}

func removeClient(conn *websocket.Conn) {
	clientsMux.Lock()
	delete(clients, conn)
	clientsMux.Unlock()
	conn.Close()
}
func broadcastPrice(price Price.SolanaPrice) {
	clientsMux.RLock()
	defer clientsMux.RUnlock()

	priceJSON, err := json.Marshal(price)
	if err != nil {
		log.Println("Error Marshaling price:", err)
		return
	}

	for client := range clients {
		if err := client.WriteMessage(websocket.TextMessage, priceJSON); err != nil {
			log.Println("Error sending message to client:", err)
			go removeClient(client)
		}
	}
}
func reader(conn *websocket.Conn) {

	defer removeClient(conn)

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			return
		}

		log.Println("Received message:", string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}

	log.Println("Client connected to WebSocket")
	addClient(ws)
	reader(ws)
}

func main() {
	cfg := config.LoadConfig()
	pricechanstruct := Price.SolanaPrice{}
	printerService := printer.NewPriceService(&pricechanstruct)

	printerService.StartFetchingPrices()

	go func() {
		for price := range printerService.PriceChannel {
			log.Println("Received price from channel:", price)
			broadcastPrice(*price)
		}
	}()

	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		apikey := cfg.GetApiKey("apikey")

		if price, err := service.FetchSolanaPrice(apikey); err == nil {
			log.Println("sending price to channel:")
			printerService.PriceChannel <- price
		} else {
			log.Println("Failed to fetch Solana price:", err)
		}
		for range ticker.C {
			if price, err := service.FetchSolanaPrice(apikey); err == nil {
				log.Println("sending price to channel:")
				printerService.PriceChannel <- price
			} else {
				log.Println("Failed to fetch Solana price:", err)
			}
		}
		log.Println("Application started successfully")
	}()
	setupRoutes()

	log.Println("Starting server on :8080")

	go func() {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	select {}
}
