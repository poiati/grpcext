package main

import (
	"context"
	"log/slog"
	"net"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/poiati/grpcext/gen/proto"
)

func BenchmarkWithoutInterceptor(b *testing.B) {
	slog.SetDefault(
		slog.New(
			slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}),
		),
	)

	server := grpc.NewServer()
	proto.RegisterFooServiceServer(server, &fooServer{})

	lis, err := net.Listen("tcp", ":8085")
	assert.NoError(b, err)

	go func() {
		assert.NoError(b, server.Serve(lis))
	}()
	defer server.Stop()

	conn, err := grpc.Dial(":8085", grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(b, err)

	client := proto.NewFooServiceClient(conn)

	b.Run("grpc call", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			client.DoFoo(context.Background(), &proto.FooRequest{
				Text:   "foo",
				Number: 11,
				Flag:   true,
			})
			assert.NoError(b, err)
		}
	})
}

func BenchmarkWithIterceptor(b *testing.B) {
	logger := slog.New(
		slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}),
	)

	server := grpc.NewServer(
		grpc.UnaryInterceptor(loggerInterceptor(logger)),
	)
	proto.RegisterFooServiceServer(server, &fooServer{})

	lis, err := net.Listen("tcp", ":8085")
	assert.NoError(b, err)

	go func() {
		assert.NoError(b, server.Serve(lis))
	}()
	defer server.Stop()

	conn, err := grpc.Dial(":8085", grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(b, err)

	client := proto.NewFooServiceClient(conn)

	b.Run("grpc call", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := client.DoFoo(context.Background(), &proto.FooRequest{
				Text:   "foo",
				Number: 11,
				Flag:   true,
			})
			assert.NoError(b, err)
		}
	})
}
