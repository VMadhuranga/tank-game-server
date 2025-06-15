package main

import (
	"encoding/json"
	"log"
)

const (
	EventNewPlayer    = "new_player"
	EventAllPlayers   = "all_players"
	EventRemovePlayer = "remove_player"
)

type event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

func handleAllPlayersEvent(ev event, p *player) error {
	players := p.hub.getAllPlayers()
	data, err := json.Marshal(players)
	if err != nil {
		log.Printf("error marshaling players: %v\n", err)
		return err
	}

	outgoingEV := event{
		Type:    EventAllPlayers,
		Payload: data,
	}

	p.egress <- outgoingEV

	return nil
}

func handleNewPlayerEvent(ev event, p *player) error {
	data, err := json.Marshal(p)
	if err != nil {
		log.Printf("error marshaling player: %v\n", err)
		return err
	}

	outgoingEV := event{
		Type:    EventNewPlayer,
		Payload: data,
	}

	for player := range p.hub.players {
		player.egress <- outgoingEV
	}

	err = handleAllPlayersEvent(ev, p)
	if err != nil {
		log.Printf("error handling all players event: %v", err)
		return err
	}

	return nil
}

func handleRemovePlayerEvent(ev event, p *player) error {
	data, err := json.Marshal(p)
	if err != nil {
		log.Printf("error marshaling player: %v\n", err)
	}

	outgoingEV := event{
		Type:    EventRemovePlayer,
		Payload: data,
	}
	for player := range p.hub.players {
		player.egress <- outgoingEV
	}

	return nil
}
