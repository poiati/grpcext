package trace

import (
	"reflect"

	"google.golang.org/protobuf/reflect/protoreflect"
)

var supportedTypes = map[protoreflect.Kind]reflect.Kind{
	protoreflect.StringKind: reflect.String,
	protoreflect.Int32Kind:  reflect.Int32,
	protoreflect.BoolKind:   reflect.Bool,
}

type Field struct {
	Name string
	Kind reflect.Kind
}

func FieldsFor(msg protoreflect.Message) ([]Field, error) {
	descriptor := msg.Descriptor()
	fieldsDesc := descriptor.Fields()

	fields := make([]Field, fieldsDesc.Len())

	for i := 0; i < fieldsDesc.Len(); i++ {
		pbField := fieldsDesc.Get(i)

		if pbField.IsList() {
			continue
		}

		if goKind, ok := supportedTypes[pbField.Kind()]; ok {
			fields[i] = Field{
				Name: pbField.JSONName(),
				Kind: goKind,
			}
		}
	}

	return fields, nil
}
