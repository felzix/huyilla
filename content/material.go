package content

import "github.com/felzix/huyilla/types"

var MaterialDefinitions = map[uint64]*types.MaterialDefinition{
	0: {Name: "adminium"},

	1: {Name: "dirt"},
	2: {Name: "water",
		PH:      7,
		SolidAt: 273,
		GasAt:   373},

	3: {Name: "siltstone"},

	4: {Name: "copper",
		Metallic: true,
		SolidAt:  1357,
		GasAt:    2835},

	5: {Name: "cow skin",
		Fleshy: true},
	6: {Name: "cow leather",
		Fleshy: true},

	7: {Name: "oakwood",
		Wooden: true},
	8: {Name: "oakbark",
		Wooden: true},

	9: {Name: "salt",
		Salty: true},
	10: {Name: "tannin",
		PH: 3},
}
