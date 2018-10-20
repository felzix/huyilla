package main

import (
    "github.com/felzix/huyilla/types"
    "github.com/loomnetwork/go-loom/plugin"
    contract "github.com/loomnetwork/go-loom/plugin/contractpb"
    "github.com/pkg/errors"
)


var PLAYERS = []byte("Players")


func (c *Huyilla) GetPlayerList (ctx contract.StaticContext, req *plugin.Request) (*types.PlayerList, error) {
    players, err := c.getPlayers(ctx)
    if err != nil { return nil, err }

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

    // NOTE: entity can be nil
    details := &types.PlayerDetails{Player: player, Entity: entity}
    return details, nil
}

func (c *Huyilla) getPlayer (ctx contract.StaticContext, name string) (*types.Player, error) {
    players, err := c.getPlayers(ctx)
    if err != nil { return nil, err }
    player := players.Players[name]
    return player, nil
}

func (c *Huyilla) getPlayers (ctx contract.StaticContext) (*types.Players, error) {
    var players types.Players
    if err := ctx.Get(PLAYERS, &players); err != nil {return nil, err}
    return &players, nil
}

func (c *Huyilla) getActivePlayers (ctx contract.StaticContext) ([]*types.PlayerDetails, error) {
    players, err := c.getPlayers(ctx)
    if err != nil { return nil, err }

    var activePlayers []*types.PlayerDetails

    for _, player := range players.Players {
        if player.LoggedIn {
            entity, err := c.getEntity(ctx, player.Id)
            if err != nil { return nil, err }
            activePlayers = append(activePlayers, &types.PlayerDetails{Player: player, Entity: entity})
        }
    }

    return activePlayers, nil
}

//
// func (c *Huyilla) getPlayersDetails (ctx contract.StaticContext) ([]*types.PlayerDetails, error) {
//     players, err := c.getPlayers(ctx)
//     if err != nil { return nil, err }
//
//     details := make([]*types.PlayerDetails, len(players.Players))
//
//     for i, player := range players.Players {
//         entity, err := c.getEntity(ctx, player.Id)
//         if err != nil { return nil, err }
//         details[i] = &types.PlayerDetails{Player: player, Entity: entity}
//     }
//
//     return &players, nil
// }
