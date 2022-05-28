package rpc

import (
	"context"

	"github.com/crypto-pe/api"
	"github.com/crypto-pe/api/proto"
)

// Ping is a healthcheck that returns an empty message.
func (s *RPC) Ping(ctx context.Context) (bool, error) {
	return true, nil
}

// Version returns service version details
func (s *RPC) Version(ctx context.Context) (*proto.Version, error) {
	return &proto.Version{
		WebrpcVersion: proto.WebRPCVersion(),
		SchemaVersion: proto.WebRPCSchemaVersion(),
		SchemaHash:    proto.WebRPCSchemaHash(),
		AppVersion:    api.GITCOMMIT,
	}, nil
}
