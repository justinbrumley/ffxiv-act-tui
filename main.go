package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	ws "github.com/gorilla/websocket"
)

type Message map[string]interface{}

type Payload struct {
	Type        string  `json:"type"`
	MessageType string  `json:"msgtype"`
	Message     Message `json:"msg"`
}

var socketUrl string

func init() {
	fmt.Println("Setting up vars...")

	socketUrl = os.Getenv("SOCKET_URL")
	if socketUrl == "" {
		socketUrl = "ws://localhost:10501/MiniParse"
	}
}

func main() {
	defer logFile.Close()

	fmt.Println("Connecting to ACT socket server...")

	headers := http.Header{}
	dialer := ws.Dialer{}
	conn, resp, err := dialer.Dial(socketUrl, headers)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	if resp.StatusCode != 101 {
		log.Fatal("Did not receive 101 status code when dialing socket server")
	}

	for {
		msgType, reader, err := conn.NextReader()
		if err != nil {
			log.Fatal(err)
		}

		switch msgType {
		case ws.TextMessage:
			payload := &Payload{}
			decoder := json.NewDecoder(reader)
			if err := decoder.Decode(&payload); err != nil {
				continue
			}

			logMessage(payload)
			handleMessage(payload)
		}
	}
}
