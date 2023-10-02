package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"reflect"
	"time"

	"github.com/poiati/grpcext/gen/proto"
	"github.com/poiati/grpcext/inspect"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type loggerKey string

const logKey loggerKey = "logKey"
const ListenPort = 8085

func LoggerFromContext(ctx context.Context) *slog.Logger {
	logger := ctx.Value(logKey)

	if logger != nil {
		return logger.(*slog.Logger)
	}

	return slog.Default()
}

type fooServer struct {
	proto.UnimplementedFooServiceServer
}

func (s *fooServer) DoFoo(ctx context.Context, req *proto.FooRequest) (*proto.FooResponse, error) {
	logger := LoggerFromContext(ctx)
	logger.Info("Doing Foo!")

	return &proto.FooResponse{Bar: "bar!"}, nil
}

func loggerWith(logger *slog.Logger, protoMsg protoreflect.ProtoMessage) *slog.Logger {
	grpcFields := inspect.FieldsFor(protoMsg)

	var logAttrs []any = make([]any, 0, len(grpcFields))

	for _, field := range grpcFields {
		switch field.Kind {
		case reflect.String:
			logAttrs = append(logAttrs, slog.String(field.Name, field.Value.(string)))
		case reflect.Int32:
			logAttrs = append(logAttrs, slog.Int(field.Name, int(field.Value.(int32))))
		case reflect.Bool:
			logAttrs = append(logAttrs, slog.Bool(field.Name, field.Value.(bool)))
		}
	}

	return logger.With(logAttrs...)
}

// loggerInterceptor returns a grpc.UnaryServerInterceptor that logs the GRPC request and response using the slog package.
func loggerInterceptor(rootLogger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		logger := rootLogger.WithGroup("grpc")
		logger = logger.With(slog.String("method", info.FullMethod))

		reqLogger := loggerWith(logger.WithGroup("req"), req.(protoreflect.ProtoMessage))
		reqLogger.Log(ctx, slog.LevelDebug, "received GRPC Request")

		res, err := handler(context.WithValue(ctx, logKey, logger), req)

		if err != nil {
			errCode := status.Code(err)

			logger.LogAttrs(
				ctx,
				slog.LevelError,
				"failed to handle GRPC Request",
				slog.String("error", err.Error()),
				slog.Int("errCode", int(errCode)),
			)
		} else {
			resLogger := loggerWith(logger.WithGroup("res"), res.(protoreflect.ProtoMessage))
			resLogger.Log(ctx, slog.LevelDebug, "returning GRPC Response")
		}

		return res, err
	}
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))

	server := grpc.NewServer(
		grpc.UnaryInterceptor(loggerInterceptor(logger)),
	)

	proto.RegisterFooServiceServer(server, &fooServer{})

	address := fmt.Sprintf(":%d", ListenPort)

	lis, err := net.Listen("tcp", address)
	if err != nil {
		logger.Error("failed to listen: %v", err)
		os.Exit(1)
	}

	go func() {
		logger.Error(server.Serve(lis).Error())
		os.Exit(1)
	}()
	defer server.GracefulStop()

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error("failed to GRPC dial: %v", err)
		os.Exit(1)
	}

	client := proto.NewFooServiceClient(conn)

	for {
		client.DoFoo(context.Background(), &proto.FooRequest{Text: "foo", Number: 11, Flag: true})
		time.Sleep(2 * time.Second)
	}
}
