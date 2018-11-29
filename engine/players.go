package main

import (
	"fmt"
	"github.com/felzix/huyilla/types"
	"github.com/pkg/errors"
)

func (engine *Engine) GetPlayerList() []string {
	list := make([]string, len(engine.Players))

	i := 0
	for name, _ := range engine.Players {
		list[i] = name
		i++
	}

	return list
}

func (engine *Engine) GetPlayer(name string) (*types.PlayerDetails, error) {
	if player, ok := engine.Players[name]; ok {
		entity, _ := engine.Entities[player.EntityId]
		// NOTE: entity can be nil
		return &types.PlayerDetails{Player: player, Entity: entity}, nil
	} else {
		return nil, errors.New(fmt.Sprintf(`No such player "%s"`, name))
	}

}

func (engine *Engine) getActivePlayers() ([]*types.PlayerDetails, error) {
	var activePlayers []*types.PlayerDetails

	for _, player := range engine.Players {
		if player.LoggedIn {
			entity := engine.Entities[player.EntityId]
			activePlayers = append(activePlayers, &types.PlayerDetails{Player: player, Entity: entity})
		}
	}

	return activePlayers, nil
}
