package ndaxapi

import (
	"encoding/json"
	"log"
	"net/url"
	"os"

	"github.com/bAngerman/auto-investor/pkg/discord"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

// Payload used in request object, it contains the data of the request
type Payload map[string]interface{}

// Request is the body of a generic request.
type Request struct {
	Type     int    `json:"m"`
	Sequence int    `json:"i"`
	Function string `json:"n"`
	Payload  string `json:"o"`
}

// GenericResponse of a request.
type GenericResponse struct {
	Type     int    `json:"m"`
	Sequence int    `json:"i"`
	Function string `json:"n"`
	Response string `json:"o"`
}

type pingResponse map[string]interface{}

var (
	instruments    map[string]int
	sequence       int
	wsHost, wsPath string
)

func init() {
	godotenv.Load()

	wsHost = os.Getenv("NDAXIO_WS_HOST")
	wsPath = os.Getenv("NDAXIO_WS_PATH")

	sequence = 2
}

// Connect to the NDAX websocket connection
func Connect() (*websocket.Conn, error) {
	u := url.URL{
		Scheme: "wss",
		Host:   wsHost,
		Path:   wsPath,
	}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)

	if err != nil {
		log.Fatal("Error connecting to ws:", err)
	}

	return conn, err
}

// SendRequest executes a write to the websocket connection
// using a provided json encoded string
func SendRequest(req []byte, conn *websocket.Conn) []byte {
	err := conn.WriteMessage(websocket.TextMessage, []byte(req))

	if err != nil {
		log.Panic("write:", err)
	}

	_, message, err := conn.ReadMessage()

	if err != nil {
		log.Panic("Failed to get message:", err)
	}

	return message
}

// EncodeRequest JSON encodes the request, and payload params.
func EncodeRequest(r Request, p Payload) []byte {
	// Encode payload.
	var pJSON []byte
	pJSON, err := json.Marshal(p)

	// Add payload to request.
	r.Payload = string(pJSON)

	// Encode request.
	var rJSON []byte
	rJSON, err = json.Marshal(r)
	if err != nil {
		log.Panic("json:", err)
	}

	return rJSON
}

// DecodeResponse unmarshals the response from NDAX api.
func DecodeResponse(res []byte) GenericResponse {

	var data GenericResponse
	err := json.Unmarshal(res, &data)

	if err != nil {
		discord.Message("Error decoding response: " + err.Error())
		log.Panic("Error decoding response:", err)
	}

	return data
}

// Send sends a request to the NDAX api.
func Send(r Request, p Payload, conn *websocket.Conn) GenericResponse {
	rJSON := EncodeRequest(r, p)
	res := SendRequest(rJSON, conn)
	data := DecodeResponse(res)

	return data
}

// GetSequence get's the next number in sequence for requests.
func GetSequence() int {
	return (sequence + 2)
}

// Ping sends a ping to the websocket connection to keep it alive.
func Ping(conn *websocket.Conn) {

	p := Payload{}
	r := Request{
		Type:     0,
		Sequence: GetSequence(),
		Function: "Ping",
	}

	data := Send(r, p, conn)

	var pingData pingResponse
	err := json.Unmarshal([]byte(data.Response), &pingData)

	if err != nil || pingData["errormsg"] != nil {
		discord.Message("Bad ping: " + pingData["errormsg"].(string))
		log.Panic("Error in ping:", err)
	}
}
