package engine


type Entity struct {
    Type EntityType
    Inventory []Item
}

type EntityType uint
type Inventory []Item


func MakeEntity (eType EntityType) *Entity {
    return &Entity{Type: eType}
}

func (e *Entity) InventoryCapacity (content *Content) uint {
    props := content.EP[e.Type]
    return props.InventoryCapacity
}