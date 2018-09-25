package engine

type Entity struct {
    Type      EntityType
    Id        UniqueId
    Inventory Inventory
}

type EntityType uint
type UniqueId   string
type Inventory  []Item

type EntityProperty struct {
    Tags              []string `json:"tags"`
    InventoryCapacity uint     `json:"inventory_capacity"`
}


func MakeEntity (eType EntityType) *Entity {
    return &Entity{Type: eType}
}

func (e *Entity) InventoryCapacity (content *Content) uint {
    props := content.EP[e.Type]
    return props.InventoryCapacity
}