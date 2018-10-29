package main

import (
    "encoding/json"
    "github.com/felzix/huyilla/content"
    "github.com/felzix/huyilla/types"
    "github.com/loomnetwork/go-loom/plugin"
    contract "github.com/loomnetwork/go-loom/plugin/contractpb"
    "github.com/pkg/errors"
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
    // So that recipes and terrain generator can reference content by name.
    content.PopulateContentNameMaps()

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

    // Add *any* entry to contract state of PLAYERS so it exists
    err = ctx.Set(PLAYERS, &types.Players{
        Players: map[string]*types.Player{
            "FAKE": {EntityId: -1,
                     Name: "FAKE",
                     Address: "FAKE",
                     LoggedIn: false}},
        })
    if err != nil { return err }

    return nil
}


func (c *Huyilla) SignUp(ctx contract.Context, req *types.PlayerName) error {
    // Make sure the player doesn't already exist
    players, err := c.getPlayers(ctx)
    if err != nil { return err }

    player := players.Players[c.myAddress(ctx)]
    if player != nil {
        return errors.New("You are already signed up.")
    }

    // Create new player
    defaultLocation := newAbsolutePoint(0, 0, 0, 0, 0, 0)

    entity := c.newEntity(ctx, ENTITY["human"], req.Name, defaultLocation)
    err = c.setEntity(ctx, entity)
    if err != nil { return errors.Wrap(err, "Entity could not be created") }

    player = &types.Player{
        EntityId: entity.Id,
        Name:     req.Name,
        Address:  c.myAddress(ctx),
        Spawn:    defaultLocation,
        LoggedIn: false,
    }
    players.Players[player.Address] = player
    err = ctx.Set(PLAYERS, players)
    if err != nil {return err}

    // Tell the client that the player was created (as well as everyone else listening, as a side effect)
    emitMsg := struct {
        Method string
        Owner  string
        Addr   string
    }{"SignUp", player.Name, player.Address}
    emitMsgJSON, err := json.Marshal(emitMsg)
    if err != nil {return err}
    ctx.EmitTopics(emitMsgJSON, "huyilla:" + string(emitMsg.Addr))

    return nil
}

func (c *Huyilla) LogIn (ctx contract.Context, req *plugin.Request) (*types.PlayerDetails, error) {
    players, err := c.getPlayers(ctx)
    if err != nil { return nil, err }

    player := players.Players[c.myAddress(ctx)]

    if player == nil {
        return nil, errors.New("You have not yet signed up")
    }

    entity, err := c.getEntity(ctx, player.EntityId)
    if err != nil { return nil, err }

    if player.LoggedIn {
        return nil, errors.New("You are already logged in.")
    }

    player.LoggedIn = true
    err = ctx.Set(PLAYERS, players)
    if err != nil { return nil, err }

    err = c.addEntityToChunk(ctx, entity)
    if err != nil { return nil, err }

    return &types.PlayerDetails{Player: player, Entity: entity}, nil
}

func (c *Huyilla) LogOut (ctx contract.Context, req *plugin.Request) error {
    players, err := c.getPlayers(ctx)
    if err != nil { return err }

    player := players.Players[c.myAddress(ctx)]

    if player.Address != c.myAddress(ctx) {
        return errors.New("Username is not associated with your address/key/account.")
    }

    if !player.LoggedIn {
        return errors.New("You are already logged out.")
    }

    player.LoggedIn = false
    err = ctx.Set(PLAYERS, players)
    if err != nil {
        return err
    }

    entity, err := c.getEntity(ctx, player.EntityId)

    err = c.removeEntityFromChunk(ctx, entity)
    if err != nil { return err }

    return nil
}

// must be Context not StaticContext because ctx.message().Sender is 0x0 under static context
func (c *Huyilla) MyAddress (ctx contract.Context, req *plugin.Request) (*types.Address, error) {
    addr := types.Address{c.myAddress(ctx)}
    return &addr, nil
}

func (c *Huyilla) myAddress (ctx contract.StaticContext) string {
    return ctx.Message().Sender.Local.String()
}