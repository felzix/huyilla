package client

import (
	"github.com/felzix/huyilla/constants"
	"github.com/felzix/huyilla/types"
	"sync"
)

type WorldCache struct {
	sync.Mutex

	age    uint64
	chunks map[types.ComparablePoint]*types.DetailedChunk
}

func NewWorldCache() *WorldCache {
	return &WorldCache{
		age:    0,
		chunks: make(map[types.ComparablePoint]*types.DetailedChunk, constants.ACTIVE_CHUNK_CUBE),
	}
}

func (cache *WorldCache) GetAge() uint64 {
	return cache.age
}

func (cache *WorldCache) SetAge(age uint64) {
	cache.Lock()
	defer cache.Unlock()
	cache.age = age
}

func (cache *WorldCache) GetChunk(coords *types.Point) *types.DetailedChunk {
	return cache.chunks[*types.NewComparablePoint(coords)]
}

func (cache *WorldCache) SetChunk(coords *types.Point, chunk *types.DetailedChunk) {
	cache.Lock()
	defer cache.Unlock()
	cache.chunks[*types.NewComparablePoint(coords)] = chunk
}
