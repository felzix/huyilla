package main

import (
	"fmt"
	C "github.com/felzix/huyilla/constants"
	"github.com/felzix/huyilla/types"
	"github.com/pkg/errors"
)

func (engine *Engine) Tick() error {
	players, err := engine.getActivePlayers()
	if err != nil {
		return errors.Wrap(err, "get active players")
	}

	activeChunks := make(map[types.Point]*types.Chunk, len(players)*C.ACTIVE_CHUNK_CUBE)
	for i := 0; i < len(players); i++ {
		player := players[i]
		loc := player.Entity.Location.Chunk
		for x := loc.X - C.ACTIVE_CHUNK_RADIUS; x < 1+loc.X+C.ACTIVE_CHUNK_RADIUS; x++ {
			for y := loc.Y - C.ACTIVE_CHUNK_RADIUS; y < 1+loc.Y+C.ACTIVE_CHUNK_RADIUS; y++ {
				for z := loc.Z - C.ACTIVE_CHUNK_RADIUS; z < 1+loc.Z+C.ACTIVE_CHUNK_RADIUS; z++ {
					point := newPoint(x, y, z)
					if chunk, err := engine.getChunkGuaranteed(point); err == nil {
						activeChunks[*point] = chunk
					} else {
						return errors.Wrap(err, "failed to get/gen chunk")
					}
				}
			}
		}
	}

	vitalizedVoxels := make([]types.Point, C.PASSIVE_VITALITY)
	for i := 0; i < C.PASSIVE_VITALITY; i++ {
		vitalizedVoxels[i] = *randomPoint()
	}

	for p, chunk := range activeChunks {
		for i := 0; i < len(chunk.ActiveVoxels); i++ {
			point := types.AbsolutePoint{Chunk: &p, Voxel: chunk.ActiveVoxels[i]}
			if err := engine.voxelPhysics(chunk, &point); err != nil {
				return errors.Wrap(err, "voxel physics of active voxels")
			}
		}

		for i := 0; i < C.PASSIVE_VITALITY; i++ {
			point := types.AbsolutePoint{Chunk: &p, Voxel: &vitalizedVoxels[i]}
			if err := engine.voxelPhysics(chunk, &point); err != nil {
				return errors.Wrap(err, "voxel physics of random voxels")
			}
		}

		for i := 0; i < len(chunk.Entities); i++ {
			entity := engine.Entities[chunk.Entities[i]]
			if err := engine.entityPhysics(chunk, entity); err != nil {
				return errors.Wrap(err, "entity physics")
			}
		}
	}

	for i := 0; i < len(engine.Actions); i++ {
		action := engine.Actions[i]

		var fn func(*types.Action) (bool, error)

		switch a := action.Action.(type) {
		case *types.Action_Move:
			fn = engine.move
		default:
			// only log error - if the action is broken then don't block the engine
			return errors.New(fmt.Sprintf("Invalid action %v", a))
		}

		if success, err := fn(action); success {
			// TODO success no error
		} else if err == nil {
			// TODO failure no error
		} else {
			return errors.Wrap(err, "action failure")
		}
	}

	// clear action queue
	engine.Actions = make([]*types.Action, 0)

	// save chunks
	for p, chunk := range activeChunks {
		engine.SetChunk(&p, chunk)
	}

	// advance age by one tick
	engine.Age++

	return nil
}

func (engine *Engine) voxelPhysics(chunk *types.Chunk, location *types.AbsolutePoint) error {
	return nil
}

func (engine *Engine) entityPhysics(chunk *types.Chunk, entity *types.Entity) error {
	return nil
}
