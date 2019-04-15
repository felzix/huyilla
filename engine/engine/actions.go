package engine

import (
	"fmt"
	"github.com/felzix/huyilla/types"
	"github.com/pkg/errors"
)

func (engine *Engine) RegisterAction(action *types.Action) {
	engine.Lock()
	defer engine.Unlock()

	engine.Actions = append(engine.Actions, action)
}

// returns true if move succeeded; false otherwise
func (engine *Engine) move(action *types.Action) (bool, error) {
	player, err := engine.World.Player(action.PlayerName)
	if err != nil {
		return false, err
	}

	if player.EntityId == 0 {
		return false, errors.New("player doesn't have an entity (player has not yet finished signup)")
	}

	entity, err := engine.World.Entity(player.EntityId)
	if err != nil {
		return false, err
	} else if entity == nil {
		return false, errors.New(fmt.Sprintf(`Entity "%d" does not exist`, player.EntityId))
	}

	if err := engine.World.RemoveEntityFromChunk(player.EntityId, entity.Location.Chunk); err != nil {
		return false, err
	}

	entity.Location = action.GetMove().WhereTo
	if err := engine.World.SetEntity(entity.Id, entity); err != nil {
		return false, err
	}

	if err := engine.World.AddEntityToChunk(entity); err != nil {
		return false, err
	}

	return true, nil
}
