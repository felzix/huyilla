package types

type Age uint64

func NewAge(initial uint64) Age {
	return Age(initial)
}

func (a *Age) Increment() *Age {
	*a++
	return a
}

func (a Age) Marshal() ([]byte, error) {
	return ToBytes(a)
}

func (a *Age) Unmarshal(blob []byte) error {
	return FromBytes(blob, &a)
}
