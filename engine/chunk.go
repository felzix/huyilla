package engine


type Chunk struct {
	voxels [][][]Voxel
}

func (chunk *Chunk) Get (p Point) (*Voxel) {
    return &chunk.voxels[p.Y][p.X][p.Z]
}


func (chunk *Chunk) Set (p Point, v Voxel) {
    chunk.voxels[p.Y][p.X][p.Z] = v
}
