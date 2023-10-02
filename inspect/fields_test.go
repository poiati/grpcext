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
