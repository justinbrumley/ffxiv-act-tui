package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	ws "github.com/gorilla/websocket"
)

type Message map[string]interface{}

type Payload struct {
	Type        string  `json:"type"`
	MessageType string  `json:"msgtype"`
	Message     Message `json:"msg"`
}

var tuiState State
var socketUrl string

func init() {
	socketUrl = os.Getenv("SOCKET_URL")
	if socketUrl == "" {
		socketUrl = "ws://localhost:10501/MiniParse"
	}
}

var conn *ws.Conn

func main() {
	defer logFile.Close()

	var err error
	dialer := ws.Dialer{}
	conn, _, err = dialer.Dial(socketUrl, http.Header{})
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	p := tea.NewProgram(NewState(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func (s *State) Init() tea.Cmd {
	return tick()
}

type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(time.Second/60, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Update responds to key presses, and keeps the state up-to-date.
func (s *State) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		msgType, reader, err := conn.NextReader()
		if err != nil {
			log.Fatal(err)
		}

		switch msgType {
		case ws.TextMessage:
			payload := &Payload{}
			decoder := json.NewDecoder(reader)
			if err := decoder.Decode(&payload); err != nil {
				return s, tick()
			}

			// logMessage(payload)
			s.handleMessage(payload)
		}

		return s, tick()

	case tea.KeyMsg:
		{
			switch msg.String() {
			case "ctrl+c", "q":
				{
					return s, tea.Quit
				}
			}
		}
	}

	return s, nil
}
