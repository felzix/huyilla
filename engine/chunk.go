package main

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
	if err := gettum(world, chunkKey(p), &chunk); err == nil {
		return &chunk, nil
	} else if fileIsNotFound(err) {
		return nil, nil
	} else {
		return nil, err
	}
}

func (world *World) CreateChunk(p *types.Point, chunk *types.Chunk) error {
	return settum(world, chunkKey(p), chunk)
}

func (world *World) SetChunk(p *types.Point, chunk *types.Chunk) error {
	return world.CreateChunk(p, chunk)
}

func (world *World) DeleteChunk(p *types.Point) error {
	return enddum(world, chunkKey(p))
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
	chunkSeed, _ := hashstructure.Hash(p, nil)
	seed := int64(world.Seed * chunkSeed)

	chunk := types.Chunk{Voxels: make([]uint64, C.CHUNK_LENGTH)}
	var x, y, z int64
	for x = 0; x < C.CHUNK_SIZE; x++ {
		for y = 0; y < C.CHUNK_SIZE; y++ {
			for z = 0; z < C.CHUNK_SIZE; z++ {
				rand.Seed(seed) // so voxels can use randomness
				index := (x * C.CHUNK_SIZE * C.CHUNK_SIZE) + (y * C.CHUNK_SIZE) + z
				location := &types.AbsolutePoint{
					Chunk: p,
					Voxel: &types.Point{X: x, Y: y, Z: z},
				}
				chunk.Voxels[index] = genVoxel(location)
			}
		}
	}

	if err := world.SetChunk(p, &chunk); err != nil {
		return nil, err
	}
	return &chunk, nil
}

func genVoxel(p *types.AbsolutePoint) uint64 {
	v := VOXEL

	if p.Chunk.Z < 0 {
		return v["dirt"]
	}

	if p.Chunk.Z > 0 {
		return v["air"]
	}

	center := randomPoint()
	center.Z = 0

	d := distance(p.Voxel, center)
	if p.Voxel.Z == center.Z && d <= float64(3) {
		return v["water"]
	}
	if p.Voxel.Z == 0 {
		return v["barren_earth"]
	}
	return v["air"]
}
