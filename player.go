package main

import (
	"encoding/json"
	"log"
	"math/rand/v2"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var (
	readLimit    int64 = 512
	pongWait           = 10 * time.Second
	pingInterval       = (pongWait * 9) / 10
)

type player struct {
	connection *websocket.Conn
	hub        *hub
	egress     chan event
	ID         string `json:"id"`
	PX         int    `json:"pX"`
	PY         int    `json:"pY"`
}

func (p *player) readMessages() {
	defer func() {
		p.hub.removeFromPlayers(p)
	}()

	p.connection.SetReadLimit(readLimit)

	err := p.connection.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		log.Printf("error setting pong wait: %v\n", err)
		return
	}

	p.connection.SetPongHandler(func(appData string) error {
		log.Println("pong")
		return p.connection.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		_, payload, err := p.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v\n", err)
			}
			break
		}

		ev := event{}
		err = json.Unmarshal(payload, &ev)
		if err != nil {
			log.Printf("error unmarshaling payload: %v\n", err)
			break
		}

		err = p.hub.handleEvents(ev, p)
		if err != nil {
			log.Printf("error handling events: %v\n", err)
		}
	}
}

func (p *player) writeMessages() {
	ticker := time.NewTicker(pingInterval)

	defer func() {
		ticker.Stop()
		p.hub.removeFromPlayers(p)
	}()

	for {
		select {
		case ev, ok := <-p.egress:
			if !ok {
				err := p.connection.WriteMessage(websocket.CloseMessage, nil)
				if err != nil {
					log.Printf("error writing message: %v\n", err)
				}

				err = handleRemovePlayerEvent(ev, p)
				if err != nil {
					log.Printf("error handling remove player event: %v\n", err)
				}

				return
			}

			data, err := json.Marshal(ev)
			if err != nil {
				log.Printf("error marshaling message: %v", err)
				return
			}

			err = p.connection.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.Printf("error writing message: %v", err)
				return
			}
		case <-ticker.C:
			log.Println("ping")
			err := p.connection.WriteMessage(websocket.PingMessage, []byte{})
			if err != nil {
				log.Printf("error writing message: %v", err)
				return
			}
		}
	}
}

func newPlayer(conn *websocket.Conn, h *hub) *player {
	return &player{
		connection: conn,
		hub:        h,
		egress:     make(chan event),
		ID:         uuid.NewString(),
		PX:         rand.IntN(1024 - 5),
		PY:         rand.IntN(768 - 5),
	}
}
