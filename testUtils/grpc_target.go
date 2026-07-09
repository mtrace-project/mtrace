package testutils

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bufbuild/protocompile"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/dynamicpb"
)

// GRPCRequestRecord records details of a gRPC request received by the mock server.
type GRPCRequestRecord struct {
	Method   string
	Metadata metadata.MD
	Request  *dynamicpb.Message
}

// GRPCTargetServer is a mock gRPC server.
type GRPCTargetServer struct {
	Address  string
	Server   *grpc.Server
	Requests []GRPCRequestRecord
}

// StartGRPCTargetServer starts a mock gRPC server on a local port.
// It uses the proto files to decode incoming requests dynamically.
func StartGRPCTargetServer(t *testing.T, baseDir string, protoPath string, serviceName string, responseData map[string]any) *GRPCTargetServer {
	t.Helper()

	importPaths := []string{".", baseDir}
	compilePath := protoPath

	if filepath.IsAbs(protoPath) {
		if rel, err := filepath.Rel(baseDir, protoPath); err == nil && !strings.HasPrefix(rel, "..") {
			compilePath = rel
		} else {
			importPaths = append(importPaths, filepath.Dir(protoPath))
			compilePath = filepath.Base(protoPath)
		}
	}

	resolver := &protocompile.SourceResolver{
		ImportPaths: importPaths,
	}
	compiler := protocompile.Compiler{
		Resolver: protocompile.WithStandardImports(resolver),
	}

	ctx := context.Background()
	allFiles, err := compiler.Compile(ctx, compilePath)
	if err != nil {
		t.Fatalf("Failed to compile proto file in mock server: %v", err)
	}
	if len(allFiles) == 0 {
		t.Fatalf("No compiled files returned for %s", protoPath)
	}
	fileDesc := allFiles[0]

	var serviceDesc protoreflect.ServiceDescriptor
	for i := 0; i < fileDesc.Services().Len(); i++ {
		s := fileDesc.Services().Get(i)
		if string(s.Name()) == serviceName || string(s.FullName()) == serviceName {
			serviceDesc = s
			break
		}
	}
	if serviceDesc == nil {
		t.Fatalf("Service %s not found in proto file", serviceName)
	}

	target := &GRPCTargetServer{
		Requests: []GRPCRequestRecord{},
	}

	// Dynamic handler for unary methods
	handler := func(srv interface{}, stream grpc.ServerStream) error {
		fullMethod, ok := grpc.MethodFromServerStream(stream)
		if !ok {
			return fmt.Errorf("failed to get method from server stream")
		}

		parts := strings.Split(fullMethod, "/")
		shortMethod := parts[len(parts)-1]

		methodDesc := serviceDesc.Methods().ByName(protoreflect.Name(shortMethod))
		if methodDesc == nil {
			return fmt.Errorf("method %s not found in service %s", shortMethod, serviceName)
		}

		// Read headers
		md, _ := metadata.FromIncomingContext(stream.Context())

		// Receive request message dynamically
		inputMsg := dynamicpb.NewMessage(methodDesc.Input())
		if err := stream.RecvMsg(inputMsg); err != nil {
			return fmt.Errorf("failed to receive request: %w", err)
		}

		target.Requests = append(target.Requests, GRPCRequestRecord{
			Method:   fullMethod,
			Metadata: md,
			Request:  inputMsg,
		})

		// Build response message dynamically
		outputMsg := dynamicpb.NewMessage(methodDesc.Output())
		jsonData, err := json.Marshal(responseData)
		if err != nil {
			return fmt.Errorf("failed to marshal output data: %w", err)
		}

		unmarshaler := protojson.UnmarshalOptions{
			DiscardUnknown: true,
			Resolver:       protoregistry.GlobalTypes,
		}
		if err := unmarshaler.Unmarshal(jsonData, outputMsg); err != nil {
			return fmt.Errorf("failed to populate dynamic response: %w", err)
		}

		// Send response
		if err := stream.SendMsg(outputMsg); err != nil {
			return fmt.Errorf("failed to send response: %w", err)
		}

		return nil
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0") // nolint:noctx
	if err != nil {
		t.Fatalf("Failed to listen on dynamic TCP port: %v", err)
	}

	server := grpc.NewServer(
		grpc.UnknownServiceHandler(handler),
	)

	go func() {
		if err := server.Serve(listener); err != nil && err != grpc.ErrServerStopped {
			t.Errorf("gRPC mock server error: %v", err)
		}
	}()

	t.Cleanup(func() {
		server.GracefulStop()
	})

	target.Address = listener.Addr().String()
	target.Server = server

	return target
}
