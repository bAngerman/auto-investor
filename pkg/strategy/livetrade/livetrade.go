package livetrade

import (
	"log"

	"github.com/gorilla/websocket"
)

// Start new live trade strategy
func Start(conn *websocket.Conn) {
	log.Println("Using livetrade strategy")
	// ndaxapi.Initloop(loop)
}

func loop() {
	log.Println("Running loop")
}
