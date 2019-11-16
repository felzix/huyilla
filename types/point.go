package types

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
)

type Point struct {
	X int64
	Y int64
	Z int64
}

func NewPoint(x, y, z int64) Point {
	return Point{
		X: x,
		Y: y,
		Z: z,
	}
}

func (p *Point) ToString() string {
	return fmt.Sprintf("(%d,%d,%d)", p.X, p.Y, p.Z)
}

func (p Point) Equals(other Point) bool {
	return p.X == other.X &&
		p.Y == other.Y &&
		p.Z == other.Z
}

func StringToPoint(x, y, z string) (Point, error) {
	X, err := strconv.ParseInt(x, 10, 64)
	if err != nil {
		return Point{}, err
	}

	Y, err := strconv.ParseInt(y, 10, 64)
	if err != nil {
		return Point{}, err
	}

	Z, err := strconv.ParseInt(z, 10, 64)
	if err != nil {
		return Point{}, err
	}

	return NewPoint(X, Y, Z), nil
}

func (p Point) Clone() Point {
	return NewPoint(p.X, p.Y, p.Z)
}

func (p Point) Marshal() ([]byte, error) {
	return ToBytes(&p)
}

func (p *Point) Unmarshal(input []byte) error {
	return FromBytes(input, &p)
}


func (p Point) DeriveVector(other Point) Point {
	return Point{
		X: other.X - p.X,
		Y: other.Y - p.Y,
		Z: other.Z - p.Z,
	}
}

func (p Point) Distance(other Point) float64 {
	x := float64(p.X - other.X)
	y := float64(p.Y - other.Y)
	z := float64(p.Z - other.Z)

	return math.Sqrt(x*x + y*y + z*z)
}

func (p Point) GridDistance(other Point) int64 {
	x := absInt64(p.X - other.X)
	y := absInt64(p.Y - other.Y)
	z := absInt64(p.Z - other.Z)

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
func (p Point) RelativeDistance(other Point) float64 {
	x := float64(p.X - other.X)
	y := float64(p.Y - other.Y)
	z := float64(p.Z - other.Z)

	return x*x + y*y + z*z
}

func RandomPoint(chunkSize uint64) Point {
	size := int64(chunkSize)
	x := rand.Int63n(size)
	y := rand.Int63n(size)
	z := rand.Int63n(size)
	return NewPoint(x, y, z)
}
