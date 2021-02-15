package ndaxapi

import (
	"encoding/json"
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

func authenticate(c *websocket.Conn) {
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
	}

	sToken = authenticate2FA(c, res)
}

func authenticate2FA(c *websocket.Conn, res Payload) string {
	time := time.Now()
	key, err := totp.GenerateCode(twoFA, time)

	if err != nil {
		discord.Message("Bad 2FA generation: " + err.Error())
	}

	p := Payload{
		"Code": key,
	}

	r := Data{
		Type:     0,
		Sequence: getSequence(),
		Function: "Authenticate2FA",
	}

	res = sendRequest(c, r, p)

	if res["errormsg"] != nil {
		discord.Message("Bad 2FA request: " + res["errormsg"].(string))
	}

	return res["SessionToken"].(string)
}

func ping(c *websocket.Conn) {

	p := Payload{}
	r := Data{
		Type:     0,
		Sequence: getSequence(),
		Function: "Ping",
	}

	res := sendRequest(c, r, p)

	log.Println("Ping:", res)
}

// Start the ndaxio websocket process
func Start() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{
		Scheme: "wss",
		Host:   wsHost,
		Path:   wsPath,
	}
	log.Printf("connecting to %s", u)

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)

	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	authenticate(c)

	return

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			t, message, err := c.ReadMessage()

			if err != nil {
				log.Println("read:", err)
				return
			}

			log.Printf("type: %d", t)
			log.Printf("recv: %s", message)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {

		case <-done:
			return

		case <-ticker.C:
			ping(c)

		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
