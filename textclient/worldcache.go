package main

import (
	"github.com/felzix/huyilla/constants"
	"github.com/felzix/huyilla/types"
)

type WorldCache struct {
	age    uint64
	chunks map[types.ComparablePoint]*types.ChunkDetail
}

func (world *WorldCache) Init() {
	world.age = 0
	world.chunks = make(map[types.ComparablePoint]*types.ChunkDetail, constants.ACTIVE_CHUNK_CUBE)
}

func (client *Client) SetChunk(coords *types.Point, chunk *types.ChunkDetail) {
	client.Lock()
	defer client.Unlock()

	client.world.chunks[*types.NewComparablePoint(coords)] = chunk
}
