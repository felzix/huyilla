package engine

import (
	C "github.com/felzix/huyilla/constants"
	"github.com/felzix/huyilla/types"
)

type LakeWorldGenerator struct {
	worldSeed uint64
	lakeRadius uint8
	lakeCenter *types.Point
}

func NewLakeWorldGenerator(lakeRadius uint8) *LakeWorldGenerator {
	return &LakeWorldGenerator{
		lakeRadius: lakeRadius,
	}
}

func (gen *LakeWorldGenerator) SetupForWorld() {}

func (gen *LakeWorldGenerator) SetupForChunk(_ *types.Point) {
	center := types.RandomPoint(C.CHUNK_SIZE)
	center.Z = 0
	gen.lakeCenter = center
}

func (gen *LakeWorldGenerator) GenVoxel(p *types.AbsolutePoint) uint64 {
	v := VOXEL

	if p.Chunk.Z < 0 {
		return v["dirt"]
	}

	if p.Chunk.Z > 0 {
		return v["air"]
	}

	d := p.Voxel.Distance(gen.lakeCenter)
	if p.Voxel.Z == gen.lakeCenter.Z && d <= float64(3) {
		return v["water"]
	}
	if p.Voxel.Z == 0 {
		return v["barren_earth"]
	}
	return v["air"]
}
