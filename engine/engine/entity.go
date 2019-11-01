package engine

import (
	"fmt"
	"github.com/felzix/huyilla/types"
	"math/rand"
)

func entityKey(id int64) string {
	return fmt.Sprintf(`Entity.%d`, id)
}

func (world *World) Entity(id int64) (*types.Entity, error) {
	var entity types.Entity
	err := world.DB.Get(entityKey(id), &entity)
	switch err.(type) {
	case nil: // entity found
		return &entity, nil
	case ThingNotFoundError: // entity not found
		return nil, nil
	default: // something went wrong
		return nil, err
	}
}

func (world *World) CreateEntity(typeInt uint64, playerName string, location *types.AbsolutePoint) (*types.Entity, error) {
	entity := types.Entity{
		Id:       world.genUniqueEntityId(),
		Type:     typeInt,
		Location: location,
	}
	if playerName == "" {
		entity.Control = types.Entity_NPC
	} else {
		entity.Control = types.Entity_PLAYER
		entity.PlayerName = playerName
	}

	if err := world.DB.Set(entityKey(entity.Id), &entity); err == nil {
		return &entity, nil
	} else {
		return nil, err
	}
}

func (world *World) SetEntity(id int64, entity *types.Entity) error {
	return world.DB.Set(entityKey(id), entity)
}

func (world *World) DeleteEntity(id int64) error {
	return world.DB.End(entityKey(id))
}

func (world *World) EntityExists(id int64) bool {
	return world.DB.Has(entityKey(id))
}

func (world *World) genUniqueEntityId() int64 {
	var id int64
	for {
		id = rand.Int63()
		if !world.EntityExists(id) {
			break
		}
	}
	return id
}
