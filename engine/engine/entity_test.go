package engine

import (
	C "github.com/felzix/huyilla/constants"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestEntity(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Entity Suite")
}

var _ = Describe("Entity", func() {
	var h *Engine

	BeforeEach(func() {
		engine, err := NewEngine(C.SEED, NewLakeWorldGenerator(3), NewMemoryDatabase())
		h = engine
		Expect(err).To(BeNil())
	})

	It("loads human type", func() {
		err := h.SignUp("felzix", "PASS")
		Expect(err).To(BeNil())
		_, err = h.LogIn("felzix", "PASS")
		Expect(err).To(BeNil())

		player, err := h.World.Player("felzix")
		Expect(err).To(BeNil())
		Expect(player).ToNot(BeNil())

		entity, err := h.World.Entity(player.EntityId)
		Expect(err).To(BeNil())
		Expect(entity).ToNot(BeNil())
	})
})
