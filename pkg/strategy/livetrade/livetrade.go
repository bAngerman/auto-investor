package livetrade

import (
	"log"

	"github.com/bAngerman/auto-investor/pkg/ndaxapi/ticker"
	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/websocket"
)

// Start new live trade strategy
func Start(conn *websocket.Conn, tickerChan <-chan ticker.StructuredData) {
	log.Println("Using livetrade strategy")

	for {
		data, err := <-tickerChan

		if err == false {
			log.Panic("Error receiving ticker payload:", err)
		}

		spew.Dump(data)
	}
}
