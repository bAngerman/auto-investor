package ticker

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/bAngerman/auto-investor/pkg/ndaxapi"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

// Data is the response from the
// websocket subscribe for a ticker
type Data [][]float64

// StructuredData is a structured object representing
// data sent in the ticker subscription.
type StructuredData struct {
	DateTime       float64 // 1501603632000, \\DateTime - UTC - Milliseconds since 1/1/1970
	High           float64 // 2700.33, \\High
	Low            float64 // 2701.2687.01, \\Low
	Open           float64 // 2687.01, \\Open
	Close          float64 // 2687.01, \\Close
	Volume         float64 // 24.86100992, \\Volume
	InsideBidPrice float64 // 0, \\Inside Bid Price
	InsideAskPrice float64 // 2870.95, \\Inside Ask Price
	InstrumentID   float64 // 1 \\InstrumentId
}

var (
	instrumentIDs    map[string]int
	tickerSymbolsArr []string
)

func init() {
	godotenv.Load()

	instrumentIDs = map[string]int{
		"BTCCAD":  1,
		"BCHCAD":  2,
		"ETHCAD":  3,
		"XRPCAD":  4,
		"LTCCAD":  5,
		"EOSCAD":  75,
		"XLMCAD":  76,
		"DOGECAD": 77,
		"ADACAD":  78,
	}

	tickerSymbolsArr = strings.Split(os.Getenv("NDAX_TICKERS"), ",")
}

func getInstrumentID(instrumentString string) int {
	if val, ok := instrumentIDs[instrumentString]; ok {
		return val
		//do something here
	}

	log.Fatal("Ticker does not exist:", instrumentString)
	return 0
}

// subscribeTicker listens to ticker events.
func subscribeTicker(conn *websocket.Conn, instrumentString string) {
	p := ndaxapi.Payload{
		"OMSId":            1,
		"InstrumentId":     getInstrumentID(instrumentString),
		"Interval":         60,
		"IncludeLastCount": 100,
	}
	r := ndaxapi.Request{
		Type:     0,
		Sequence: ndaxapi.GetSequence(),
		Function: "SubscribeTicker",
	}

	// Send subscribe request
	_ = ndaxapi.Send(r, p, conn)
}

// GetTickerSymbols reads from the .env file to
// get ticker that user cares about.
func GetTickerSymbols() []string {
	return tickerSymbolsArr
}

// SubscribeToTickers listens to ticker events, and sends them to the channel
func SubscribeToTickers(conn *websocket.Conn, c chan<- StructuredData, tickers []string) {
	for _, t := range tickers {
		subscribeTicker(conn, t)
	}

	for {
		_, res, err := conn.ReadMessage()
		if err == nil {
			data := ndaxapi.DecodeResponse(res)

			var tickerData Data
			err := json.Unmarshal([]byte(data.Response), &tickerData)

			if err != nil {
				log.Panic("Error in ticker response:", err)
			}

			if len(tickerData) > 0 {
				fData := tickerData[0]

				sData := StructuredData{
					DateTime:       fData[0],
					High:           fData[1],
					Low:            fData[2],
					Open:           fData[3],
					Close:          fData[4],
					Volume:         fData[5],
					InsideBidPrice: fData[6],
					InsideAskPrice: fData[7],
					InstrumentID:   fData[8],
				}

				c <- sData
			}

		}
	}
}
