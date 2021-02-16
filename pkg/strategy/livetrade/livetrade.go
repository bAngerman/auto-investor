package livetrade

import (
	"log"

	"github.com/bAngerman/auto-investor/pkg/ndaxapi"
)

// Start new live trade strategy
func Start() {
	log.Println("Using livetrade strategy")
	ndaxapi.Initloop(loop)
}

func loop() {
	log.Println("Running loop")
}
