package engine

import (
	C "github.com/felzix/huyilla/constants"
	"github.com/felzix/huyilla/types"
	"github.com/pkg/errors"
	"sync"
)

type Engine struct {
	sync.Mutex

	World   *types.World
	Actions []*types.Action
	Secret  []byte
}

func NewEngine(seed uint64, generator types.WorldGenerator, db types.Database) (*Engine, error) {
	world, err := types.NewWorld(seed, generator, db)
	if err != nil {
		return nil, err
	}
	return &Engine{
		World: world,
		Actions: make([]*types.Action, 0),
		Secret: []byte(`&$0C-7#o4sK"W*&Q7;8PD_pz^8%]"v),zY(b-3.v`), // TODO not hard-coded
	}, nil
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

	vitalizedVoxels := make([]types.Point, C.PASSIVE_VITALITY)
	for i := 0; i < C.PASSIVE_VITALITY; i++ {
		vitalizedVoxels[i] = types.RandomPoint(C.CHUNK_SIZE)
	}

	for i := 0; i < len(players); i++ {
		player := players[i]
		loc := player.Entity.Location.Chunk
		for x := loc.X - C.ACTIVE_CHUNK_RADIUS; x < 1+loc.X+C.ACTIVE_CHUNK_RADIUS; x++ {
			for y := loc.Y - C.ACTIVE_CHUNK_RADIUS; y < 1+loc.Y+C.ACTIVE_CHUNK_RADIUS; y++ {
				for z := loc.Z - C.ACTIVE_CHUNK_RADIUS; z < 1+loc.Z+C.ACTIVE_CHUNK_RADIUS; z++ {
					p := types.NewPoint(x, y, z)
					if chunk, err := engine.World.Chunk(p); err == nil {
						// voxel physics
						for i := 0; i < C.PASSIVE_VITALITY; i++ {
							point := types.AbsolutePoint{Chunk: p, Voxel: vitalizedVoxels[i]}
							if err := engine.voxelPhysics(chunk, point); err != nil {
								return errors.Wrap(err, "voxel physics of random voxels")
							}
						}

						// entity physics
						for _, id := range chunk.Entities {
							entity, err := engine.World.Entity(id)
							if err != nil {
								return err
							}
							if err := engine.entityPhysics(chunk, entity); err != nil {
								return errors.Wrap(err, "entity physics")
							}
						}
					} else {
						return errors.Wrap(err, "failed to get/gen chunk")
					}
				}
			}
		}
	}

	engine.Lock()
	actions := engine.Actions
	engine.Actions = make([]*types.Action, 0) // clear action queue so engine can unlock
	engine.Unlock()

	for _, action := range actions {
		if success, err := action.Apply(engine.World); success {
			Log.Debug("success no error")
		} else if err == nil {
			Log.Debug("failure no error")
		} else {
			return errors.Wrap(err, "action failure")
		}
	}

	return nil
}

func (engine *Engine) voxelPhysics(chunk *types.Chunk, location types.AbsolutePoint) error {
	return nil
}

func (engine *Engine) entityPhysics(chunk *types.Chunk, entity *types.Entity) error {
	return nil
}
