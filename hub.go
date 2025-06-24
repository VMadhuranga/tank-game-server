package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type hub struct {
	players    map[*player]bool
	evHandlers map[string]func(ev event, c *player) error
	mu         sync.RWMutex
}

func (h *hub) setupEventHandlers() {
	h.evHandlers[EventNewPlayer] = handleNewPlayerEvent
	h.evHandlers[EventOtherPlayers] = handleOtherPlayersEvent
	h.evHandlers[EventMovePlayer] = handleMovePlayerEvent
	h.evHandlers[EventBulletHit] = handleBulletHitEvent
	h.evHandlers[EventShoot] = handleShootEvent
}

func (h *hub) handleEvents(ev event, p *player) error {
	evHandler, ok := h.evHandlers[ev.Type]
	if !ok {
		return fmt.Errorf("unsupported event type")
	}

	err := evHandler(ev, p)
	if err != nil {
		log.Printf("error handling event: %v\n", err)
		return err
	}

	return nil
}

func (h *hub) addToPlayers(p *player) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.players[p] = true
}

func (h *hub) removeFromPlayers(p *player) {
	h.mu.Lock()
	defer h.mu.Unlock()

	_, ok := h.players[p]
	if ok {
		p.connection.Close()
		close(p.egress)
		delete(h.players, p)
	}
}

func (h *hub) getOtherPlayers(p *player) []*player {
	players := []*player{}
	for player := range p.hub.players {
		if player.ID != p.ID {
			players = append(players, player)
		}
	}

	return players
}

func (h *hub) handlePlayable(w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(map[string]bool{
		"playable": len(h.players) <= 50,
	})
	if err != nil {
		log.Printf("error marshaling players length")
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(data)
}

func (h *hub) serveWS(w http.ResponseWriter, r *http.Request) {
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error upgrading to websocket protocol: %v\n", err)
		return
	}

	c := newPlayer(conn, h)
	h.addToPlayers(c)

	go c.readMessages()
	go c.writeMessages()
}

func newHub() *hub {
	h := &hub{
		players:    make(map[*player]bool),
		evHandlers: make(map[string]func(ev event, c *player) error),
		mu:         sync.RWMutex{},
	}
	h.setupEventHandlers()
	return h
}
