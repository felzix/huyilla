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