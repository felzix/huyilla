package types

import "bytes"
import "github.com/davecgh/go-xdr/xdr2"

type Serializable interface {
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
}

func FromBytes(input []byte, thing interface{}) error {
	if _, err := xdr.Unmarshal(bytes.NewReader(input), &thing); err != nil {
		return err
	}
	return nil
}

func ToBytes(thing interface{}) ([]byte, error) {
	var buffer bytes.Buffer
	if _, err := xdr.Marshal(&buffer, &thing); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
