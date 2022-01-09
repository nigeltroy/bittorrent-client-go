package main

import (
	"log"
	"os"
)

const logFileName string = "logs.txt"

func main() {
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)

	if err != nil {
		log.Fatal(err)
	}

	defer logFile.Close()
	log.SetOutput(logFile)
	log.Printf("\nNew program run\n\n")
}
