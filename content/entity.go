package content

import "github.com/felzix/huyilla/types"

var EntityDefinitions = map[uint64]*types.EntityDefinition {
    0: {Name: "human",
        Falls: true,
        InventoryCapacity: 10},
    1: {Name: "snake",
        Falls: true},
    2: {Name: "wisp",
        Falls: false,
        InventoryCapacity: 1},
}
