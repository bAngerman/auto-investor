package taskrunner

import (
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"

	"github.com/bAngerman/auto-investor/pkg/ndaxapi"
	"github.com/bAngerman/auto-investor/pkg/ndaxapi/auth"
)

var (
	wsHost, wsPath string
)

func init() {
	godotenv.Load()

	wsHost = os.Getenv("NDAXIO_WS_HOST")
	wsPath = os.Getenv("NDAXIO_WS_PATH")
}

// Start the ndaxio websocket process
func Start() (*websocket.Conn, error) {
	u := url.URL{
		Scheme: "wss",
		Host:   wsHost,
		Path:   wsPath,
	}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)

	if err != nil {
		log.Fatal("Error connecting to ws:", err)
	}

	err = auth.Authenticate(conn)

	return conn, err
}

// Initloop the loop used for trading
func Initloop(conn *websocket.Conn, callback func()) {

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	done := make(chan struct{})

	pause := time.Second * 1

	ticker := time.NewTicker(pause)
	defer ticker.Stop()

	for {
		select {

		case <-done:
			return

		case <-ticker.C:

			// Always ping to keep conn alive.
			ndaxapi.Ping(conn)

			callback()

		case <-interrupt:

			log.Println("Shutting down.")

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
