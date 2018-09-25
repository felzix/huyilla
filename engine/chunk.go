package engine


type Chunk struct {
	voxels VoxelCube
}
type VoxelCube [CHUNK_SIZE][CHUNK_SIZE][CHUNK_SIZE]Voxel

func MakeChunk (voxels VoxelCube) *Chunk {
    return &Chunk{voxels}
}

func (chunk *Chunk) Get (p Point) (*Voxel) {
    return &chunk.voxels[p.Y][p.X][p.Z]
}


func (chunk *Chunk) Set (p Point, v Voxel) {
    chunk.voxels[p.Y][p.X][p.Z] = v
}
