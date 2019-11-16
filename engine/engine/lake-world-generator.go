package engine

import (
	C "github.com/felzix/huyilla/constants"
	"github.com/felzix/huyilla/types"
)

type LakeWorldGenerator struct {
	worldSeed  uint64
	lakeRadius float64
	lakeCenter types.Point
}

func NewLakeWorldGenerator(lakeRadius uint8) *LakeWorldGenerator {
	return &LakeWorldGenerator{
		lakeRadius: float64(lakeRadius),
	}
}

func (gen *LakeWorldGenerator) SetupForWorld() {}

func (gen *LakeWorldGenerator) SetupForChunk(_ types.Point) {
	center := types.RandomPoint(C.CHUNK_SIZE)
	center.Z = 0
	gen.lakeCenter = center
}

func (gen *LakeWorldGenerator) GenVoxel(p types.AbsolutePoint) types.Voxel {
	m := MATERIAL
	f := FORM

	if p.Chunk.Z < 0 {
		return types.ExpandedVoxel{
			Form:        f["cube"],
			Material:    m["dirt"],
			Temperature: types.RoomTemperature,
		}.Compress()
	}

	if p.Chunk.Z > 0 {
		return types.ExpandedVoxel{
			Form:        f["cube"],
			Material:    m["air"],
			Temperature: types.RoomTemperature,
		}.Compress()
	}

	if p.Voxel.Z == gen.lakeCenter.Z {
		if p.Voxel.Distance(gen.lakeCenter) <= gen.lakeRadius {
			return types.ExpandedVoxel{
				Form:        f["cube"],
				Material:    m["water"],
				Temperature: types.RoomTemperature,
			}.Compress()
		} else {
			return types.ExpandedVoxel{
				Form:        f["cube"],
				Material:    m["dirt"],
				Temperature: types.RoomTemperature,
			}.Compress()
		}
	} else if p.Voxel.Z > gen.lakeCenter.Z {
		return types.ExpandedVoxel{
			Form:        f["cube"],
			Material:    m["air"],
			Temperature: types.RoomTemperature,
		}.Compress()
	} else {
		return types.ExpandedVoxel{
			Form:        f["cube"],
			Material:    m["dirt"],
			Temperature: types.RoomTemperature,
		}.Compress()
	}
}
