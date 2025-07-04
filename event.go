package main

import (
	"encoding/json"
	"log"
)

const (
	EventNewPlayer    = "new_player"
	EventOtherPlayers = "other_players"
	EventRemovePlayer = "remove_player"
	EventMovePlayer   = "move_player"
	EventBulletHit    = "bullet_hit"
	EventShoot        = "shoot"
)

type event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
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

	return nil
}

func handleOtherPlayersEvent(ev event, p *player) error {
	players := p.hub.getOtherPlayers(p)
	data, err := json.Marshal(players)
	if err != nil {
		log.Printf("error marshaling players: %v\n", err)
		return err
	}

	outgoingEV := event{
		Type:    EventOtherPlayers,
		Payload: data,
	}
	p.egress <- outgoingEV

	return nil
}

func handleRemovePlayerEvent(ev event, p *player) error {
	data, err := json.Marshal(p)
	if err != nil {
		log.Printf("error marshaling player: %v\n", err)
		return err
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

func handleMovePlayerEvent(ev event, p *player) error {
	pld := player{}
	err := json.Unmarshal(ev.Payload, &pld)
	if err != nil {
		log.Printf("error unmarshaling payload: %v\n", err)
		return err
	}

	p.PX = pld.PX
	p.PY = pld.PY
	p.Angle = pld.Angle

	data, err := json.Marshal(p)
	if err != nil {
		log.Printf("error marshaling player: %v\n", err)
	}

	outgoingEV := event{
		Type:    EventMovePlayer,
		Payload: data,
	}
	for player := range p.hub.players {
		if player.ID != p.ID {
			player.egress <- outgoingEV
		}
	}

	return nil
}

func handleBulletHitEvent(ev event, p *player) error {
	outgoingEV := event{
		Type:    EventBulletHit,
		Payload: ev.Payload,
	}
	for player := range p.hub.players {
		if player.ID != p.ID {
			player.egress <- outgoingEV
		}
	}

	return nil
}

func handleShootEvent(ev event, p *player) error {
	outgoingEV := event{
		Type:    EventShoot,
		Payload: ev.Payload,
	}
	for player := range p.hub.players {
		if player.ID != p.ID {
			player.egress <- outgoingEV
		}
	}

	return nil
}
