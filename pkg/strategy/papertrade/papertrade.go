package papertrade

import (
	"log"

	"github.com/bAngerman/auto-investor/pkg/ndaxapi"
	account "github.com/bAngerman/auto-investor/pkg/ndaxapi/account"
)

// Start new paper trade strategy
func Start() {
	log.Println("Using papertrade strategy")
	ndaxapi.Initloop(run)
}

func run() {
	checkAccount()
	// Check account status
	// Evaluate current trade statuses
	// Offer new trade statuses
	// If any events to report, discord DM them
}

func checkAccount() {
	acc := account.GetAccountPosition()
}
