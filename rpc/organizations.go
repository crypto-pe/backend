package rpc

import (
	"context"
	"errors"
	"strings"

	"github.com/0xsequence/go-sequence/lib/prototyp"
	"github.com/crypto-pe/backend/data"
	"github.com/crypto-pe/backend/data/sqlc"
	"github.com/crypto-pe/backend/proto"
	"github.com/crypto-pe/backend/rpc/middleware"
	"github.com/google/uuid"
)

func (s *RPC) CreateOrganization(ctx context.Context, name string, tokenType *proto.TokenType) (bool, *proto.Organization, error) {
	user, ok := ctx.Value(middleware.UserCtxKey).(*sqlc.Accounts)
	// wont fail cause we ensure in middleware, ok
	if !ok {
		s.Log.Err(errors.New("User does not exist")).Msg("Could not get user.")
		return false, nil, proto.WrapError(proto.ErrPermissionDenied, errors.New("User does not exist"), "Could not get user")
	}
	tokenAddress, err := GetTokenAddressFromTokenType(tokenType)
	if err != nil {
		return false, nil, proto.WrapError(proto.ErrInvalidArgument, err, "could not create organization")
	}
	// create org here
	org, err := data.DB.CreateOrganization(
		ctx,
		sqlc.CreateOrganizationParams{
			Name:         name,
			OwnerAddress: user.Address,
			Token:        []byte(tokenAddress), // create a constant method
		},
	)
	if err != nil {
		s.Log.Err(err).Msg("unable to create account")
		if strings.Contains(err.Error(), "SQLSTATE 23505") {
			return false, nil, proto.Errorf(proto.ErrAlreadyExists, "org already exists")
		}
		return false, nil, proto.WrapError(proto.ErrInternal, err, "unable to create org")
	}

	respOrg := proto.Organization{
		Id:           org.ID.String(),
		Name:         org.Name,
		CreatedAt:    &org.CreatedAt,
		OwnerAddress: prototyp.HashFromBytes(org.OwnerAddress),
		Token:        tokenType, // reverse here
	}

	return true, &respOrg, nil
}

// create org done

func (s *RPC) GetOrganization(ctx context.Context, organization_uuid_string string) (*proto.Organization, error) {
	organization_uuid, uuiderr := uuid.FromBytes([]byte(organization_uuid_string))
	if uuiderr != nil {
		s.Log.Err(uuiderr).Msg("Invalid UUID provided")
		return nil, proto.WrapError(proto.ErrInvalidArgument, uuiderr, "Invalid UUID provided")
	}
	dbOrganization, err := data.DB.GetOrganization(ctx, organization_uuid)
	if err != nil {
		s.Log.Err(err).Msg("Organization does not exist.")
		return nil, proto.WrapError(proto.ErrNotFound, err, "Organization does not exist.")
	}

	tokenType, err := GetTokeTypeFromAddress(prototyp.HashFromBytes(dbOrganization.Token).String())

	if err != nil {
		s.Log.Err(err).Msg(prototyp.HashFromBytes(dbOrganization.Token).String() + " cannot find token")
		return nil, proto.WrapError(proto.ErrInternal, err, "could not get org")
	}

	return &proto.Organization{
		Id:           organization_uuid_string,
		Name:         dbOrganization.Name,
		CreatedAt:    &dbOrganization.CreatedAt,
		OwnerAddress: prototyp.Hash(dbOrganization.OwnerAddress),
		Token:        &tokenType,
	}, nil
}

func (s *RPC) UpdateOrganization(ctx context.Context, organization *proto.Organization) (bool, *proto.Organization, error) {
	account, ok := ctx.Value(middleware.UserCtxKey).(sqlc.Accounts)
	if !ok {
		return false, nil, proto.WrapError(proto.ErrInternal, errors.New("could not get account"), "could not update organization")
	}

	data.DB.GetOrganization(ctx, organization.Id)
	// get?

	// check if the guy is an owner or something, okok
	// wait we forgot something
	// while creating org we should create orgMember for the owner
	// and make him admin
	// and getorg member here, instead of org

	// admin is for site admins right?
	// but we can do that :mhm:
	// yea but org owner is also admin
}

func (s *RPC) DeleteOrganization(ctx context.Context, organization_uuid_string string) (bool, error) {
	organization_uuid, uuiderr := uuid.FromBytes([]byte(organization_uuid_string))
	if uuiderr != nil {
		s.Log.Err(uuiderr).Msg("Invalid UUID provided")
		return false, proto.WrapError(proto.ErrInvalidArgument, uuiderr, "Invalid UUID provided")
	}

	err := data.DB.DeleteOrganization(ctx, organization_uuid)
	if err != nil {
		return false, proto.WrapError(proto.ErrInternal, err, "Could not delete the organization.")
	}

	return true, nil

}
