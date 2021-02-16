package ndaxapi

import (
	"encoding/json"
	"errors"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/pquerna/otp/totp"

	"github.com/bAngerman/auto-investor/pkg/discord"
)

// Payload used in request object
type Payload map[string]interface{}

// Data - generic data for requests or responses
type Data struct {
	Type     int    `json:"m"`
	Sequence int    `json:"i"`
	Function string `json:"n"`
	Payload  string `json:"o"`
}

var (
	conn               *websocket.Conn
	wsHost, wsPath     string
	username, password string

	twoFA  string
	sToken string

	instruments map[string]int
	sequence    int
)

func init() {
	godotenv.Load()

	twoFA = os.Getenv("NDAXIO_2FA")
	wsHost = os.Getenv("NDAXIO_WS_HOST")
	wsPath = os.Getenv("NDAXIO_WS_PATH")
	username = os.Getenv("NDAXIO_USER")
	password = os.Getenv("NDAXIO_PASS")

	instruments = map[string]int{
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

	sequence = 2
}

func sendRequest(c *websocket.Conn, r Data, p Payload) Payload {
	var pJSON []byte
	pJSON, err := json.Marshal(p)

	r.Payload = string(pJSON)

	var rJSON []byte
	rJSON, err = json.Marshal(r)

	if err != nil {
		log.Panic("json:", err)
	}
	err = c.WriteMessage(websocket.TextMessage, rJSON)
	if err != nil {
		log.Panic("write:", err)
	}

	_, message, err := c.ReadMessage()

	if err != nil {
		log.Panic("Failed to get message:", err)
	}

	var res Data

	err = json.Unmarshal(message, &res)
	if err != nil {
		log.Panic("Failed to decode json:", err)
	}

	var data Payload

	err = json.Unmarshal([]byte(res.Payload), &data)
	if err != nil {
		log.Panic("Failed to decode payload json:", err)
	}

	return data
}

func getSequence() int {
	return (sequence + 2)
}

func authenticate(c *websocket.Conn) error {
	err := login(c)
	err = authenticate2FA(c)

	return err
}

func login(c *websocket.Conn) error {
	p := Payload{
		"UserName": username,
		"Password": password,
	}

	r := Data{
		Type:     0,
		Sequence: getSequence(),
		Function: "AuthenticateUser",
	}

	res := sendRequest(c, r, p)

	if res["errormsg"] != nil {
		discord.Message("Bad auth request: " + res["errormsg"].(string))
		return errors.New(res["errormsg"].(string))
	}

	return nil
}

func authenticate2FA(c *websocket.Conn) error {
	time := time.Now()
	key, err := totp.GenerateCode(twoFA, time)

	if err != nil {
		discord.Message("Bad 2FA generation: " + err.Error())
		return err
	}

	p := Payload{
		"Code": key,
	}

	r := Data{
		Type:     0,
		Sequence: getSequence(),
		Function: "Authenticate2FA",
	}

	res := sendRequest(c, r, p)

	if res["errormsg"] != nil {
		discord.Message("Bad 2FA request: " + res["errormsg"].(string))
		return errors.New(res["errormsg"].(string))
	}

	return nil
}

func ping() {

	p := Payload{}
	r := Data{
		Type:     0,
		Sequence: getSequence(),
		Function: "Ping",
	}

	res := sendRequest(conn, r, p)

	log.Printf("PING, %s", res["msg"])

	if res["errormsg"] != nil {
		discord.Message("Bad ping: " + res["errormsg"].(string))
	}
}

// Start the ndaxio websocket process
func Start() (*websocket.Conn, error) {
	u := url.URL{
		Scheme: "wss",
		Host:   wsHost,
		Path:   wsPath,
	}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)

	if err != nil {
		log.Fatal("dial:", err)
	}

	// Set to var to close later.
	conn = c

	err = authenticate(c)

	return c, err
}

// Initloop the loop used for trading
func Initloop(callback func()) {

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	done := make(chan struct{})

	pause := time.Second * 10

	ticker := time.NewTicker(pause)
	defer ticker.Stop()

	for {
		select {

		case <-done:
			return

		case <-ticker.C:

			// Always ping to keep conn alive.
			ping()

			callback()

		case <-interrupt:

			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}

			conn.Close()
			return
		}
	}

}
