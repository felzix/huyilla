package types

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
)

func NewMove(player string, whereTo *AbsolutePoint) *Action{
	return &Action{
		PlayerName: player,
		Action: &Action_Move{
			Move: &Action_MoveAction{
				WhereTo: whereTo,
			},
		},
	}
}


func NewAbsolutePoint(cX, cY, cZ, vX, vY, vZ int64) *AbsolutePoint {
	return &AbsolutePoint{
		Chunk: NewPoint(cX, cY, cZ),
		Voxel: NewPoint(vX, vY, vZ),
	}
}

func NewPoint(x, y, z int64) *Point {
	return &Point{X: x, Y: y, Z: z}
}

func ClonePoint(p *Point) *Point {
	return NewPoint(p.X, p.Y, p.Z)
}

func CloneAbsolutePoint(p *AbsolutePoint) *AbsolutePoint {
	return NewAbsolutePoint(p.Chunk.X, p.Chunk.Y, p.Chunk.Z, p.Voxel.X, p.Voxel.Y, p.Voxel.Z)
}

func PointEquals(p0, p1 *Point) bool {
	return p0.X == p1.X &&
		p0.Y == p1.Y &&
		p0.Z == p1.Z
}

func AbsolutePointEquals(p0, p1 *AbsolutePoint) bool {
	return PointEquals(p0.Chunk, p1.Chunk) &&
		PointEquals(p0.Voxel, p1.Voxel)
}

func RandomPoint(chunkSize int64) *Point {
	x := rand.Int63n(chunkSize)
	y := rand.Int63n(chunkSize)
	z := rand.Int63n(chunkSize)
	return NewPoint(x, y, z)
}

func Distance(p0, p1 *Point) float64 {
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

func GridDistance(p0, p1 *Point) int64 {
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
func RelativeDistance(p0, p1 *Point) float64 {
	x := float64(p0.X - p1.X)
	y := float64(p0.Y - p1.Y)
	z := float64(p0.Z - p1.Z)

	return x*x + y*y + z*z
}

func PointToString(p *Point) string {
	return fmt.Sprintf("(%d,%d,%d)", p.X, p.Y, p.Z)
}

func StringToPoint(x, y, z string) (*Point, error) {
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

func AbsolutePointToString(p *AbsolutePoint) string {
	return fmt.Sprintf("(%s,%s)", PointToString(p.Chunk), PointToString(p.Voxel))
}
