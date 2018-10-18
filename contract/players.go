package main

import (
    "github.com/felzix/huyilla/types"
    "github.com/loomnetwork/go-loom/plugin"
    contract "github.com/loomnetwork/go-loom/plugin/contractpb"
    "github.com/pkg/errors"
)


var PLAYERS = []byte("Players")


func (c *Huyilla) GetPlayerList (ctx contract.StaticContext, req *plugin.Request) (*types.PlayerList, error) {
    var players = &types.Players{}
    if err := ctx.Get(PLAYERS, players); err != nil { return nil, err }

    list := types.PlayerList{}
    for name, _ := range players.Players {
        list.Names = append(list.Names, name)
    }
    return &list, nil
}

func (c *Huyilla) GetPlayer (ctx contract.StaticContext, req *types.PlayerName) (*types.PlayerDetails, error) {
    player, err := c.getPlayer(ctx, req.Name)
    if err != nil {return nil, errors.Wrap(err, "player not found")}

    entity, err := c.getEntity(ctx, player.Id)
    if err != nil {return nil, errors.Wrap(err, "entity not found")}

    details := &types.PlayerDetails{Player: player, Entity: entity}
    return details, nil
}

func (c *Huyilla) getPlayer (ctx contract.StaticContext, name string) (*types.Player, error) {
    var players types.Players
    if err := ctx.Get(PLAYERS, &players); err != nil {return nil, err}
    player := players.Players[name]
    return player, nil
}

func (c *Huyilla) getPlayers (ctx contract.StaticContext) (*types.Players, error) {
    var players types.Players
    if err := ctx.Get(PLAYERS, &players); err != nil {return nil, err}
    return &players, nil
}

func (c *Huyilla) thisUser (ctx contract.StaticContext) []byte{
    return []byte(ctx.Message().Sender.Local)
}