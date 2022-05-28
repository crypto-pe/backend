package rpc

import (
	"context"

	api "github.com/crypto-pe/backend"
	"github.com/crypto-pe/backend/proto"
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
