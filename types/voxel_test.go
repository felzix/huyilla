package types

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestVoxel(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Voxel Suite")
}

var _ = Describe("Basic", func() {
	It("compress-expand preserves typical values", func() {
		voxel := ExpandedVoxel{
			Form:        12,
			Material:    9004,
			Other:       0,
			Temperature: RoomTemperature,
			Pressure:    0,
			Multiblock:  0,
		}.Compress().Expand()

		Expect(voxel.Form).To(Equal(uint64(12)))
		Expect(voxel.Material).To(Equal(uint64(9004)))
		Expect(voxel.Other).To(Equal(uint64(0)))
		Expect(voxel.Temperature).To(Equal(RoomTemperature))
		Expect(voxel.Pressure).To(Equal(uint64(0)))
		Expect(voxel.Multiblock).To(Equal(uint64(0)))
	})

	It("compress-expand preserves form", func() {
		voxel := ExpandedVoxel{
			Form: 0xFFFF,
		}.Compress().Expand()

		Expect(voxel.Form).To(Equal(uint64(0xFFFF)))
	})

	It("compress-expand preserves material", func() {
		voxel := ExpandedVoxel{
			Material: 0xFFFF,
		}.Compress().Expand()

		Expect(voxel.Material).To(Equal(uint64(0xFFFF)))
	})

	It("compress-expand preserves other", func() {
		voxel := ExpandedVoxel{
			Other: 0x3fff,
		}.Compress().Expand()

		Expect(voxel.Other).To(Equal(uint64(0x3fff)))
	})

	It("compress-expand preserves temperature", func() {
		voxel := ExpandedVoxel{
			Temperature: MaxTemperature,
		}.Compress().Expand()

		Expect(voxel.Temperature).To(Equal(uint64(MaxTemperature)))
	})

	It("compress-expand preserves pressure", func() {
		voxel := ExpandedVoxel{
			Pressure: MaxPressure,
		}.Compress().Expand()

		Expect(voxel.Pressure).To(Equal(MaxPressure))
	})

	It("compress-expand preserves multiblock", func() {
		voxel := ExpandedVoxel{
			Multiblock: 1,
		}.Compress().Expand()

		Expect(voxel.Multiblock).To(Equal(uint64(1)))
	})
})
