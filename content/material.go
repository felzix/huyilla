package content

import "github.com/felzix/huyilla/types"

// 16 bits available: 0 - 65535
var MaterialDefinitions = types.MaterialDefinitions{
	0: {Name: "air",
		SolidAt: uint32(types.MaxTemperature),
		GasAt:   uint32(types.MinTemperature)},

	100: {Name: "dirt"},
	101: {Name: "silt"},
	102: {Name: "grass"},

	200: {
		Name:    "water",
		Ph:      7,
		SolidAt: 273,
		GasAt:   373},

	1000: {Name: "quartz"},
	1001: {Name: "feldspar"},
	1002: {Name: "mica"},
	1003: {Name: "salt"},

	2000: {
		Name:     "copper",
		Metallic: true,
		SolidAt:  1357,
		GasAt:    2835},

	3000: {
		Name:   "cow skin",
		Fleshy: true},
	3001: {
		Name:   "cow leather",
		Fleshy: true},

	4000: {
		Name:   "oakwood",
		Wooden: true},
	4001: {
		Name:   "oakbark",
		Wooden: true},

	5000: {Name: "tannin",
		Ph: 3},

	10000: {Name: "human"},
	10001: {Name: "snake"},
	10002: {Name: "wisp"},
}
