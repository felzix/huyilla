package content

import "github.com/felzix/huyilla/types"

var FormDefinitions = map[uint64]*types.FormDefinition{
    0: {Name: "voxel"},

    1: {Name: "pile"},
    2: {Name: "ingot"},

    3: {Name: "pole"},
    4: {Name: "sphere"},
    5: {Name: "cube"},

    6: {Name: "stick"},
    7: {Name: "log"},

    8: {Name: "nail"},
    9: {Name: "leather strap"},
    10: {Name: "spearhead"},

    11: {Name: "spear",
         Wieldable2Handed: true,
         Sharpness: 50},
    12: {Name: "helm",
         WearableSlot: types.FormDefinition_WearableHead,
         PiercingProtection: 3,
         BluntProtection: 3,
         Insulation: 2},
}
