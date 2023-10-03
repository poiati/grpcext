# grpcext

grpcext is a Go package that provides extensions and utilities for working with gRPC.

## Installation

To install grpcext, you can use the `go get` command:

```
go get github.com/poiati/grpcext
```

## Usage

Here's an example of how to use the `inspect.FieldsFor` function to get the fields of a protocol buffer message:

```go
fooReq := &proto.FooRequest{
    Text:               "fooValue",
    Number:             -32,
    Flag:               true,
    LongNumber:         -64,
    UnsignedNumber:     32,
    UnsignedLongNumber: 64,
}

fields := inspect.FieldsFor(fooReq)

fmt.Println(fields)
```

This will output the following:

```
[
    {"text", reflect.String, "fooValue"},
    {"number", reflect.Int32, int32(-32)},
    {"flag", reflect.Bool, true},
    {"long_number", reflect.Int64, int64(-64)},
    {"unsigned_number", reflect.Uint32, uint32(32)},
    {"unsigned_long_number", reflect.Uint64, uint64(64)},
]
```

You can find an example in the `example` folder of this repository. It uses this package for writting an
interceptor that logs the grpc request and response.

```
$ go run example/main.go
```

```
time=2023-10-02T20:55:13.358-03:00 level=DEBUG msg="received GRPC Request" grpc.method=/proto.FooService/DoFoo grpc.req.text=foo grpc.req.number=11 grpc.req.flag=true
time=2023-10-02T20:55:13.358-03:00 level=INFO msg="Doing Foo!" grpc.method=/proto.FooService/DoFoo
time=2023-10-02T20:55:13.358-03:00 level=DEBUG msg="returning GRPC Response" grpc.method=/proto.FooService/DoFoo grpc.res.bar=bar!
```

## Benchmark

There is a simple benchmark set to compare the perf. with and without the interceptor.

```
go test -bench=. ./example/ -benchmem -benchtime=10s
```

```
goos: linux
goarch: amd64
pkg: github.com/poiati/grpcext/example
cpu: AMD Ryzen 5 5600X 6-Core Processor
BenchmarkWithoutInterceptor/grpc_call-12                  108408            112382 ns/op            8934 B/op        175 allocs/op
BenchmarkWithIterceptor/grpc_call-12                      103459            115325 ns/op           11643 B/op        225 allocs/op
PASS
ok      github.com/poiati/grpcext/example       26.414s
```

## Roadmap

For now it just supports basic fields inspection.

- Add support to nested messages
- Add support to lists (repeated)
- Add support to maps
- Add an option to exclude fields