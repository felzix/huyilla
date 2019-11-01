package engine

import (
	"fmt"
	C "github.com/felzix/huyilla/constants"
	"github.com/felzix/huyilla/types"
	"github.com/mitchellh/hashstructure"
	"github.com/pkg/errors"
	"math/rand"
)

//
// Chunk
//

func chunkKey(p *types.Point) string {
	return fmt.Sprintf(`Chunk.%d.%d.%d`, p.X, p.Y, p.Z)
}

func (world *World) Chunk(p *types.Point) (*types.Chunk, error) {
	if chunk, err := world.OnlyGetChunk(p); chunk != nil {
		return chunk, nil
	} else if err == nil {
		if chunk, err := world.GenerateChunk(p); err == nil {
			return chunk, nil
		} else {
			return nil, errors.Wrap(err, "Chunk generation failure")
		}
	} else {
		return nil, err
	}
}

func (world *World) OnlyGetChunk(p *types.Point) (*types.Chunk, error) {
	var chunk types.Chunk
	err := world.DB.Get(chunkKey(p), &chunk)
	switch err.(type) {
	case nil: // chunk found
		return &chunk, nil
	case ThingNotFoundError: // chunk not found
		return nil, nil
	default: // something went wrong
		return nil, err
	}
}

func (world *World) CreateChunk(p *types.Point, chunk *types.Chunk) error {
	return world.SetChunk(p, chunk)
}

func (world *World) SetChunk(p *types.Point, chunk *types.Chunk) error {
	age, err := world.Age()
	if err != nil {
		return err
	}

	chunk.Tick = age.Ticks

	return world.DB.Set(chunkKey(p), chunk)
}

func (world *World) DeleteChunk(p *types.Point) error {
	return world.DB.End(chunkKey(p))
}

func (world *World) AddEntityToChunk(entity *types.Entity) error {
	if chunk, err := world.Chunk(entity.Location.Chunk); err == nil {
		chunk.Entities = append(chunk.Entities, entity.Id)
		return world.SetChunk(entity.Location.Chunk, chunk)
	} else {
		return err
	}
}

func (world *World) RemoveEntityFromChunk(entityId int64, p *types.Point) error {
	chunk, err := world.Chunk(p)
	if err != nil {
		return err
	}

	entities := chunk.Entities
	for i := 0; i < len(entities); i++ {
		id := entities[i]
		if entityId == id {
			// idiomatic way of removing a list element in Go
			entities[i] = entities[len(entities)-1]
			chunk.Entities = entities[:len(entities)-1]
			break
		}
	}

	return world.SetChunk(p, chunk)
}

func (world *World) GenerateChunk(p *types.Point) (*types.Chunk, error) {
	if chunkSeed, err := hashstructure.Hash(p, nil); err == nil {
		rand.Seed(int64(world.Seed * chunkSeed))
	} else {
		return nil, err
	}

	world.WorldGenerator.SetupForChunk(p)

	chunk := types.NewChunk(0, C.CHUNK_LENGTH)
	var x, y, z int64
	for x = 0; x < C.CHUNK_SIZE; x++ {
		for y = 0; y < C.CHUNK_SIZE; y++ {
			for z = 0; z < C.CHUNK_SIZE; z++ {
				index := (x * C.CHUNK_SIZE * C.CHUNK_SIZE) + (y * C.CHUNK_SIZE) + z
				location := types.NewAbsolutePoint(p.X, p.Y, p.Z, x, y, z)
				chunk.Voxels[index] = uint64(world.WorldGenerator.GenVoxel(location))
			}
		}
	}

	if err := world.SetChunk(p, chunk); err != nil {
		return nil, err
	}
	return chunk, nil
}

// All methods are idempotent
type WorldGenerator interface {
	SetupForWorld()
	SetupForChunk(chunkLocation *types.Point)
	GenVoxel(voxelLocation *types.AbsolutePoint) types.Voxel
}
