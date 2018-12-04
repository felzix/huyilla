package types

import (
	"github.com/gogo/protobuf/proto"
	"testing"
)

func TestPrimitiveProtobuf(t *testing.T) {
	primitiveInt := &Primitive{Value: &Primitive_Int{Int: 12}}

	data, err := proto.Marshal(primitiveInt)

	if err != nil {
		t.Fatal("marshaling error: ", err)
	}
	newPrimitiveInt := &Primitive{}
	err = proto.Unmarshal(data, newPrimitiveInt)
	if err != nil {
		t.Fatal("unmarshaling error: ", err)
	}

	if primitiveInt.GetInt() != newPrimitiveInt.GetInt() {
		t.Fatalf("data mismatch %v != %v", primitiveInt.GetInt(), newPrimitiveInt.GetInt())
	}
}

func TestAgeProtobuf(t *testing.T) {
	age := &Age{Ticks: 3}

	data, err := proto.Marshal(age)

	if err != nil {
		t.Fatal("marshaling error: ", err)
	}
	newAge := &Age{}
	err = proto.Unmarshal(data, newAge)
	if err != nil {
		t.Fatal("unmarshaling error: ", err)
	}

	if age.GetTicks() != newAge.GetTicks() {
		t.Fatalf("data mismatch %v != %v", age.GetTicks(), newAge.GetTicks())
	}
}
