package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bAngerman/auto-investor/pkg/strategy/livetrade"
	"github.com/bAngerman/auto-investor/pkg/strategy/papertrade"
	"github.com/bAngerman/auto-investor/pkg/taskrunner"
	"github.com/joho/godotenv"
	"github.com/manifoldco/promptui"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Panic("Error parsing .env file:", err)
	}
}

func main() {
	log.Println("Starting auto trader.")

	conn, err := taskrunner.Start()

	if err != nil {
		log.Panic("Error connecting to NDAX:", err)
	}

	log.Println("Connected and authenticated to NDAX.")

	if os.Getenv("ENV") == "live" {
		prompt := promptui.Prompt{
			Label:     ".env states to use live API and real money, are you sure? (y/n)",
			IsConfirm: true,
		}

		result, err := prompt.Run()

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		if result == "y" {
			livetrade.Start(conn)
		}
	} else {
		papertrade.Start(conn)
	}
}
