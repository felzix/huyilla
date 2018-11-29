package main

import (
	"github.com/felzix/huyilla/types"
	"math/rand"
)

func (engine *Engine) newEntity(typeInt uint64, playerName string, location *types.AbsolutePoint) *types.Entity {
	entity := types.Entity{
		Id:       engine.genUniqueEntityId(),
		Type:     typeInt,
		Location: location,
	}
	if playerName == "" {
		entity.Control = types.Entity_NPC
	} else {
		entity.Control = types.Entity_PLAYER
		entity.PlayerName = playerName
	}

	engine.Entities[entity.Id] = &entity

	return &entity
}

func (engine *Engine) setEntity(entity *types.Entity) {
	engine.Entities[entity.Id] = entity
}

func (engine *Engine) getSpecificEntities(ids []int64) []*types.Entity {
	entities := make([]*types.Entity, len(ids))
	for i := 0; i < len(ids); i++ {
		entities[i] = engine.Entities[ids[i]]
	}
	return entities
}

func (engine *Engine) genUniqueEntityId() int64 {
	var id int64
	for {
		id = rand.Int63()
		if _, ok := engine.Entities[id]; !ok {
			break
		}
	}
	return id
}
