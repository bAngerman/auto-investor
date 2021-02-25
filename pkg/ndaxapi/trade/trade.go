package trade

import (
	"encoding/json"
	"log"

	"github.com/bAngerman/auto-investor/pkg/ndaxapi"
	"github.com/bAngerman/auto-investor/pkg/ndaxapi/account"
	"github.com/gorilla/websocket"
)

// Trade is a NDAX trade structure
type Trade struct {
	OmsID                    int64   `json:"OMSId"`                    // 1,
	ExecutionID              int64   `json:"ExecutionId"`              // 675528,
	TradeID                  int64   `json:"TradeId"`                  // 429092,
	OrderID                  int64   `json:"OrderId"`                  // 14297295238,
	AccountID                int64   `json:"AccountId"`                // 32087,
	AccountName              string  `json:"AccountName"`              // "brendan.angerman@gmail.com",
	SubAccountID             int64   `json:"SubAccountId"`             // 0,
	ClientOrderID            int64   `json:"ClientOrderId"`            // 0,
	InstrumentID             int64   `json:"InstrumentId"`             // 77,
	Side                     string  `json:"Side"`                     // "Buy",
	OrderType                string  `json:"OrderType"`                // "Limit",
	Quantity                 float64 `json:"Quantity"`                 // 2232.0000000000000000000000000,
	RemainingQuantity        float64 `json:"RemainingQuantity"`        // 0.0000000000000000000000000000,
	Price                    float64 `json:"Price"`                    // 0.0850000000000000000000000000,
	Value                    float64 `json:"Value"`                    // 189.72000000000000000000000000,
	CounterParty             string  `json:"CounterParty"`             // "37647",
	OrderTradeRevision       int64   `json:"OrderTradeRevision"`       // 1,
	Direction                string  `json:"Direction"`                // "NoChange",
	IsBlockTrade             bool    `json:"IsBlockTrade"`             // false,
	Fee                      float64 `json:"Fee"`                      // 4.4640000000000000000000000000,
	FeeProductID             int64   `json:"FeeProductId"`             // 10,
	OrderOriginator          int64   `json:"OrderOriginator"`          // 31943,
	UserName                 string  `json:"UserName"`                 // "brendan.angerman@gmail.com",
	TradeTimeMS              int64   `json:"TradeTimeMS"`              // 1613246642358,
	MakerTaker               string  `json:"MakerTaker"`               // "Maker",
	AdapterTradeID           int64   `json:"AdapterTradeId"`           // 0,
	InsideBid                float64 `json:"InsideBid"`                // 0.0851000000000000000000000000,
	InsideBidSize            float64 `json:"InsideBidSize"`            // 5000.0000000000000000000000000,
	InsideAsk                float64 `json:"InsideAsk"`                // 0.0866990000000000000000000000,
	InsideAskSize            float64 `json:"InsideAskSize"`            // 3034.0000000000000000000000000,
	IsQuote                  bool    `json:"IsQuote"`                  // false,
	CounterPartyClientUserID int64   `json:"CounterPartyClientUserId"` // 1,
	NotionalProductID        int64   `json:"NotionalProductId"`        // 6,
	NotionalRate             float64 `json:"NotionalRate"`             // 0.7874945859747214000000000000,
	NotionalValue            float64 `json:"NotionalValue"`            // 149.40347285112414400800000000,
	NotionalHoldAmount       int64   `json:"NotionalHoldAmount"`       // 0,
	TradeTime                int64   `json:"TradeTime"`                // 637488434423575198
}

// Trades is a collection of Trade
type Trades []Trade

// GetAccountTrades gets recent trades.
func GetAccountTrades(conn *websocket.Conn) Trades {
	accID := account.GetUserAccountID(conn)

	p := ndaxapi.Payload{
		"AccountId": accID,
		"OMSId":     1,
	}

	r := ndaxapi.Request{
		Type:     0,
		Sequence: ndaxapi.GetSequence(),
		Function: "GetAccountTrades",
	}

	data := ndaxapi.Send(r, p, conn)

	var trades Trades
	err := json.Unmarshal([]byte(data.Response), &trades)

	if err != nil {
		log.Panic(err)
	}

	return trades
}

// GetOpenTradePositions gets the current trade positions for the user.
func GetOpenTradePositions(conn *websocket.Conn) Trades {
	accID := account.GetUserAccountID(conn)

	p := ndaxapi.Payload{
		"AccountId": accID,
		"OMSId":     1,
	}

	r := ndaxapi.Request{
		Type:     0,
		Sequence: ndaxapi.GetSequence(),
		Function: "GetOpenTradeReports",
	}

	data := ndaxapi.Send(r, p, conn)

	var trades Trades
	err := json.Unmarshal([]byte(data.Response), &trades)

	if err != nil {
		log.Panic(err)
	}

	return trades
}
