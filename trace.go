package trace

import (
	"reflect"

	"google.golang.org/protobuf/reflect/protoreflect"
)

var supportedTypes = map[protoreflect.Kind]reflect.Kind{
	protoreflect.StringKind: reflect.String,
	protoreflect.Int32Kind:  reflect.Int32,
	protoreflect.Int64Kind:  reflect.Int64,
	protoreflect.BoolKind:   reflect.Bool,
}

// Field represents a field of a proto message.
type Field struct {
	Name  string
	Kind  reflect.Kind
	Value any
}

// FieldsFor returns a slice of fields for the given proto message.
// Not all proto types are supported by default.
func FieldsFor(msg protoreflect.ProtoMessage) []Field {
	protoReflect := msg.ProtoReflect()
	descriptor := protoReflect.Descriptor()
	fieldsDesc := descriptor.Fields()

	fields := make([]Field, 0, fieldsDesc.Len())

	for i := 0; i < fieldsDesc.Len(); i++ {
		protoField := fieldsDesc.Get(i)

		if protoField.IsList() {
			continue
		}

		if goKind, ok := supportedTypes[protoField.Kind()]; ok {
			fields = append(fields, Field{
				Name:  string(protoField.Name()),
				Kind:  goKind,
				Value: protoReflect.Get(protoField).Interface(),
			})
		}
	}

	return fields
}
