package engine

type World struct {
    Chunks  *SparseSpace
    GeneratorFn MapGenerator
    GeneratorSeed uint64
    Players map[string]Entity
    Age     Tick
    Content *Content
}

type Tick uint


func NewWorld (content *Content, gen MapGenerator) *World {
    return &World{
        Chunks:      NewSparseSpace(),
        GeneratorFn: gen,
        Players:     make(map[string]Entity, 0),
        Age:         0,
        Content:     content}
}


func (world *World) GetChunk (p Point) (*Chunk) {
    generic := world.Chunks.Get(p)
    chunk, ok := (generic).(*Chunk)
    if !ok {
        return nil
    }
    return chunk
}

func (world *World) GenerateChunk (p Point) {
    c := GenerateChunk(world.Content, world.GeneratorFn, world.GeneratorSeed, p)
    world.Chunks.Set(p, c)
}


func (world *World) GetVoxel (p AbsolutePoint) (*Voxel) {
    chunk := world.GetChunk(p.ChunkCoords)
    if chunk == nil {
        return nil
    }
    return chunk.Get(p.VoxelCoords)
}

func (world *World) SetVoxel (p AbsolutePoint, v Voxel) {
    chunk := world.GetChunk(p.ChunkCoords)
    chunk.Set(p.VoxelCoords, v)
}
