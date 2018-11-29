package main

import (
	"fmt"
	C "github.com/felzix/huyilla/constants"
	"github.com/felzix/huyilla/types"
	"github.com/mitchellh/hashstructure"
	"github.com/pkg/errors"
	"math/rand"
)

func (engine *Engine) GetChunk(p *types.Point) (*types.Chunk, error) {
	if chunk, ok := engine.Chunks[*p]; ok {
		return chunk, nil
	} else {
		return nil, errors.New(fmt.Sprintf("Chunk (%d,%d,%d) has yet to be generated", p.X, p.Y, p.Z))
	}
}

func (engine *Engine) getChunkGuaranteed(p *types.Point) (*types.Chunk, error) {
	if chunk, err := engine.GetChunk(p); err == nil {
		return chunk, nil
	} else if chunk, err := engine.GenChunk(C.SEED, p); err == nil {
		return chunk, nil
	} else {
		return nil, errors.Wrap(err, "Chunk generation failure")
	}
}

func (engine *Engine) SetChunk(p *types.Point, chunk *types.Chunk) {
	engine.Chunks[*p] = chunk
}

func (engine *Engine) addEntityToChunk(entity *types.Entity) error {
	if chunk, err := engine.getChunkGuaranteed(entity.Location.Chunk); err == nil {
		chunk.Entities = append(chunk.Entities, entity.Id)
		return nil
	} else {
		return err
	}
}

func (engine *Engine) removeEntityFromChunk(entity *types.Entity) error {
	chunk, err := engine.GetChunk(entity.Location.Chunk)

	if err != nil {
		if err.Error() == "not found" {
			return nil // chunk doesn't exist anyway so it need not be changed
		} else {
			return err // something else went wrong
		}
	}

	entities := chunk.Entities
	for i := 0; i < len(entities); i++ {
		id := entities[i]
		if entity.Id == id {
			// idiomatic way of removing a list element in Go
			entities[i] = entities[len(entities)-1]
			entities = entities[:len(entities)-1]
			break
		}
	}

	chunk.Entities = append(chunk.Entities, entity.Id)
	engine.SetChunk(entity.Location.Chunk, chunk)
	return nil
}

func (engine *Engine) GenChunk(worldSeed uint64, p *types.Point) (*types.Chunk, error) {
	chunkSeed, _ := hashstructure.Hash(p, nil)
	seed := int64(worldSeed * chunkSeed)

	chunk := types.Chunk{Voxels: make([]uint64, C.CHUNK_LENGTH)}
	var x, y, z int64
	for x = 0; x < C.CHUNK_SIZE; x++ {
		for y = 0; y < C.CHUNK_SIZE; y++ {
			for z = 0; z < C.CHUNK_SIZE; z++ {
				rand.Seed(seed) // so voxels can use randomness
				index := (x * C.CHUNK_SIZE * C.CHUNK_SIZE) + (y * C.CHUNK_SIZE) + z
				location := &types.AbsolutePoint{
					Chunk: p,
					Voxel: &types.Point{x, y, z},
				}
				chunk.Voxels[index] = genVoxel(location)
			}
		}
	}

	engine.SetChunk(p, &chunk)
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
