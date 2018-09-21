package engine


type Chunk struct {
	voxels [][][]Voxel
}

func (chunk *Chunk) Get (p Point) (*Voxel) {
    return &chunk.voxels[p.y][p.x][p.z]
}


func (chunk *Chunk) Set (p Point, v Voxel) {
    chunk.voxels[p.y][p.x][p.z] = v
}
