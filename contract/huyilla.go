package contract

import (
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
    ctx.Set(AGE, &types.Age{Ticks: 1})  // starts at 1 because 0 counts as non-existent

    config := &types.Config{
        Options: &types.PrimitiveMap{
            Map: map[string]*types.Primitive{
                "PlayerCap": {Value: &types.Primitive_Int{Int: 10}},
            },
        },
    }
    ctx.Set(CONFIG, config)

    ctx.Set(PLAYERS, &types.Players{
        Players: map[string]*types.AbsolutePoint{
            "admin": {Chunk: &types.Point{0, 0, 0}, Voxel: &types.Point{0, 0, 0}},
        }})

    return nil
}
