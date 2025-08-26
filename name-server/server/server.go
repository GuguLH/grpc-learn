package server

import (
	"context"
	"grpc-learn/name"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NameServer struct {
	name.UnimplementedNameServer
}

func (NameServer) Register(ctx context.Context, in *name.NameRequest) (*name.NameResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (NameServer) Delete(ctx context.Context, in *name.NameRequest) (*name.NameResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (NameServer) Keepalive(stream name.Name_KeepaliveServer) error {
	return status.Errorf(codes.Unimplemented, "method Keepalive not implemented")
}
func (NameServer) GetAddress(ctx context.Context, in *name.NameRequest) (*name.NameResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAddress not implemented")
}
