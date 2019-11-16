package types

type EntityController int
const (
	NPC EntityController = iota
	PLAYER
)

type Properties map[string]interface{}

type Inventory struct {
	Items []ItemId
}

type EntityId int64
type EntityType uint64
type Entity struct {
	Id EntityId
	Type EntityType
	Control EntityController
	Location AbsolutePoint
	Properties
	Inventory
	PlayerName string
}

func NewEntity(id EntityId, type_ EntityType, location AbsolutePoint) *Entity {
	return &Entity{
		Id:         id,
		Type:       type_,
		Control:    NPC,
		Location:   location,
		Properties: make(Properties, 0),
		Inventory:  Inventory{},
		PlayerName: "",
	}
}

func NewPlayerEntity(id EntityId, type_ EntityType, location AbsolutePoint, name string) *Entity {
	return &Entity{
		Id:         id,
		Type:       type_,
		Control:    PLAYER,
		Location:   location,
		Properties: make(Properties, 0),
		Inventory:  Inventory{},
		PlayerName: name,
	}
}

func (e Entity) Marshal() ([]byte, error) {
	return ToBytes(e)
}

func (e *Entity) Unmarshal(blob []byte) error {
	return FromBytes(blob, &e)
}
