package rpc

import (
	"context"
	"net/http"

	"github.com/0xsequence/go-sequence/lib/prototyp"
	"github.com/0xsequence/go-sequence/metadata"
	"github.com/crypto-pe/backend/proto"
)

var USDCAddress = "0x2791Bca1f2de4661ED88A30C99A7a9449Aa84174"
var DaiAddress = "0x8f3Cf7ad23Cd3CaDbD9735AFf958023239c6A063"

func (s *RPC) GetSupportedTokens(ctx context.Context) ([]*proto.Token, error) {
	tokens := []*proto.Token{}
	client := metadata.NewMetadataClient("https://dev-metadata.sequence.app", http.DefaultClient)
	usdc, err := client.GetContractInfo(ctx, "137", USDCAddress)
	if err != nil {
		return nil, err
	}
	dai, errr := client.GetContractInfo(ctx, "137", DaiAddress)
	if errr != nil {
		return nil, err
	}

	usdcToken := proto.Token{
		Address:  prototyp.HashFromString(USDCAddress),
		Metadata: *usdc,
	}

	daiToken := proto.Token{
		Address:  prototyp.HashFromString(DaiAddress),
		Metadata: *dai,
	}

	tokens = append(tokens, &usdcToken, &daiToken)
	return tokens, nil
}
