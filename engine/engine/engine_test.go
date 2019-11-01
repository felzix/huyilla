package engine

import (
	C "github.com/felzix/huyilla/constants"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestEngine(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Engine Suite")
}

var _ = Describe("Engine", func() {
	var h *Engine

	BeforeEach(func() {
		engine, err := NewEngine(C.SEED, NewLakeWorldGenerator(3), NewMemoryDatabase())
		h = engine
		Expect(err).To(BeNil())
	})

	It("ticks", func() {
		err := h.SignUp("felzix", "PASS")
		Expect(err).To(BeNil())

		_, err = h.LogIn("felzix", "PASS")
		Expect(err).To(BeNil())

		err = h.Tick()
		Expect(err).To(BeNil())

		content := h.GetContent()
		Expect(content.E[0].Name).To(Equal("human"))

		player, err := h.World.Player("felzix")
		Expect(err).To(BeNil())

		entity, err := h.World.Entity(player.EntityId)
		Expect(err).To(BeNil())
		Expect(entity).ToNot(BeNil())

		chunk, err := h.World.Chunk(entity.Location.Chunk)
		Expect(err).To(BeNil())
		Expect(len(chunk.Entities)).To(Equal(1))

		// active range in the positive direction
		edge := entity.Location.Chunk.Clone()
		edge.X += 3
		chunk, err = h.World.Chunk(edge)
		Expect(err).To(BeNil())
		Expect(chunk).ToNot(BeNil()) // chunk within player's range should exist
		Expect(len(chunk.Entities)).To(Equal(0))

		beyond := entity.Location.Chunk.Clone()
		beyond.X += 4
		chunk, err = h.World.OnlyGetChunk(beyond)

		Expect(chunk).To(BeNil()) // Chunk beyond player's range exists

		// active range in the negative direction
		edge = entity.Location.Chunk.Clone()
		edge.X -= 3
		chunk, err = h.World.Chunk(edge)
		Expect(err).To(BeNil())
		Expect(chunk).ToNot(BeNil()) // Chunk within player's range should exist

		beyond = entity.Location.Chunk.Clone()
		beyond.X -= 4
		chunk, err = h.World.OnlyGetChunk(beyond)
		Expect(err).To(BeNil())
		Expect(chunk).To(BeNil()) // Chunk beyond player's range exists
	})
})
