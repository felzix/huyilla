package engine

import (
    "github.com/mitchellh/hashstructure"
    "math/rand"
)

type MapGenerator interface {
    GenerateVoxel(*Content, AbsolutePoint) Voxel
}


func GenerateChunk (content *Content, gen MapGenerator, seed uint64, p Point) *Chunk {
    chunkSpecificSeed, _ := hashstructure.Hash(p, nil)

    var voxels [CHUNK_SIZE][CHUNK_SIZE][CHUNK_SIZE]Voxel

    for y := 0; y < CHUNK_SIZE; y++ {
        for x := 0; x < CHUNK_SIZE; x++ {
            for z := 0; z < CHUNK_SIZE; z++ {
                rand.Seed(int64(seed * chunkSpecificSeed))  // so voxels can use randomness
                voxels[y][x][z] = gen.GenerateVoxel(
                    content,
                    AbsolutePoint{p, Point{x, y, z}})
            }
        }
    }

    return MakeChunk(voxels)
}


type PondWorld struct {
    Seed uint
    pondSize uint
}

func NewPondWorldGenerator (seed uint, pondSize uint) *PondWorld {
    return &PondWorld{seed, pondSize}
}

func (gen *PondWorld) GenerateVoxel(content *Content, p AbsolutePoint) Voxel {
    v := content.V

    if p.ChunkCoords.Z < 0 {
        return NewVoxel(v["dirt"])
    }

    if p.ChunkCoords.Z > 0 {
        return NewVoxel(v["air"])
    }

    center := RandomPoint(16)
    center.Z = 0

    d := p.VoxelCoords.Distance(center)
    if p.VoxelCoords.Z == center.Z && d <= float64(gen.pondSize) {
        return NewVoxel(v["water"])
    }
    if p.VoxelCoords.Z == 0 {
        return NewVoxel(v["barren_earth"])
    }
    return NewVoxel(v["air"])
}
