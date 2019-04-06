// These are necessary because protobuf structs are intentionally non-comparable.

package types

type ComparablePoint struct {
	X int64
	Y int64
	Z int64
}

func NewComparablePoint(point *Point) *ComparablePoint {
	return &ComparablePoint{
		X: point.X,
		Y: point.Y,
		Z: point.Z,
	}
}