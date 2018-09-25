package engine


type World struct {
    Chunks  SparseSpace
    Players map[string]Entity
    Age     Tick
    Content *Content
}

type Tick uint


func (world *World) GetVoxel (p AbsolutePoint) (*Voxel) {
    generic := world.Chunks.Get(p.chunkCoords)
    chunk, ok := generic.(Chunk)
    if !ok {
        panic(`Expected type to be Chunk but it wasn't.'`)
    }
    return chunk.Get(p.voxelCoords)
}

func (world *World) SetVoxel (p AbsolutePoint, v Voxel) {
    generic := world.Chunks.Get(p.chunkCoords)
    chunk, ok := generic.(Chunk)
    if !ok {
        panic(`Expected type to be Chunk but it wasn't.'`)
    }
    chunk.Set(p.voxelCoords, v)
}