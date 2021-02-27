package main

import (
	"log"
	"os"
	"strings"

	"github.com/bAngerman/auto-investor/pkg/ndaxapi"
	"github.com/bAngerman/auto-investor/pkg/ndaxapi/auth"
	"github.com/bAngerman/auto-investor/pkg/ndaxapi/ticker"
	"github.com/bAngerman/auto-investor/pkg/strategy/livetrade"
	"github.com/bAngerman/auto-investor/pkg/strategy/papertrade"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Panic("Error parsing .env file:", err)
	}
}

// 1. Connect, and Authenticate with NDAX.
// 2. Subscribe to Tickers utilizing a channel.
// 3. Within strategy, complete actions based on ticker events.
func main() {
	log.Println("Starting auto trader.")

	conn, err := ndaxapi.Connect()
	defer conn.Close()

	// Create a channel for the subscribed ticker information to go through
	tickerChan := make(chan ticker.StructuredData, 4)

	// go func() {
	// 	for {
	// 		_, res, err := conn.ReadMessage()
	// 		if err == nil {
	// 			data := ndaxapi.DecodeResponse(res)

	// 			var tickerData ticker.TickerData
	// 			err := json.Unmarshal([]byte(data.Response), &tickerData)

	// 			if err != nil {
	// 				log.Panic("Error decoding ticker response:", err)
	// 			}

	// 			// log.Println(tickerData)

	// 			tickerChan <- tickerData
	// 		} else {
	// 			log.Panic("Error in ticker response:", err)
	// 		}
	// 	}
	// }()

	if err != nil {
		log.Panic("Error connecting to NDAX:", err)
	}

	err = auth.Authenticate(conn)
	if err != nil {
		log.Panic("Error authenticating with NDAX:", err)
	}
	log.Println("Connected and authenticated to NDAX.")

	// Get symbols to watch
	tickers := ticker.GetTickerSymbols()

	log.Println("Subscribing to ticker(s):", strings.Join(tickers, ", "))

	go ticker.SubscribeToTickers(conn, tickerChan, tickers)

	log.Println("Ticker Channel:", tickerChan)

	if os.Getenv("ENV") == "live" {
		// prompt := promptui.Prompt{
		// 	Label:     "env states to use live API and real money, are you sure? (y/n)",
		// 	IsConfirm: true,
		// }

		// result, err := prompt.Run()

		// if err != nil {
		// 	fmt.Printf("Prompt failed %v\n", err)
		// 	return
		// }

		// if result == "y" {
		livetrade.Start(conn, tickerChan)
		// }/
	} else {
		papertrade.Start(conn, tickerChan)
	}
}
