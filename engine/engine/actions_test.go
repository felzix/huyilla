package engine

import (
	"fmt"
	C "github.com/felzix/huyilla/constants"
	"github.com/felzix/huyilla/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestAction(t *testing.T) {

	RegisterFailHandler(Fail)
	RunSpecs(t, "Action Suite")
}

var _ = Describe("Action", func() {
		NAME := "FAKE"
		PASS := "PASS"
		var h *Engine

		BeforeEach(func() {
			engine, err := NewEngine(C.SEED, NewLakeWorldGenerator(3), NewMemoryDatabase())
			h = engine
			Expect(err).To(BeNil())

			err = h.SignUp(NAME, PASS)
			Expect(err).To(BeNil())

			Eventually(func() error {
				_, err := h.LogIn(NAME, PASS)
				return err
			}).Should(BeNil())
		})

		Describe("queue actions and tick", func() {
			It("moves within chunk", func() {
				whereTo := types.NewAbsolutePoint(0, 0, 0, 2, 4, 8)

				Expect(len(h.Actions)).To(Equal(0))

				h.RegisterAction(&types.Action{
					PlayerName: NAME,
					Action: &types.Action_Move{
						Move: &types.Action_MoveAction{
							WhereTo: whereTo,
						}}})

				Expect(len(h.Actions)).To(Equal(1))

				err := h.Tick()
				Expect(err).To(BeNil())

				Expect(len(h.Actions)).To(Equal(0))

				player, err := h.World.Player(NAME)
				Expect(err).To(BeNil())
				Expect(player).ToNot(BeNil())
				Expect(player.EntityId).ToNot(Equal(0))

				entity, err := h.World.Entity(player.EntityId)
				Expect(err).To(BeNil())
				Expect(entity).ToNot(BeNil())
				Expect(entity.Location.Equals(whereTo)).To(BeTrue(), fmt.Sprintf(
					`Player should be at "%s" but is at "%s"`,
					whereTo.ToString(),
					entity.Location.ToString(),
				))

				chunk, err := h.World.Chunk(whereTo.Chunk)
				Expect(err).To(BeNil())

				entityPresent := false
				for _, e := range chunk.Entities {
					if e == entity.Id {
						entityPresent = true
					}
				}
				Expect(entityPresent).To(BeTrue(), "Entity is not actually present in the chunk")
			})

			It("moves to another chunk", func() {
				whereTo := types.NewAbsolutePoint(-1, 0, 0, 15, 0, 0)

				move := &types.Action{
					PlayerName: NAME,
					Action: &types.Action_Move{
						Move: &types.Action_MoveAction{
							WhereTo: whereTo,
						},
					},
				}

				success, err := h.Move(move)
				Expect(err).To(BeNil())
				Expect(success).To(Equal(true))

				player, err := h.World.Player(NAME)
				Expect(err).To(BeNil())
				Expect(player).ToNot(BeNil())
				Expect(player.EntityId).ToNot(Equal(0))

				entity, err := h.World.Entity(player.EntityId)
				Expect(err).To(BeNil())
				Expect(entity).ToNot(BeNil())
				Expect(entity.Location.Equals(whereTo)).To(BeTrue(), fmt.Sprintf(
					`Player should be at "%s" but is at "%s"`,
					whereTo.ToString(),
					entity.Location.ToString(),
				))

				chunk, err := h.World.Chunk(whereTo.Chunk)
				Expect(err).To(BeNil())

				entityPresent := false
				for _, e := range chunk.Entities {
					if e == entity.Id {
						entityPresent = true
					}
				}
				Expect(entityPresent).To(BeTrue(), "Entity is not actually present in the chunk")
			})
		})

		It("moves to another chunk, via action registration", func() {
			whereTo := types.NewAbsolutePoint(-1, 0, 0, 0, 0, 0)

			Expect(len(h.Actions)).To(Equal(0))

			h.RegisterAction(&types.Action{
				PlayerName: NAME,
				Action: &types.Action_Move{
					Move: &types.Action_MoveAction{
						WhereTo: whereTo,
					}}})

			Expect(len(h.Actions)).To(Equal(1))

			err := h.Tick()
			Expect(err).To(BeNil())

			Expect(len(h.Actions)).To(Equal(0))

			player, err := h.World.Player(NAME)
			Expect(err).To(BeNil())
			Expect(player).ToNot(BeNil())
			Expect(player.EntityId).ToNot(Equal(0))

			entity, err := h.World.Entity(player.EntityId)
			Expect(err).To(BeNil())
			Expect(entity).ToNot(BeNil())
			Expect(entity.Location.Equals(whereTo)).To(BeTrue(), fmt.Sprintf(
				`Player should be at "%s" but is at "%s"`,
				whereTo.ToString(),
				entity.Location.ToString(),
			))

			chunk, err := h.World.Chunk(whereTo.Chunk)
			Expect(err).To(BeNil())

			entityPresent := false
			for _, e := range chunk.Entities {
				if e == entity.Id {
					entityPresent = true
				}
			}
			Expect(entityPresent).To(BeTrue(), "Entity is not actually present in the chunk")
		})

	})
