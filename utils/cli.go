package cli

import (
	"fmt"
	"log"
	"strings"
)

// Askforconfirmation prompts the user to confirm their action
func Askforconfirmation() bool {
	var response string

	_, err := fmt.Scanln(&response)
	if err != nil {
		log.Fatal(err)
	}

	switch strings.ToLower(response) {
	case "y", "yes":
		return true
	case "n", "no":
		return false
	default:
		fmt.Println("I'm sorry but I didn't get what you meant, please type (y)es or (n)o and then press enter:")
		return Askforconfirmation()
	}
}
