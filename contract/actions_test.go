package main

import (
    "github.com/felzix/huyilla/types"
    "github.com/loomnetwork/go-loom"
    "github.com/loomnetwork/go-loom/plugin"
    "github.com/loomnetwork/go-loom/plugin/contractpb"
    "testing"
)



func TestHuyilla_Actions (t *testing.T) {
    h := &Huyilla{}

    addr1 := loom.MustParseAddress(ADDR_FROM_LOOM_EXAMPLE)
    ctx := contractpb.WrapPluginContext(plugin.CreateFakeContext(addr1, addr1))

    h.Init(ctx, &plugin.Request{})

    action := types.Action{
        PlayerName: "admin",
        Action: &types.Action_Move{
            Move: &types.Action_MoveAction{
                WhereTo: &types.AbsolutePoint{
                    &types.Point{1, 12, 144},
                    &types.Point{2, 4, 8},
                },
            },
        },
    }

    // tests behavior when there are no queued actions
    actions, err := h.getActions(ctx)
    if err != nil {
        t.Fatalf("Error: %v", err)
    }
    if len(actions.Actions) != 0 {
        t.Errorf(`Expected 0 action but found %d`, len(actions.Actions))
    }

    err = h.RegisterAction(ctx, &action)
    if err != nil {
        t.Fatalf("Error: %v", err)
    }

    actions, err = h.getActions(ctx)
    if err != nil {
        t.Fatalf("Error: %v", err)
    }
    if len(actions.Actions) != 1 {
        t.Errorf(`Expected 1 action but found %d`, len(actions.Actions))
    }

    // tests behavior when there are queued actions
    err = h.RegisterAction(ctx, &action)
    if err != nil {
        t.Fatalf("Error: %v", err)
    }

    actions, err = h.getActions(ctx)
    if err != nil {
        t.Fatalf("Error: %v", err)
    }
    if len(actions.Actions) != 2 {
        t.Errorf(`Expected 2 actions but found %d`, len(actions.Actions))
    }

    // tests behavior when action queue is reset
    err = h.Tick(ctx, &plugin.Request{})
    if err != nil {
        t.Fatalf("Error: %v", err)
    }

    actions, err = h.getActions(ctx)
    if err != nil {
        t.Fatalf("Error: %v", err)
    }
    if len(actions.Actions) != 0 {
        t.Errorf(`Expected 0 actions but found %d`, len(actions.Actions))
    }
}


func TestHuyilla_Move (t *testing.T) {
    h := &Huyilla{}

    addr1 := loom.MustParseAddress(ADDR_FROM_LOOM_EXAMPLE)
    ctx := contractpb.WrapPluginContext(plugin.CreateFakeContext(addr1, addr1))

    h.Init(ctx, &plugin.Request{})

    NAME := "felzix"
    CHUNK_POINT := newPoint(0, 0, 0)
    VOXEL_POINT := newPoint(2, 4, 8)

    h.SignUp(ctx, &types.PlayerName{Name: NAME})

    err := h.RegisterAction(ctx, &types.Action{
        PlayerName: NAME,
        Action: &types.Action_Move{
            Move: &types.Action_MoveAction{
                WhereTo: &types.AbsolutePoint{CHUNK_POINT, VOXEL_POINT},
            },
        },
    })
    if err != nil {
        t.Fatalf("Error: %v", err)
    }

    err = h.Tick(ctx, &plugin.Request{})
    if err != nil {
        t.Fatalf("Error: %v", err)
    }

    player, err := h.GetPlayer(ctx, &types.PlayerName{Name: NAME})
    if err != nil {
        t.Error("Error:", err)
    }
    entity := player.Entity

    if !(pointEquals(entity.Location.Chunk, CHUNK_POINT) && pointEquals(entity.Location.Voxel, VOXEL_POINT)) {
        t.Errorf(`Player should be at "%s" but is at "%s"`,
            absolutePointToString(&types.AbsolutePoint{Chunk: CHUNK_POINT, Voxel: VOXEL_POINT}),
            absolutePointToString(entity.Location))
    }
}
