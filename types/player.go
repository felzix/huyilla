package types

type Player struct {
	Name string
	EntityId EntityId
	Password []byte
	Token string
	Spawn AbsolutePoint
}

func NewPlayer(name string, password []byte, entityId EntityId, spawn AbsolutePoint) *Player {
	return &Player{
		Name:     name,
		EntityId: entityId,
		Password: password,
		Token:    "",
		Spawn:    spawn,
	}
}

func (p Player) Marshal() ([]byte, error) {
	return ToBytes(p)
}

func (p *Player) Unmarshal(blob []byte) error {
	return FromBytes(blob, &p)
}
