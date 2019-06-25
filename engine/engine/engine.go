package engine

import (
	"fmt"
	C "github.com/felzix/huyilla/constants"
	"github.com/felzix/huyilla/content"
	"github.com/felzix/huyilla/types"
	"github.com/pkg/errors"
	"sync"
)

type Engine struct {
	sync.Mutex

	World   *World
	Actions []*types.Action
	Secret  []byte
}

func (engine *Engine) Init(saveDir string) error {
	// So that recipes and terrain generator can reference content by name.
	content.PopulateContentNameMaps()

	engine.World = &World{Seed: C.SEED}

	if err := engine.World.Init(saveDir, 16*1024*1024); err != nil { // 16 MB
		return err
	}

	engine.Actions = make([]*types.Action, 0)

	// TODO not hard-coded
	engine.Secret = []byte(`&$0C-7#o4sK"W*&Q7;8PD_pz^8%]"v),zY(b-3.v`)

	return nil
}

func (engine *Engine) Tick() error {
	// advance age by one tick
	_, err := engine.World.IncrementAge()
	if err != nil {
		return err
	}

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
					point := types.NewPoint(x, y, z)
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
		vitalizedVoxels[i] = *types.RandomPoint(C.CHUNK_SIZE)
	}

	for cp, chunk := range activeChunks {
		p := types.NewPoint(cp.X, cp.Y, cp.Z)
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

	engine.Lock()
	actions := engine.Actions
	engine.Actions = make([]*types.Action, 0) // clear action queue so engine can unlock
	engine.Unlock()

	for _, action := range actions {
		var fn func(*types.Action) (bool, error)

		switch a := action.Action.(type) {
		case *types.Action_Move:
			fn = engine.Move
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

	return nil
}

func (engine *Engine) voxelPhysics(chunk *types.Chunk, location *types.AbsolutePoint) error {
	return nil
}

func (engine *Engine) entityPhysics(chunk *types.Chunk, entity *types.Entity) error {
	return nil
}
