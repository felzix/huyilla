package content

import "github.com/felzix/huyilla/types"

var EntityDefinitions = map[uint64]*types.EntityDefinition{
	0: {Name: "human",
		Form: 10000, // "human"
		Material: 10000, // "human"
		Falls:             true,
		InventoryCapacity: 10},
	1: {Name: "snake",
		Form: 10001, // "snake"
		Material: 10001, // "snake"
		Falls: true},
	2: {Name: "wisp",
		Form: 10002, // "wisp"
		Material: 10002, // "snake"
		Falls:             false,
		InventoryCapacity: 1},
}
