package paperholdings

import (
	"log"

	"github.com/asdine/storm/q"
	"github.com/asdine/storm/v3"

	"github.com/bAngerman/auto-investor/pkg/ndaxapi/account"
	"github.com/bAngerman/auto-investor/pkg/ndaxapi/order"
	"github.com/bAngerman/auto-investor/pkg/ndaxapi/ticker"
)

var dbName string = "paper.db" // name of the bolt db

func init() {
	db := openDB()
	defer db.Close()
}

func openDB() *storm.DB {
	db, err := storm.Open(dbName)
	if err != nil {
		log.Panic("Error creating DB:", err)
	}
	return db
}

func setAccountPosition(accs account.Accounts, db *storm.DB) error {
	for _, account := range accs {
		err := db.Save(&account)
		if err != nil {
			return err
		}
	}
	return nil
}

func createAccountPositions(db *storm.DB) account.Accounts {
	accs := account.Accounts{
		{
			OMSId:         1,
			ID:            1,
			ProductSymbol: "CAD",
			Amount:        100.00000,
		},
		{
			OMSId:         1,
			ID:            2,
			ProductSymbol: "DOGE",
			Amount:        0,
		},
	}

	err := setAccountPosition(accs, db)

	if err != nil {
		log.Panic("Error setting initial account positions:", err)
	}

	return accs
}

// GetAccountPositions spoofs some accounts, from a stored key/value db.
func GetAccountPositions() account.Accounts {
	db := openDB()
	defer db.Close()

	tickers := ticker.GetTickerSymbols()

	var accs account.Accounts
	err := db.Select(q.In("ProductSymbol", tickers)).Find(&accs)

	if err != nil {
		log.Println("Created account positions as they were not present.")
		accs = createAccountPositions(db)
	}

	return accs
}

// GetOpenOrders spoofs some orders, from a stored key / value db.
func GetOpenOrders() order.Orders {
	db := openDB()
	defer db.Close()

	var orders order.Orders
	err := db.All(&orders)

	if err != nil {
		log.Panic("Error getting orders.", err)
	}

	return orders
}
