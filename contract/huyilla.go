package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "github.com/pkg/errors"
    "github.com/felzix/huyilla/types"
    "github.com/loomnetwork/go-loom/plugin"
    contract "github.com/loomnetwork/go-loom/plugin/contractpb"
)


func main() {
    plugin.Serve(Contract)
}

type Huyilla struct {}


var Contract = contract.MakePluginContract(&Huyilla{})


func (c *Huyilla) Meta () (plugin.Meta, error) {
    return plugin.Meta{
        Name:    "Huyilla",
        Version: "0.0.1",
    }, nil
}

func (c *Huyilla) Init (ctx contract.Context, req *plugin.Request) error {
    err := ctx.Set(AGE, &types.Age{Ticks: 1})  // starts at 1 because 0 counts as non-existent
    if err != nil { return err }

    config := &types.Config{
        Options: &types.PrimitiveMap{
            Map: map[string]*types.Primitive{
                "PlayerCap": {Value: &types.Primitive_Int{Int: 10}},
            },
        },
    }
    err = ctx.Set(CONFIG, config)
    if err != nil { return err }

    adminEntity := c.newEntity(ctx, 1, "admin")
    err = c.setEntity(ctx, adminEntity)
    if err != nil {return err}

    err = ctx.Set(PLAYERS, &types.Players{
        Players: map[string]*types.Player{
            "admin": {Id: adminEntity.Id,
                      Name: "admin",
                      LoggedIn: false}},
        })
    if err != nil { return err }

    return c.genChunk(ctx, &types.Point{0, 0, 0})
}


func (c *Huyilla) SignUp(ctx contract.Context, req *types.PlayerName) error {
    players, err := c.getPlayers(ctx)
    if err != nil {return err}

    player := players.Players[req.Name]
    if player != nil {
        return errors.New(fmt.Sprintf(`Name "%v" is taken. Try another.`, req.Name))
    }

    entity := c.newEntity(ctx, uint32(1), req.Name)
    err = c.setEntity(ctx, entity)
    if err != nil {return err}

    player = &types.Player{
        Id:      entity.Id,
        Name:    req.Name,
        Address: c.thisUser(ctx),
    }
    players.Players[player.Name] = player
    err = ctx.Set(PLAYERS, players)
    if err != nil {return err}

    ctx.Logger().Info("Created player", "name", player.Name, "address", player.Address)

    emitMsg := struct {
        Method string
        Owner  string
        Addr   []byte
    }{"CreatePlayer", player.Name, player.Address}
    emitMsgJSON, err := json.Marshal(emitMsg)
    if err != nil {return err}

    ctx.EmitTopics(emitMsgJSON, "huyilla:" + emitMsg.Method)
    return nil
}

func (c *Huyilla) LogIn (ctx contract.Context, req *types.PlayerName) (*types.PlayerDetails, error) {
    players, err := c.getPlayers(ctx)
    player := players.Players[req.Name]
    if err != nil { return nil, err }

    if !bytes.Equal(player.Address, c.thisUser(ctx)) {
        return nil, errors.New("Username is not associated with your address/key/account.")
    }

    if !player.LoggedIn {
        player.LoggedIn = true
        err = ctx.Set(PLAYERS, players)
        if err != nil {
            return nil, err
        }
    }

    entity, err := c.getEntity(ctx, player.Id)

    return &types.PlayerDetails{Player: player, Entity: entity}, nil
}