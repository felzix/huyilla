package main

import (
	"github.com/felzix/huyilla/constants"
	"github.com/felzix/huyilla/types"
)

type WorldCache struct {
	age      uint64
	chunks   map[types.Point]*types.Chunk
	entities map[int64]*types.Entity
}

func (world *WorldCache) Init() {
	world.age = 0
	world.chunks = make(map[types.Point]*types.Chunk, constants.ACTIVE_CHUNK_CUBE)
	world.entities = make(map[int64]*types.Entity)
}
