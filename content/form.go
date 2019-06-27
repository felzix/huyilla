package content

import "github.com/felzix/huyilla/types"

var FormDefinitions = map[uint64]*types.FormDefinition{
	// Forms typical of voxels

	0: {Name: "cube"},

	// partial blocks assume the nearby material fills in the rest (typically that's air)
	1: {Name: "slab, bottom"},
	2: {Name: "slab, top"},
	3: {Name: "slab, north"},
	4: {Name: "slab, south"},
	5: {Name: "slab, east"},
	6: {Name: "slab, west"},

	7:  {Name: "stairs, floor, north"},
	8:  {Name: "stairs, floor, south"},
	9:  {Name: "stairs, floor, east"},
	10: {Name: "stairs, floor, west"},
	11: {Name: "stairs, floor, northeast"},
	12: {Name: "stairs, floor, northwest"},
	13: {Name: "stairs, floor, southeast"},
	14: {Name: "stairs, floor, southwest"},
	15: {Name: "stairs, ceiling, north"},
	16: {Name: "stairs, ceiling, south"},
	17: {Name: "stairs, ceiling, east"},
	18: {Name: "stairs, ceiling, west"},
	19: {Name: "stairs, ceiling, northeast"},
	20: {Name: "stairs, ceiling, northwest"},
	21: {Name: "stairs, ceiling, southeast"},
	22: {Name: "stairs, ceiling, southwest"},

	23: {Name: "ramp, floor, north"},
	24: {Name: "ramp, floor, south"},
	25: {Name: "ramp, floor, east"},
	26: {Name: "ramp, floor, west"},
	27: {Name: "ramp, floor, northeast"},
	28: {Name: "ramp, floor, northwest"},
	29: {Name: "ramp, floor, southeast"},
	30: {Name: "ramp, floor, southwest"},
	31: {Name: "ramp, ceiling, north"},
	32: {Name: "ramp, ceiling, south"},
	33: {Name: "ramp, ceiling, east"},
	34: {Name: "ramp, ceiling, west"},
	35: {Name: "ramp, ceiling, northeast"},
	36: {Name: "ramp, ceiling, northwest"},
	37: {Name: "ramp, ceiling, southeast"},
	38: {Name: "ramp, ceiling, southwest"},

	39: {Name: "pole, vertical"},
	40: {Name: "pole, north-south"},
	41: {Name: "pole, east-west"},

	// relevant for solids but not fluids
	32000: {Name: "big chunks"},
	32001: {Name: "large gravel"},
	32002: {Name: "gravel"},
	32003: {Name: "small gravel"},
	32004: {Name: "gritty sand"},
	32005: {Name: "sand"},
	32006: {Name: "fine sand"},
	32007: {Name: "dust"},

	// Forms typical of items

	32767: {Name: "pile"},
	32768: {Name: "ingot"},

	32800: {Name: "pole"},
	32801: {Name: "sphere"},
	32802: {Name: "cube"},

	32900: {Name: "stick"},
	32901: {Name: "log"},

	33000: {Name: "nail"},
	33001: {Name: "leather strap"},
	33002: {Name: "spearhead"},

	34000: {Name: "spear",
		Wieldable2Handed: true,
		Sharpness:        50},

	35000: {Name: "helm",
		WearableSlot:       types.FormDefinition_WearableHead,
		PiercingProtection: 3,
		BluntProtection:    3,
		Insulation:         2},

	36000: {Name: "shovel",
		Wieldable2Handed: true,
		Digging:          90},
}
