package account

import (
	"encoding/json"
	"log"
	"os"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"

	"github.com/bAngerman/auto-investor/pkg/ndaxapi"
	"github.com/bAngerman/auto-investor/pkg/ndaxapi/auth"
)

// Account is the structure of an account in NDAX
type Account struct {
	OMSId                      int64   `json:"OMSId"`                      // 1,
	AccountID                  int64   `json:"AccountId"`                  // 32087,
	ProductSymbol              string  `json:"ProductSymbol"`              // "BTC",
	ProductID                  int64   `json:"ProductId"`                  // 1,
	Amount                     float64 `json:"Amount"`                     // 0,
	Hold                       float64 `json:"Hold"`                       // 0,
	PendingDeposits            float64 `json:"PendingDeposits"`            // 0,
	PendingWithdraws           float64 `json:"PendingWithdraws"`           // 0,
	TotalDayDeposits           float64 `json:"TotalDayDeposits"`           // 0,
	TotalMonthDeposits         float64 `json:"TotalMonthDeposits"`         // 0,
	TotalYearDeposits          float64 `json:"TotalYearDeposits"`          // 0,
	TotalDayDepositNotional    float64 `json:"TotalDayDepositNotional"`    // 0,
	TotalMonthDepositNotional  float64 `json:"TotalMonthDepositNotional"`  // 0,
	TotalYearDepositNotional   float64 `json:"TotalYearDepositNotional"`   // 0,
	TotalDayWithdraws          int64   `json:"TotalDayWithdraws"`          // 0,
	TotalMonthWithdraws        int64   `json:"TotalMonthWithdraws"`        // 0,
	TotalYearWithdraws         int64   `json:"TotalYearWithdraws"`         // 0,
	TotalDayWithdrawNotional   int64   `json:"TotalDayWithdrawNotional"`   // 0,
	TotalMonthWithdrawNotional int64   `json:"TotalMonthWithdrawNotional"` // 0,
	TotalYearWithdrawNotional  int64   `json:"TotalYearWithdrawNotional"`  // 0,
	NotionalProductID          int64   `json:"NotionalProductId"`          // 6,
	NotionalProductSymbol      string  `json:"NotionalProductSymbol"`      // "USD",
	NotionalValue              float64 `json:"NotionalValue"`              // 267.86320000000000000000000000,
	NotionalHoldAmount         float64 `json:"NotionalHoldAmount"`         // 0.00,
	NotionalRate               float64 `json:"NotionalRate"`               // 50329.00
}

// Array of account structs
type Accounts []Account

type getUserAccountsResponse []int

func init() {
	godotenv.Load()
}

// GetUserAccountID will call the API, and get user accounts
func GetUserAccountID(conn *websocket.Conn) int {
	omsID := 1
	username := os.Getenv("NDAXIO_USER")
	userID := auth.UserID

	p := ndaxapi.Payload{
		"OMSId":    omsID,
		"UserId":   userID,
		"UserName": username,
	}

	r := ndaxapi.Request{
		Type:     0,
		Sequence: ndaxapi.GetSequence(),
		Function: "GetUserAccounts",
	}

	data := ndaxapi.Send(r, p, conn)

	var accounts getUserAccountsResponse
	err := json.Unmarshal([]byte(data.Response), &accounts)

	if err != nil {
		log.Panic("Error getting user account ID:", err)
	}

	// Return first in arr.
	if len(accounts) > 0 {
		return accounts[0]
	}

	return 0
}

// GetAccountPosition will call the API, and return the account status
func getAccountPosition(conn *websocket.Conn) Accounts {

	accID := GetUserAccountID(conn)

	p := ndaxapi.Payload{
		"AccountId": accID,
		"OMSId":     1,
	}

	r := ndaxapi.Request{
		Type:     0,
		Sequence: ndaxapi.GetSequence(),
		Function: "GetAccountPositions",
	}

	data := ndaxapi.Send(r, p, conn)

	var accounts Accounts
	err := json.Unmarshal([]byte(data.Response), &accounts)

	if err != nil {
		log.Panic(err)
	}

	return accounts
}
