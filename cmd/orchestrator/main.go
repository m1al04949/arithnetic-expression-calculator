package main

import (
	"log"

	"github.com/m1al04949/arithnetic-expression-calculator/internal/orchestrator"
)

func main() {

	if err := orchestrator.RunServer(); err != nil {
		log.Fatal(err)
	}

}
