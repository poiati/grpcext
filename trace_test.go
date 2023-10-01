package trace_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"trace"
	"trace/gen/proto"
)

func TestFieldsFor(t *testing.T) {
	fooReq := proto.FooRequest{
		Foo: "fooValue",
	}

	fields, err := trace.FieldsFor(fooReq.ProtoReflect())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.Equal(
		t,
		[]trace.Field{
			{"foo", reflect.String},
			{"number", reflect.Int32},
			{"flag", reflect.Bool},
		},
		fields,
	)
}
