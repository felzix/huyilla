package main

import (
    "github.com/felzix/huyilla/types"
    contract "github.com/loomnetwork/go-loom/plugin/contractpb"
)


var ACTIONS = []byte("Actions")


func (c *Huyilla) RegisterAction (ctx contract.Context, req *types.Action) error {
    actions, err := c.getActions(ctx)
    if err != nil { return err }

    actions.Actions = append(actions.Actions, req)

    return ctx.Set(ACTIONS, actions)
}

func (c *Huyilla) getActions (ctx contract.StaticContext) (*types.Actions, error) {
    var actions types.Actions
    var err error = nil

    if ctx.Has(ACTIONS) {
        err = ctx.Get(ACTIONS, &actions)
    }

    return &actions, err
}

// returns true if move succeeded; false otherwise
func (c *Huyilla) move (ctx contract.Context, action *types.Action) (bool, error) {
    player, err := c.getPlayer(ctx, action.PlayerName)
    if err != nil { return false, err }

    if player.Entity == nil {
        return false, nil  // player doesn't have an entity (player has not yet finished signup)
    }

    err = c.removeEntityFromChunk(ctx, player.Entity)
    if err != nil { return false, err }

    player.Entity.Location = action.GetMove().WhereTo

    err = c.setEntity(ctx, player.Entity)
    if err != nil { return false, err }

    err = c.addEntityToChunk(ctx, player.Entity)
    if err != nil { return false, err }

    return true, nil
}