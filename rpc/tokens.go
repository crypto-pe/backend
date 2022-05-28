package rpc

import (
	"context"

	"github.com/crypto-pe/backend/proto"
)

var USDCAddress = "0x2791Bca1f2de4661ED88A30C99A7a9449Aa84174"
var DaiAddress = "0x2791Bca1f2de4661ED88A30C99A7a9449Aa84174"

func (s *RPC) GetSupportedTokens(ctx context.Context) ([]*proto.Token, error) {
	tokens := []*proto.Token{}

	return tokens, nil
}
