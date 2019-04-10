package engine

import (
	"fmt"
	C "github.com/felzix/huyilla/constants"
	"github.com/felzix/huyilla/content"
	"github.com/felzix/huyilla/types"
	"github.com/pkg/errors"
)

type Engine struct {
	World   *World
	Actions []*types.Action // TODO locking
	Secret  []byte
}

func (engine *Engine) Init(saveDir string) error {
	// So that recipes and terrain generator can reference content by name.
	content.PopulateContentNameMaps()

	engine.World = &World{Seed: C.SEED}

	if err := engine.World.Init(saveDir, 1024*1024); err != nil { // 1 MB
		return err
	}

	engine.Actions = make([]*types.Action, 0)

	// TODO not hard-coded
	engine.Secret = []byte(`&$0C-7#o4sK"W*&Q7;8PD_pz^8%]"v),zY(b-3.v`)

	return nil
}

func (engine *Engine) Tick() error {
	players, err := engine.World.GetActivePlayers()
	if err != nil {
		return errors.Wrap(err, "get active players")
	}

	activeChunks := make(map[types.ComparablePoint]*types.Chunk, len(players)*C.ACTIVE_CHUNK_CUBE)
	for i := 0; i < len(players); i++ {
		player := players[i]
		loc := player.Entity.Location.Chunk
		for x := loc.X - C.ACTIVE_CHUNK_RADIUS; x < 1+loc.X+C.ACTIVE_CHUNK_RADIUS; x++ {
			for y := loc.Y - C.ACTIVE_CHUNK_RADIUS; y < 1+loc.Y+C.ACTIVE_CHUNK_RADIUS; y++ {
				for z := loc.Z - C.ACTIVE_CHUNK_RADIUS; z < 1+loc.Z+C.ACTIVE_CHUNK_RADIUS; z++ {
					point := NewPoint(x, y, z)
					if chunk, err := engine.World.Chunk(point); err == nil {
						activeChunks[*types.NewComparablePoint(point)] = chunk
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

	for cp, chunk := range activeChunks {
		p := NewPoint(cp.X, cp.Y, cp.Z)
		for i := 0; i < len(chunk.ActiveVoxels); i++ {
			point := types.AbsolutePoint{Chunk: p, Voxel: chunk.ActiveVoxels[i]}
			if err := engine.voxelPhysics(chunk, &point); err != nil {
				return errors.Wrap(err, "voxel physics of active voxels")
			}
		}

		for i := 0; i < C.PASSIVE_VITALITY; i++ {
			point := types.AbsolutePoint{Chunk: p, Voxel: &vitalizedVoxels[i]}
			if err := engine.voxelPhysics(chunk, &point); err != nil {
				return errors.Wrap(err, "voxel physics of random voxels")
			}
		}

		for i := 0; i < len(chunk.Entities); i++ {
			entity, err := engine.World.Entity(chunk.Entities[i])
			if err != nil {
				return err
			}
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
	for cp, chunk := range activeChunks {
		p := NewPoint(cp.X, cp.Y, cp.Z)
		if err := engine.World.SetChunk(p, chunk); err != nil {
			return err
		}
	}

	// advance age by one tick
	if err := engine.World.IncrementAge(); err != nil {
		return err
	}

	return nil
}

func (engine *Engine) voxelPhysics(chunk *types.Chunk, location *types.AbsolutePoint) error {
	return nil
}

func (engine *Engine) entityPhysics(chunk *types.Chunk, entity *types.Entity) error {
	return nil
}
