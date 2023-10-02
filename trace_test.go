package trace_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"trace"
	"trace/gen/proto"
)

func TestFieldsFor(t *testing.T) {
	fooReq := &proto.FooRequest{
		Text:       "fooValue",
		Number:     11,
		Flag:       true,
		LongNumber: 64,
	}

	fields := trace.FieldsFor(fooReq)

	assert.Equal(
		t,
		[]trace.Field{
			{"text", reflect.String, "fooValue"},
			{"number", reflect.Int32, int32(11)},
			{"flag", reflect.Bool, true},
			{"long_number", reflect.Int64, int64(64)},
		},
		fields,
	)
}
