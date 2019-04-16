package engine

import (
	"fmt"
	C "github.com/felzix/huyilla/constants"
	"github.com/felzix/huyilla/types"
	"math"
	"math/rand"
	"strconv"
)

func NewAbsolutePoint(cX, cY, cZ, vX, vY, vZ int64) *types.AbsolutePoint {
	return &types.AbsolutePoint{
		Chunk: NewPoint(cX, cY, cZ),
		Voxel: NewPoint(vX, vY, vZ),
	}
}

func NewPoint(x, y, z int64) *types.Point {
	return &types.Point{X: x, Y: y, Z: z}
}

func clonePoint(p *types.Point) *types.Point {
	return NewPoint(p.X, p.Y, p.Z)
}

func CloneAbsolutePoint(p *types.AbsolutePoint) *types.AbsolutePoint {
	return NewAbsolutePoint(p.Chunk.X, p.Chunk.Y, p.Chunk.Z, p.Voxel.X, p.Voxel.Y, p.Voxel.Z)
}

func pointEquals(p0, p1 *types.Point) bool {
	return p0.X == p1.X &&
		p0.Y == p1.Y &&
		p0.Z == p1.Z
}

func absolutePointEquals(p0, p1 *types.AbsolutePoint) bool {
	return pointEquals(p0.Chunk, p1.Chunk) &&
		pointEquals(p0.Voxel, p1.Voxel)
}

func randomPoint() *types.Point {
	x := rand.Int63n(C.CHUNK_SIZE)
	y := rand.Int63n(C.CHUNK_SIZE)
	z := rand.Int63n(C.CHUNK_SIZE)
	return NewPoint(x, y, z)
}

func distance(p0, p1 *types.Point) float64 {
	x := float64(p0.X - p1.X)
	y := float64(p0.Y - p1.Y)
	z := float64(p0.Z - p1.Z)

	return math.Sqrt(x*x + y*y + z*z)
}

func absInt64(i int64) int64 {
	if i >= 0 {
		return i
	} else {
		return -i
	}
}

func gridDistance(p0, p1 *types.Point) int64 {
	x := absInt64(p0.X - p1.X)
	y := absInt64(p0.Y - p1.Y)
	z := absInt64(p0.Z - p1.Z)

	most := x
	if y > most {
		most = y
	}
	if z > most {
		most = z
	}

	return most
}

// If you don't need to know the specific distance,
// just if it's more or less than another distance,
// there's no need to take the square root.
func relativeDistance(p0, p1 *types.Point) float64 {
	x := float64(p0.X - p1.X)
	y := float64(p0.Y - p1.Y)
	z := float64(p0.Z - p1.Z)

	return x*x + y*y + z*z
}

func pointToString(p *types.Point) string {
	return fmt.Sprintf("(%d,%d,%d)", p.X, p.Y, p.Z)
}

func stringToPoint(x, y, z string) (*types.Point, error) {
	X, err := strconv.ParseInt(x, 10, 64)
	if err != nil {
		return nil, err
	}

	Y, err := strconv.ParseInt(y, 10, 64)
	if err != nil {
		return nil, err
	}

	Z, err := strconv.ParseInt(z, 10, 64)
	if err != nil {
		return nil, err
	}

	return NewPoint(X, Y, Z), nil
}

func absolutePointToString(p *types.AbsolutePoint) string {
	return fmt.Sprintf("(%s,%s)", pointToString(p.Chunk), pointToString(p.Voxel))
}
