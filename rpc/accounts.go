package rpc

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/0xsequence/go-ethauth"
	"github.com/0xsequence/go-sequence/lib/prototyp"
	"github.com/crypto-pe/backend/data"
	"github.com/crypto-pe/backend/data/sqlc"
	"github.com/crypto-pe/backend/proto"
)

func (s *RPC) CreateAccount(ctx context.Context, ethAuthProofString string, name string, email string) (bool, string, *proto.Account, error) {
	var jwtToken string

	ethAuth, err := ethauth.New()
	if err != nil {
		return false, "", nil, err
	}

	valid, proof, err := ethAuth.DecodeProof(ethAuthProofString)
	if err != nil {
		return false, "", nil, proto.WrapError(proto.ErrPermissionDenied, err, "invalid ethauth proof")
	}
	if !valid || proof == nil {
		return false, "", nil, proto.Errorf(proto.ErrPermissionDenied, "invalid ethauth proof")
	}

	// LATER
	// // Validate the origin in the proof claims against the http request origin header
	// if proof.Claims.Origin != "" {
	// 	httpReq, _ := ctx.Value(proto.HTTPRequestCtxKey).(*http.Request)
	// 	if httpReq.Header.Get("Origin") != proof.Claims.Origin {
	// 		return false, "", "", nil, proto.Errorf(proto.ErrInvalidArgument, "ethauth proof origin does not match the http request")
	// 	}
	// }

	jwtClaims := map[string]interface{}{
		"account": strings.ToLower(proof.Address),
		"iat":     time.Now().Unix(),
		"exp":     proof.Claims.ExpiresAt,
		"app":     proof.Claims.App,
	}

	if proof.Claims.IssuedAt != 0 {
		jwtClaims["iat"] = proof.Claims.IssuedAt
	}
	// if proof.Claims.Origin != "" {
	// 	jwtClaims["ogn"] = proof.Claims.Origin
	// }
	_, jwtToken, err = s.JWTAuth.Encode(jwtClaims)
	if err != nil {
		return false, "", nil, proto.Errorf(proto.ErrPermissionDenied, "unable to create jwt")
	}

	addr := prototyp.HashFromString(proof.Address).String()
	fmt.Println("addr", addr)
	// create account
	dbAccount := sqlc.CreateUserParams{
		Name:    name,
		Email:   email,
		Address: []byte(addr),
	}

	account, err := data.DB.CreateUser(ctx, dbAccount)
	if err != nil {
		s.Log.Err(err).Msg("unable to create account")
		return false, "", nil, proto.WrapError(proto.ErrInternal, err, "unable to create account")
	}

	responseAccount := &proto.Account{
		Address:   prototyp.HashFromBytes(account.Address),
		Name:      account.Name,
		Email:     account.Email.(string),
		CreatedAt: &account.CreatedAt.Time,
	}

	return true, jwtToken, responseAccount, nil
}

func (s *RPC) Login(ctx context.Context, ethAuthProofString string) (string, *proto.Account, error) {
	return "", &proto.Account{}, nil
}
