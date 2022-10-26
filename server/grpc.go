package server

import "google.golang.org/grpc"

func NewGrpc(opts ...grpc.ServerOption) *grpc.Server {
	return grpc.NewServer(opts...)
}
