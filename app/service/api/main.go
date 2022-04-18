package main

import "github.com/Mahamadou828/tgs_with_golang/foundation/logger"

func main() {
	log, err := logger.NewLogger("tgs-api")

	if err != nil {
		panic(err)
	}

	log.Info("Testing the logger package")
}
