package engine

import (
	C "github.com/felzix/huyilla/constants"
	"github.com/felzix/huyilla/types"
	"github.com/ojrac/opensimplex-go"
)

type GrassWorldGenerator struct {
	worldSeed     uint64
	simplex_noise opensimplex.Noise
}

func NewGrassWorldGenerator() *GrassWorldGenerator {
	return &GrassWorldGenerator{
		worldSeed:     C.SEED,
		simplex_noise: opensimplex.NewNormalized(int64(C.SEED)),
	}
}

func (gen *GrassWorldGenerator) SetupForWorld() {
}

func (gen *GrassWorldGenerator) SetupForChunk(_ *types.Point) {
	center := types.RandomPoint(C.CHUNK_SIZE)
	center.Z = 0
}

func (gen *GrassWorldGenerator) GenVoxel(p *types.AbsolutePoint) types.Voxel {
	m := MATERIAL
	f := FORM

	scale := 0.05
	pointX := p.X()
	pointY := p.Y()
	pointZ := p.Z()

	h1 := gen.simplex_noise.Eval2(1+float64(pointX)*scale, 1+float64(pointY)*scale)
	h2 := gen.simplex_noise.Eval2(1+float64(pointX*2)*scale, 1+float64(pointY*2)*scale)
	h3 := gen.simplex_noise.Eval2(1+float64(pointX*4)*scale, 1+float64(pointY*4)*scale)

	maxHeight := 8.0
	height := 0.0
	height -= h1 * 0.5
	height -= h2 * 0.25
	height -= h3 * 0.125
	height *= maxHeight
	h := int64(height)

	if pointZ < h {
		return types.ExpandedVoxel{
			Form:        f["cube"],
			Material:    m["dirt"],
			Temperature: types.RoomTemperature,
		}.Compress()
	} else if pointZ == h {
		return types.ExpandedVoxel{
			Form:        f["cube"],
			Material:    m["grass"],
			Temperature: types.RoomTemperature,
		}.Compress()
	} else {
		return types.ExpandedVoxel{
			Form:        f["cube"],
			Material:    m["air"],
			Temperature: types.RoomTemperature,
		}.Compress()
	}
}
