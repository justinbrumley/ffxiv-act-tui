package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

var logFile *os.File

func init() {
	// Log file
	var err error
	logFile, err = os.OpenFile("event-log.txt", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func logMessage(msg *Payload) error {
	str, err := json.MarshalIndent(msg, "", "  ")
	if err != nil {
		return err
	}

	logFile.WriteString(fmt.Sprintf("%v\n", string(str)))
	return nil
}
