package auth

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"time"

	"github.com/bAngerman/auto-investor/pkg/discord"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/pquerna/otp/totp"

	"github.com/bAngerman/auto-investor/pkg/ndaxapi"
)

// loginResponse is the response of the login request.
type loginResponse struct {
	Authenticated bool   `json:"Authenticated"` // false
	Requres2FA    bool   `json:"requires2FA"`   // true
	AuthType      string `json:"twoFAType"`     // "Google"
	AddtlInfo     string `json:"AddtlInfo"`     // Additional info
	ErrorMessage  string `json:"errormsg"`      // Error if any
}

type twoFAResponse struct {
	Authenticated bool   `json:"Authenticated"` // true
	UserID        int64  `json:"UserId"`        // 1
	SessionToken  string `json:"SessionToken"`  // "e4d..."
	ErrorMessage  string `json:"errormsg"`      // Error if any
}

var (
	username, password string // username, password from .env
	twoFA              string // Two factor token from .env
	// sToken             string // Session token, likely not to be used.

	// UserID used elsewhere
	UserID int64 // User id used elsewhere
)

func init() {
	godotenv.Load()

	username = os.Getenv("NDAXIO_USER")
	password = os.Getenv("NDAXIO_PASS")
	twoFA = os.Getenv("NDAXIO_2FA")
}

// Authenticate logs the user in.
func Authenticate(conn *websocket.Conn) error {
	err := login(conn)
	if err != nil {
		log.Panic(err)
	}
	err = authenticate2FA(conn)
	if err != nil {
		log.Panic(err)
	}

	return err
}

// login sends user name and password params to
// authenticate the user.
func login(conn *websocket.Conn) error {
	p := ndaxapi.Payload{
		"UserName": username,
		"Password": password,
	}

	r := ndaxapi.Request{
		Type:     0,
		Sequence: ndaxapi.GetSequence(),
		Function: "AuthenticateUser",
	}

	data := ndaxapi.Send(r, p, conn)

	var loginData loginResponse

	err := json.Unmarshal([]byte(data.Response), &loginData)

	if err != nil || loginData.ErrorMessage != "" {
		discord.Message("Bad ping: " + loginData.ErrorMessage)
		return errors.New("Error authenticating: " + loginData.ErrorMessage)
	}

	return nil
}

// authenticate2FA determines the 2FA token, and completes
// the authentication request.
func authenticate2FA(conn *websocket.Conn) error {
	time := time.Now()
	key, err := totp.GenerateCode(twoFA, time)

	if err != nil {
		discord.Message("Bad 2FA generation: " + err.Error())
		return err
	}

	p := ndaxapi.Payload{
		"Code": key,
	}

	r := ndaxapi.Request{
		Type:     0,
		Sequence: ndaxapi.GetSequence(),
		Function: "Authenticate2FA",
	}

	data := ndaxapi.Send(r, p, conn)

	var twoFAData twoFAResponse

	err = json.Unmarshal([]byte(data.Response), &twoFAData)

	if err != nil || twoFAData.ErrorMessage != "" {
		discord.Message("Error authenticating with 2FA: " + twoFAData.ErrorMessage)
		return errors.New("Error authenticating with 2FA: " + twoFAData.ErrorMessage)
	}

	// Store for use later.
	UserID = twoFAData.UserID

	return nil
}
