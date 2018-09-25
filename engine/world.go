package engine

import "fmt"

type World struct {
    Chunks  *SparseSpace
    GeneratorFn string
    Players map[string]Entity
    Age     Tick
    Content *Content
}

type Tick uint


func MakeWorld (content *Content) *World {
    space := NewSparseSpace()
    players := make(map[string]Entity, 0)
    return &World{Chunks: space, Players: players, Age: 0, Content: content}
}


func (world *World) GetChunk (p Point) (*Chunk) {
    generic := world.Chunks.Get(p)
    chunk, ok := (generic).(*Chunk)
    if !ok {
        panic(fmt.Sprintf(`Expected type to be Chunk but it was not: %v`, generic))
    }
    return chunk
}

func (world *World) GenerateChunk (p Point) {
    world.Chunks.Set(p, nil)
}


func (world *World) GetVoxel (p AbsolutePoint) (*Voxel) {
    chunk := world.GetChunk(p.ChunkCoords)
    return chunk.Get(p.VoxelCoords)
}

func (world *World) SetVoxel (p AbsolutePoint, v Voxel) {
    chunk := world.GetChunk(p.ChunkCoords)
    chunk.Set(p.VoxelCoords, v)
}
