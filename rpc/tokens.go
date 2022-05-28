package rpc

import (
	"context"
	"errors"
	"net/http"

	"github.com/0xsequence/go-sequence/lib/prototyp"
	"github.com/0xsequence/go-sequence/metadata"
	"github.com/crypto-pe/backend/proto"
)

const (
	USDC_ADDRESS = "0x2791Bca1f2de4661ED88A30C99A7a9449Aa84174"
	DAI_ADDRESS  = "0x8f3Cf7ad23Cd3CaDbD9735AFf958023239c6A063"
)

const (
	USDC_DECIMALS = 6
	DAI_DECIMALS  = 12
)

func GetTokenAddressFromTokenType(tokenType *proto.TokenType) (string, error) {
	switch *tokenType {
	case proto.TokenType_USDC:
		return USDC_ADDRESS, nil
	case proto.TokenType_DAI:
		return DAI_ADDRESS, nil
	default:
		return "", errors.New("Invalid token type")
	}

}

func GetTokeTypeFromAddress(tokenAddress string) (proto.TokenType, error) {
	switch tokenAddress {
	case USDC_ADDRESS:
		return proto.TokenType_USDC, nil
	case DAI_ADDRESS:
		return proto.TokenType_DAI, nil
	default:
		return -1, errors.New("Invalid token address")
	}
}

func (s *RPC) GetSupportedTokens(ctx context.Context) ([]*proto.Token, error) {
	tokens := []*proto.Token{}
	client := metadata.NewMetadataClient("https://dev-metadata.sequence.app", http.DefaultClient)
	usdc, err := client.GetContractInfo(ctx, "137", USDC_ADDRESS)
	if err != nil {
		return nil, err
	}
	dai, errr := client.GetContractInfo(ctx, "137", DAI_ADDRESS)
	if errr != nil {
		return nil, err
	}

	usdcToken := proto.Token{
		Address:  prototyp.HashFromString(USDC_ADDRESS),
		Metadata: *usdc,
	}

	daiToken := proto.Token{
		Address:  prototyp.HashFromString(DAI_ADDRESS),
		Metadata: *dai,
	}

	tokens = append(tokens, &usdcToken, &daiToken)
	return tokens, nil
}
