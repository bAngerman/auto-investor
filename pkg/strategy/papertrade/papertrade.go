package papertrade

import (
	"log"

	"github.com/bAngerman/auto-investor/pkg/ndaxapi/account"
	"github.com/bAngerman/auto-investor/pkg/ndaxapi/order"
	"github.com/bAngerman/auto-investor/pkg/ndaxapi/ticker"
	"github.com/bAngerman/auto-investor/pkg/strategy/papertrade/paperholdings"

	"github.com/gorilla/websocket"
)

var (
	c *websocket.Conn
)

// Start new paper trade strategy
func Start(conn *websocket.Conn, tickerChan chan ticker.StructuredData) {
	// Set to package var.
	c = conn

	log.Println("Using papertrade strategy")
}

func run() {
	// log.Println("Running tasks.")

	// Get current account holdings
	// accounts, orders := evaluateHoldings()

	// Evaluate whether trades should execute.
	// shouldTradesExecute(accounts, orders)

	// // Freshen accounts / orders
	// accounts, orders = evaluateHoldings()

	// // Offer new trade statuses based on account / orders
	// submitNewTrades(accounts, orders)
}

func evaluateHoldings() (account.Accounts, order.Orders) {

	// Get fake values for paper trading
	accounts := paperholdings.GetAccountPositions()
	orders := paperholdings.GetOpenOrders()

	return accounts, orders
}

// func shouldTradesExecute(account.Accounts, order.Orders) bool {

// }

// func submitNewTrades(account.Accounts, order.Orders) {

// }
