package types

type Auth struct {
	Name string
	Password []byte
}

func NewAuth(name string, password []byte) Auth {
	return Auth{
		Name: name,
		Password: password,
	}
}

func (a Auth) Marshal() ([]byte, error) {
	return ToBytes(a)
}

func (a *Auth) Unmarshal(blob []byte) error {
	return FromBytes(blob, &a)
}
