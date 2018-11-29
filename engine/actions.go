package main

import (
	"github.com/felzix/huyilla/types"
)

func (engine *Engine) RegisterAction(action *types.Action) {
	engine.Actions = append(engine.Actions, action)
}

// returns true if move succeeded; false otherwise
func (engine *Engine) move(action *types.Action) (bool, error) {
	player, err := engine.GetPlayer(action.PlayerName)
	if err != nil {
		return false, err
	}

	if player.Entity == nil {
		return false, nil // player doesn't have an entity (player has not yet finished signup)
	}

	err = engine.removeEntityFromChunk(player.Entity)
	if err != nil {
		return false, err
	}

	player.Entity.Location = action.GetMove().WhereTo

	err = engine.addEntityToChunk(player.Entity)
	if err != nil {
		return false, err
	}

	return true, nil
}
