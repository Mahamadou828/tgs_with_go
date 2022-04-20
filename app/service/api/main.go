package main

import (
	"fmt"
	"github.com/Mahamadou828/tgs_with_golang/foundation/logger"
	"os"
)

func main() {
	_, err := logger.NewLogger("tgs-api")

	n, err := os.Stdin.Write([]byte("Try to write to stdin"))

	if err != nil {
		panic(err)
	}

	fmt.Println(n)

	if err != nil {
		panic(err)
	}

	//log.Info("Testing the logger package")
}
