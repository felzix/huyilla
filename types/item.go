package types

type Form uint64
type Material uint64

type ItemId int64
type Item struct {
	Id ItemId
	Form
	Substance
	Properties
	Inventory
}
type Substance struct {
	Material
	Components
}
type Components map[Form]Item

func NewSimpleItem(id ItemId, form Form, material Material) *Item {
	return &Item{
		Id:   id,
		Form: form,
		Substance: Substance{Material: material},
	}
}

func (i Item) Marshal() ([]byte, error) {
	return ToBytes(i)
}

func (i *Item) Unmarshal(blob []byte) error {
	return FromBytes(blob, &i)
}
