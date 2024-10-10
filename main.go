package main

import (
	"log"
	"main/cmd"
	"os"
)

func main() {
	logFile, err := os.OpenFile("logs/go.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Error opening log file: ", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	log.Println("Starting Trigger CLI")
	cmd.Execute()
}