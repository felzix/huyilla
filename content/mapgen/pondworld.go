package mapgen

import (
    "github.com/felzix/huyilla-dappchain/engine"
)

type PondWorld struct {
    world *engine.World
    seed uint
    pondSize uint
}

func NewPondWorldGenerator (world *engine.World, seed uint, pondSize uint) *PondWorld {
    return &PondWorld{world, seed, pondSize}
}

// func GeneratorChunk(p engine.Point) *engine.Chunk {
//
// }
//
// func ( gen*PondWorld) GenerateVoxel(p engine.AbsolutePoint) engine.Voxel {
//     voxels := gen.world.Content.
//     if p.ChunkCoords.Z < 0 {
//         return engine.NewVoxel(gen.world.Content)
//     }
// }