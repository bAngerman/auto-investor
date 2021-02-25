package order

import (
	"encoding/json"
	"log"

	"github.com/bAngerman/auto-investor/pkg/ndaxapi"
	"github.com/bAngerman/auto-investor/pkg/ndaxapi/account"
	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/websocket"
)

// Order is the order data from the API
type Order struct {
	Side                 string  `json:"Side"`                 // "Sell",
	OrderID              int64   `json:"OrderId"`              // 14307301634,
	Price                int64   `json:"Price"`                // 1,
	Quantity             float64 `json:"Quantity"`             // 100,
	DisplayQuantity      float64 `json:"DisplayQuantity"`      // 100,
	Instrument           int64   `json:"Instrument"`           // 77,
	Account              int64   `json:"Account"`              // 32087,
	AccountName          string  `json:"AccountName"`          // "brendan.angerman@gmail.com",
	OrderType            string  `json:"OrderType"`            // "Limit",
	ClientOrderID        int64   `json:"ClientOrderId"`        // 0,
	OrderState           string  `json:"OrderState"`           // "Working",
	ReceiveTime          int64   `json:"ReceiveTime"`          // 1614226199689,
	ReceiveTimeTicks     int64   `json:"ReceiveTimeTicks"`     // 637498229996886294,
	LastUpdatedTime      int64   `json:"LastUpdatedTime"`      // 1614226199688,
	LastUpdatedTimeTicks int64   `json:"LastUpdatedTimeTicks"` // 637498229996883117,
	OrigQuantity         int64   `json:"OrigQuantity"`         // 100,
	QuantityExecuted     int64   `json:"QuantityExecuted"`     // 0,
	GrossValueExecuted   int64   `json:"GrossValueExecuted"`   // 0,
	ExecutableValue      int64   `json:"ExecutableValue"`      // 0,
	AvgPrice             int64   `json:"AvgPrice"`             // 0,
	CounterPartyID       int64   `json:"CounterPartyId"`       // 0,
	ChangeReason         string  `json:"ChangeReason"`         // "NewInputAccepted",
	OrigOrderID          int64   `json:"OrigOrderId"`          // 14307301634,
	OrigClOrdID          int64   `json:"OrigClOrdId"`          // 0,
	EnteredBy            int64   `json:"EnteredBy"`            // 31943,
	UserName             string  `json:"UserName"`             // "brendan.angerman@gmail.com",
	IsQuote              bool    `json:"IsQuote"`              // false,
	InsideAsk            float64 `json:"InsideAsk"`            // 0.06996993,
	InsideAskSize        int64   `json:"InsideAskSize"`        // 19000,
	InsideBid            float64 `json:"InsideBid"`            // 0.068308,
	InsideBidSize        float64 `json:"InsideBidSize"`        // 13333,
	LastTradePrice       float64 `json:"LastTradePrice"`       // 0.06994568,
	RejectReason         string  `json:"RejectReason"`         // "",
	IsLockedIn           bool    `json:"IsLockedIn"`           // false,
	CancelReason         string  `json:"CancelReason"`         // "",
	OrderFlag            string  `json:"OrderFlag"`            // "AddedToBook",
	UseMargin            bool    `json:"UseMargin"`            // false,
	StopPrice            int64   `json:"StopPrice"`            // 0,
	PegPriceType         string  `json:"PegPriceType"`         // "Last",
	PegOffset            int64   `json:"PegOffset"`            // 0,
	PegLimitOffset       int64   `json:"PegLimitOffset"`       // 0,
	OMSId                int64   `json:"OMSId"`                // 1
}

// Orders is an array of Order
type Orders []Order

// GetOpenOrders gets the open orders for the user.
func GetOpenOrders(conn *websocket.Conn) Orders {
	accID := account.GetUserAccountID(conn)

	p := ndaxapi.Payload{
		"AccountId": accID,
		"OMSId":     1,
	}

	r := ndaxapi.Request{
		Type:     0,
		Sequence: ndaxapi.GetSequence(),
		Function: "GetOpenOrders",
	}

	data := ndaxapi.Send(r, p, conn)

	var orders Orders
	err := json.Unmarshal([]byte(data.Response), &orders)

	if err != nil {
		log.Panic(err)
	}

	spew.Dump(orders)

	return orders
}
