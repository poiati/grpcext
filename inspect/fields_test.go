package inspect_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/poiati/grpcext/gen/proto"
	"github.com/poiati/grpcext/inspect"
)

func TestFieldsFor(t *testing.T) {
	fooReq := &proto.FooRequest{
		Text:               "fooValue",
		Number:             -32,
		Flag:               true,
		LongNumber:         -64,
		UnsignedNumber:     32,
		UnsignedLongNumber: 64,
	}

	fields := inspect.FieldsFor(fooReq)

	assert.Equal(
		t,
		[]inspect.Field{
			{"text", reflect.String, "fooValue"},
			{"number", reflect.Int32, int32(-32)},
			{"flag", reflect.Bool, true},
			{"long_number", reflect.Int64, int64(-64)},
			{"unsigned_number", reflect.Uint32, uint32(32)},
			{"unsigned_long_number", reflect.Uint64, uint64(64)},
		},
		fields,
	)
}

func TestFieldsForUnsupportedType(t *testing.T) {
	fooReq := &proto.FooRequest{
		Text: "the field below in unsupported",
		Data: []byte("data here"),
	}

	fields := inspect.FieldsFor(fooReq)

	assert.Equal(
		t,
		[]inspect.Field{
			{"text", reflect.String, "the field below in unsupported"},
			{"number", reflect.Int32, int32(0)},
			{"flag", reflect.Bool, false},
			{"long_number", reflect.Int64, int64(0)},
			{"unsigned_number", reflect.Uint32, uint32(0)},
			{"unsigned_long_number", reflect.Uint64, uint64(0)},
		},
		fields,
	)
}

func TestFieldsForNestedMessage(t *testing.T) {
	msg := &proto.ComplexMesasge{
		Nested: &proto.NestedMessage{
			NestedText: "nested text",
			DeeperNested: &proto.DeeperNestedMessage{
				DeeperNestedText: "deeper nested text",
			},
		},
	}

	fields := inspect.FieldsFor(msg)

	assert.Equal(
		t,
		[]inspect.Field{
			{"flag", reflect.Bool, false},
			{"nested.nested_text", reflect.String, "nested text"},
			{"nested.deeper_nested.deeper_nested_text", reflect.String, "deeper nested text"},
		},
		fields,
	)
}
