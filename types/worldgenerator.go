package types


// All methods are idempotent
type WorldGenerator interface {
	SetupForWorld()
	SetupForChunk(chunkLocation Point)
	GenVoxel(voxelLocation AbsolutePoint) Voxel
}

