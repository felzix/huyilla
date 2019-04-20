package types

import (
	"fmt"
	"math/rand"
	"strconv"
)

func NewChunk(tick, chunkLength uint64) *Chunk {
	return &Chunk{
		Tick: tick,
		Voxels: make([]uint64, chunkLength),
	}
}

func NewMove(player string, whereTo *AbsolutePoint) *Action {
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

func (m *Point) Clone() *Point {
	return NewPoint(m.X, m.Y, m.Z)
}


func RandomPoint(chunkSize int64) *Point {
	x := rand.Int63n(chunkSize)
	y := rand.Int63n(chunkSize)
	z := rand.Int63n(chunkSize)
	return NewPoint(x, y, z)
}

func (m *Point) ToString() string {
	return fmt.Sprintf("(%d,%d,%d)", m.X, m.Y, m.Z)
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

func (m *AbsolutePoint) ToString() string {
	return fmt.Sprintf("(%s,%s)", m.Chunk.ToString(), m.Voxel.ToString())
}
