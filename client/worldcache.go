package client

import (
	"github.com/felzix/huyilla/constants"
	"github.com/felzix/huyilla/types"
	"sync"
)

type WorldCache struct {
	sync.Mutex

	age            types.Age
	chunks         map[types.Point]*types.Chunk
	previousChunks map[types.Point]*types.Chunk
	entities map[types.EntityId]types.Entity
	items map[types.ItemId]types.Item
}

func NewWorldCache() *WorldCache {
	return &WorldCache{
		age:            0,
		chunks:         make(map[types.Point]*types.Chunk, constants.ACTIVE_CHUNK_CUBE),
		previousChunks: make(map[types.Point]*types.Chunk, constants.ACTIVE_CHUNK_CUBE),
	}
}

func (cache *WorldCache) GetAge() types.Age {
	return cache.age
}

func (cache *WorldCache) SetAge(age types.Age) {
	cache.Lock()
	defer cache.Unlock()
	cache.age = age
}

func (cache *WorldCache) GetChunk(point types.Point) *types.Chunk {
	return cache.chunks[point]
}

func (cache *WorldCache) GetPreviousChunk(point types.Point) *types.Chunk {
	return cache.previousChunks[point]
}

func (cache *WorldCache) SetChunk(point types.Point, chunk *types.Chunk) {
	cache.Lock()
	defer cache.Unlock()
	cache.previousChunks[point] = cache.chunks[point]
	cache.chunks[point] = chunk
}

func (cache *WorldCache) GetEntity(id types.EntityId) types.Entity {
	return cache.entities[id]
}