package main

import (
	"github.com/felzix/huyilla/types"
	"testing"
)

func TestHuyilla_Actions(t *testing.T) {
	h := &Engine{}
	h.Init()
	defer h.World.WipeDatabase()

	NAME := "FAKE"
	PASS := "PASS"

	if err := h.SignUp(NAME, PASS); err != nil {
		t.Fatal(err)
	}
	if _, err := h.LogIn(NAME, PASS); err != nil {
		t.Fatal(err)
	}

	action := types.Action{
		PlayerName: NAME,
		Action: &types.Action_Move{
			Move: &types.Action_MoveAction{
				WhereTo: &types.AbsolutePoint{
					Chunk: &types.Point{X: 1, Y: 12, Z: 144},
					Voxel: &types.Point{X: 2, Y: 4, Z: 8},
				},
			},
		},
	}

	// tests behavior when there are no queued actions
	if len(h.Actions) != 0 {
		t.Errorf(`Expected 0 action but found %d`, len(h.Actions))
	}

	h.RegisterAction(&action)

	if len(h.Actions) != 1 {
		t.Errorf(`Expected 1 action but found %d`, len(h.Actions))
	}

	// tests behavior when there are queued actions
	h.RegisterAction(&action)

	if len(h.Actions) != 2 {
		t.Errorf(`Expected 2 actions but found %d`, len(h.Actions))
	}

	// tests behavior when action queue is reset
	if err := h.Tick(); err != nil {
		t.Fatal(err)
	}

	if len(h.Actions) != 0 {
		t.Errorf(`Expected 0 actions but found %d`, len(h.Actions))
	}
}

func TestHuyilla_Move(t *testing.T) {
	h := &Engine{}
	h.Init()
	defer h.World.WipeDatabase()

	NAME := "felzix"
	PASS := "PASS"
	CHUNK_POINT := newPoint(0, 0, 0)
	VOXEL_POINT := newPoint(2, 4, 8)

	if err := h.SignUp(NAME, PASS); err != nil {
		t.Fatal(err)
	}
	if _, err := h.LogIn(NAME, PASS); err != nil {
		t.Fatal(err)
	}

	h.RegisterAction(&types.Action{
		PlayerName: NAME,
		Action: &types.Action_Move{
			Move: &types.Action_MoveAction{
				WhereTo: &types.AbsolutePoint{Chunk: CHUNK_POINT, Voxel: VOXEL_POINT},
			}}})

	if err := h.Tick(); err != nil {
		t.Fatal(err)
	}

	player, err := h.World.Player(NAME)
	if err != nil {
		t.Error("Error:", err)
	}

	entity, err := h.World.Entity(player.EntityId)
	if entity == nil {
		t.Fatal("Entity should exist but doesn't")
	} else if err != nil {
		t.Fatal(err)
	}

	if !(pointEquals(entity.Location.Chunk, CHUNK_POINT) && pointEquals(entity.Location.Voxel, VOXEL_POINT)) {
		t.Errorf(`Player should be at "%s" but is at "%s"`,
			absolutePointToString(&types.AbsolutePoint{Chunk: CHUNK_POINT, Voxel: VOXEL_POINT}),
			absolutePointToString(entity.Location))
	}
}
