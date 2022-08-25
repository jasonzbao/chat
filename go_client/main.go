package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	Red   = "\033[31m"
	Green = "\033[32m"
	Reset = "\033[0m"
)

// StringPrompt asks for a string value using the label
func StringPrompt(label string) string {
	var s string
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stderr, label+" ")
		s, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}
	return strings.TrimSpace(s)
}

type Message struct {
	Message string `json:"message"`
}

// TODO: repackage so that api structs are separate package
type HistoryMessage struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`

	Contents string `json:"content"`
	Name     string `json:"name"`
}

type History struct {
	Messages []*HistoryMessage `json:"past_messages"`
}

func main() {
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodGet, "http://localhost:9001/history", nil)
	if err != nil {
		log.Fatalf("Error composing request: %v", err)
	}

	resp, err := client.Do(request)
	if err != nil {
		log.Fatalf("Error communicating with server: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Error received from server: %s", resp)
	}
	defer resp.Body.Close()

	var historyMessages History
	if err := json.NewDecoder(resp.Body).Decode(&historyMessages); err != nil {
		log.Fatalf("Error decoding message from server: %s", resp)
	}

	for i := len(historyMessages.Messages) - 1; i >= 0; i-- {
		msg := historyMessages.Messages[i]
		fmt.Printf("%s%s%s: %s \n", Red, msg.Name, Reset, msg.Contents)
	}

	ws, _, err := websocket.DefaultDialer.Dial("ws://localhost:9001/socket", nil)
	if err != nil {
		log.Fatalf("Error dialing : %v", err)
	}
	defer ws.Close()

	wg := &sync.WaitGroup{}
	defer wg.Wait()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				_, msg, err := ws.ReadMessage()
				if err != nil {
					fmt.Println("Error receiving message from socket %v", err)
					cancel()
					return
				}
				fmt.Println(string(msg))
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			message := StringPrompt("")
			msg := &Message{
				Message: message,
			}
			jsonBody, err := json.Marshal(msg)
			if err != nil {
				fmt.Println("Error with marshalling json body")
				continue
			}

			if err := ws.WriteMessage(websocket.TextMessage, jsonBody); err != nil {
				fmt.Println("Error sending websocket message")
			}
			if message == "/exit" {
				return
			}
		}
	}
}
