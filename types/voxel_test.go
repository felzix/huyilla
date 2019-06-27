package types

import (
	. "github.com/felzix/goblin"
	"testing"
)

func TestVoxel(t *testing.T) {
	g := Goblin(t)

	g.Describe("Basic", func() {
		g.It("compress-expand preserves values", func() {
			v := (ExpandedVoxel{
				Form:        12,
				Material:    9004,
				Temperature: RoomTemperature,
				Pressure:    0,
			}).Compress()

			voxel := v.Expand()

			g.Assert(voxel.Form).Equal(uint64(12))
			g.Assert(voxel.Material).Equal(uint64(9004))
			g.Assert(voxel.Temperature).Equal(RoomTemperature)
			g.Assert(voxel.Pressure).Equal(uint64(0))
		})
	})
}
