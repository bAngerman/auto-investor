package main

import (
	"log"
	"os"

	"github.com/bAngerman/auto-investor/pkg/ndaxapi"
	"github.com/bAngerman/auto-investor/pkg/strategy/livetrade"
	"github.com/bAngerman/auto-investor/pkg/strategy/papertrade"
	cli "github.com/bAngerman/auto-investor/utils"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Panic("Error parsing .env file:", err)
	}
}
func main() {
	log.Println("Starting auto trader.")

	_, err := ndaxapi.Start()

	if err != nil {
		log.Panic("Error connecting to NDAX.")
	}
	log.Println("Connected, and authenticated to NDAX.")

	if os.Getenv("ENV") == "live" {
		log.Println(".env states to use live API and real money, are you sure? (y/n)")
		if cli.Askforconfirmation() {
			livetrade.Start()
		}
	} else {
		papertrade.Start()
	}
}
