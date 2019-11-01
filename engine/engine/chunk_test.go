package engine

import (
	C "github.com/felzix/huyilla/constants"
	"github.com/felzix/huyilla/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestChunk(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Client Suite")
}

var _ = Describe("Chunk", func() {
	var h *Engine

	BeforeEach(func() {
		engine, err := NewEngine(C.SEED, NewLakeWorldGenerator(3), NewMemoryDatabase())
		h = engine
		Expect(err).To(BeNil())
	})

	It("implicitly generates a chunk", func() {
		chunk, err := h.World.Chunk(&types.Point{X: 0, Y: 0, Z: 0})
		Expect(err).To(BeNil())

		expectedVoxelCount := C.CHUNK_SIZE * C.CHUNK_SIZE * C.CHUNK_SIZE
		Expect(len(chunk.Voxels)).To(Equal(expectedVoxelCount))
	})

	It("generates 1,000 chunks", func() {
		expectedVoxelCount := C.CHUNK_SIZE * C.CHUNK_SIZE * C.CHUNK_SIZE

		for i := 0; i < 1000; i++ {
			chunk, err := h.World.GenerateChunk(&types.Point{X: 0, Y: 0, Z: 0})
			Expect(err).To(BeNil())
			Expect(len(chunk.Voxels)).To(Equal(expectedVoxelCount))
		}
	})

	It("adds entity to chunk", func() {
		p := types.NewAbsolutePoint(0, 0, 0, 0, 0, 0)

		_, err := h.World.GenerateChunk(p.Chunk)
		Expect(err).To(BeNil())

		entity, err := h.World.CreateEntity(0, "", p)
		Expect(err).To(BeNil())

		err = h.World.AddEntityToChunk(entity)
		Expect(err).To(BeNil())

		chunk, err := h.World.Chunk(p.Chunk)
		Expect(err).To(BeNil())
		Expect(chunk).ToNot(BeNil())

		Expect(len(chunk.Entities)).To(Equal(1))
	})

	It("removes entity from chunk", func() {
		p := types.NewAbsolutePoint(0, 0, 0, 0, 0, 0)

		_, err := h.World.GenerateChunk(p.Chunk)
		Expect(err).To(BeNil())

		entity, err := h.World.CreateEntity(0, "", p)
		Expect(err).To(BeNil())

		err = h.World.AddEntityToChunk(entity)
		Expect(err).To(BeNil())

		err = h.World.RemoveEntityFromChunk(entity.Id, p.Chunk)
		Expect(err).To(BeNil())

		chunk, err := h.World.Chunk(p.Chunk)
		Expect(err).To(BeNil())

		Expect(len(chunk.Entities)).To(Equal(0))
	})
})
