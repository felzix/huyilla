package main

import (
	"fmt"
	C "github.com/felzix/huyilla/constants"
	"github.com/felzix/huyilla/types"
	"math"
	"math/rand"
)

func newAbsolutePoint(cX, cY, cZ, vX, vY, vZ int64) *types.AbsolutePoint {
	return &types.AbsolutePoint{
		Chunk: newPoint(cX, cY, cZ),
		Voxel: newPoint(vX, vY, vZ),
	}
}

func newPoint(x, y, z int64) *types.Point {
	return &types.Point{x, y, z}
}

func clonePoint(p *types.Point) *types.Point {
	return newPoint(p.X, p.Y, p.Z)
}

func pointEquals(p0, p1 *types.Point) bool {
	return p0.X == p1.X &&
		p0.Y == p1.Y &&
		p0.Z == p1.Z
}

func randomPoint() *types.Point {
	x := rand.Int63n(C.CHUNK_SIZE)
	y := rand.Int63n(C.CHUNK_SIZE)
	z := rand.Int63n(C.CHUNK_SIZE)
	return newPoint(x, y, z)
}

func distance(p0, p1 *types.Point) float64 {
	x := float64(p0.X - p1.X)
	y := float64(p0.Y - p1.Y)
	z := float64(p0.Z - p1.Z)

	return math.Sqrt(x*x + y*y + z*z)
}

func pointToString(p *types.Point) string {
	return fmt.Sprintf("(%d,%d,%d)", p.X, p.Y, p.Z)
}

func absolutePointToString(p *types.AbsolutePoint) string {
	return fmt.Sprintf("(%s,%s)", pointToString(p.Chunk), pointToString(p.Voxel))
}
