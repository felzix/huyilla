package types

import (
	. "github.com/felzix/goblin"
	"testing"
)

func TestVoxel(t *testing.T) {
	g := Goblin(t)

	g.Describe("Basic", func() {
		g.It("compress-expand preserves typical values", func() {
			voxel := ExpandedVoxel{
				Form:        12,
				Material:    9004,
				Other:       0,
				Temperature: RoomTemperature,
				Pressure:    0,
				Multiblock:  0,
			}.Compress().Expand()

			g.Assert(voxel.Form).Equal(uint64(12))
			g.Assert(voxel.Material).Equal(uint64(9004))
			g.Assert(voxel.Other).Equal(uint64(0))
			g.Assert(voxel.Temperature).Equal(RoomTemperature)
			g.Assert(voxel.Pressure).Equal(uint64(0))
			g.Assert(voxel.Multiblock).Equal(uint64(0))
		})

		g.It("compress-expand preserves form", func() {
			voxel := ExpandedVoxel{
				Form: 0xFFFF,
			}.Compress().Expand()

			g.Assert(voxel.Form).Equal(uint64(0xFFFF))
		})

		g.It("compress-expand preserves material", func() {
			voxel := ExpandedVoxel{
				Material: 0xFFFF,
			}.Compress().Expand()

			g.Assert(voxel.Material).Equal(uint64(0xFFFF))
		})

		g.It("compress-expand preserves other", func() {
			voxel := ExpandedVoxel{
				Other: 0x3fff,
			}.Compress().Expand()

			g.Assert(voxel.Other).Equal(uint64(0x3fff))
		})

		g.It("compress-expand preserves temperature", func() {
			voxel := ExpandedVoxel{
				Temperature: MaxTemperature,
			}.Compress().Expand()

			g.Assert(voxel.Temperature).Equal(uint64(MaxTemperature))
		})

		g.It("compress-expand preserves pressure", func() {
			voxel := ExpandedVoxel{
				Pressure: MaxPressure,
			}.Compress().Expand()

			g.Assert(voxel.Pressure).Equal(MaxPressure)
		})

		g.It("compress-expand preserves multiblock", func() {
			voxel := ExpandedVoxel{
				Multiblock: 1,
			}.Compress().Expand()

			g.Assert(voxel.Multiblock).Equal(uint64(1))
		})
	})
}
