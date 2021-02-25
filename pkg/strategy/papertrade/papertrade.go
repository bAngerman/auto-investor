package papertrade

import (
	"log"

	"github.com/bAngerman/auto-investor/pkg/ndaxapi/account"
	"github.com/bAngerman/auto-investor/pkg/taskrunner"
	"github.com/boltdb/bolt"
	"github.com/gorilla/websocket"
)

var (
	c *websocket.Conn
)

func init() {
	// Initialize the database, if it does not have a db, it will create one.
	_, err := bolt.Open("paper", 0600, nil)
	if err != nil {
		log.Panic("Error creating paper db:", err)
	}
}

// Start new paper trade strategy
func Start(conn *websocket.Conn) {
	// Set to package var.
	c = conn

	log.Println("Using papertrade strategy")
	taskrunner.Initloop(conn, run)

}

func run() {
	log.Println("Running tasks.")

	checkAccount()
	// Check account status
	// Evaluate current trade statuses
	// Offer new trade statuses
	// If any events to report, discord DM them
}

func checkAccount() {
	account.GetAccountPosition(c)
}
