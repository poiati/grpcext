package inspect

import (
	"reflect"
	"strings"

	"google.golang.org/protobuf/reflect/protoreflect"
)

var supportedTypes = map[protoreflect.Kind]reflect.Kind{
	protoreflect.StringKind: reflect.String,
	protoreflect.Int32Kind:  reflect.Int32,
	protoreflect.Int64Kind:  reflect.Int64,
	protoreflect.Uint32Kind: reflect.Uint32,
	protoreflect.Uint64Kind: reflect.Uint64,
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
	return fieldsForMessage(msg.ProtoReflect(), "")
}

func fieldsForMessage(msg protoreflect.Message, prefix string) []Field {
	descriptor := msg.Descriptor()
	fieldsDesc := descriptor.Fields()

	fields := make([]Field, 0, fieldsDesc.Len())

	for i := 0; i < fieldsDesc.Len(); i++ {
		protoField := fieldsDesc.Get(i)

		if protoField.IsList() {
			continue
		}

		qualifiedName := string(protoField.Name())
		if prefix != "" {
			qualifiedName = strings.Join([]string{prefix, string(protoField.Name())}, ".")
		}

		if goKind, ok := supportedTypes[protoField.Kind()]; ok {
			fields = append(fields, Field{
				Name:  qualifiedName,
				Kind:  goKind,
				Value: msg.Get(protoField).Interface(),
			})
		} else {
			if protoField.Kind() != protoreflect.MessageKind {
				continue
			}

			// Handle nested messages recursively.
			fields = append(
				fields,
				fieldsForMessage(
					msg.Get(protoField).Interface().(protoreflect.Message),
					qualifiedName,
				)...,
			)
		}
	}

	return fields
}
