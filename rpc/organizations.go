package rpc

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/0xsequence/go-sequence/lib/prototyp"
	"github.com/crypto-pe/backend/data"
	"github.com/crypto-pe/backend/data/sqlc"
	"github.com/crypto-pe/backend/proto"
	"github.com/crypto-pe/backend/rpc/middleware"
	"github.com/google/uuid"
)

func (s *RPC) CreateOrganization(ctx context.Context, name string, tokenType *proto.TokenType) (bool, *proto.Organization, error) {
	user, ok := ctx.Value(middleware.UserCtxKey).(sqlc.Accounts)
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
			Token:        tokenAddress, // create a constant method
		},
	)
	if err != nil {
		s.Log.Err(err).Msg("unable to create account")
		if strings.Contains(err.Error(), "SQLSTATE 23505") {
			return false, nil, proto.Errorf(proto.ErrAlreadyExists, "org already exists")
		}
		return false, nil, proto.WrapError(proto.ErrInternal, err, "unable to create org")
	}

	_, err = data.DB.CreateOrganizationMember(ctx, sqlc.CreateOrganizationMemberParams{
		OrganizationID: org.ID,
		MemberAddress:  user.Address,
		Role:           "owner",
		IsAdmin:        sql.NullBool{Bool: true, Valid: true},
		Salary: sql.NullString{
			String: fmt.Sprintf("%d", 0),
			Valid:  true,
		},
	})

	if err != nil {
		s.Log.Err(err).Msg("unable to create account")
	}

	respOrg := proto.Organization{
		Id:           org.ID.String(),
		Name:         org.Name,
		CreatedAt:    &org.CreatedAt,
		OwnerAddress: prototyp.HashFromString(org.OwnerAddress),
		Token:        tokenType, // reverse here
	}

	return true, &respOrg, nil
}

// create org done

func (s *RPC) GetOrganization(ctx context.Context, organizationID string) (*proto.Organization, error) {
	organizationUUID, err := uuid.FromBytes([]byte(organizationID))
	if err != nil {
		s.Log.Err(err).Msg("Invalid UUID provided")
		return nil, proto.WrapError(proto.ErrInvalidArgument, err, "Invalid UUID provided")
	}

	dbOrganization, err := data.DB.GetOrganization(ctx, organizationUUID)
	if err != nil {
		s.Log.Err(err).Msg("Organization does not exist.")
		return nil, proto.WrapError(proto.ErrNotFound, err, "Organization does not exist.")
	}

	tokenType, err := GetTokeTypeFromAddress(prototyp.HashFromString(dbOrganization.Token).String())

	if err != nil {
		s.Log.Err(err).Msg(prototyp.HashFromString(dbOrganization.Token).String() + " cannot find token")
		return nil, proto.WrapError(proto.ErrInternal, err, "could not get org")
	}

	return &proto.Organization{
		Id:           organizationUUID.String(),
		Name:         dbOrganization.Name,
		CreatedAt:    &dbOrganization.CreatedAt,
		OwnerAddress: prototyp.Hash(dbOrganization.OwnerAddress),
		Token:        &tokenType,
	}, nil
}

func (s *RPC) UpdateOrganization(ctx context.Context, organization *proto.Organization) (bool, *proto.Organization, error) {
	organizationUuid, err := uuid.FromBytes([]byte(organization.Id))
	if err != nil {
		s.Log.Err(err).Msg("Invalid UUID provided")
		return false, nil, proto.WrapError(proto.ErrInvalidArgument, err, "Invalid UUID provided")
	}

	isAdmin, err := s.checkOrgAdmin(ctx, organizationUuid)
	if err != nil {
		return false, nil, err
	}

	if !isAdmin {
		return false, nil, proto.WrapError(proto.ErrPermissionDenied, errors.New("not an admin"), "not an admin")
	}

	dbOrg, err := data.DB.UpdateOrganization(ctx, sqlc.UpdateOrganizationParams{
		ID:    organizationUuid,
		Name:  organization.Name,
		Token: organization.Token.String(),
	})

	if err != nil {
		s.Log.Err(err).Msg("Could not update organization")
		return false, nil, proto.WrapError(proto.ErrInternal, err, "Could not update organization")
	}
	return true, &proto.Organization{
		Id:           organizationUuid.String(),
		Name:         dbOrg.Name,
		CreatedAt:    &dbOrg.CreatedAt,
		OwnerAddress: prototyp.Hash(dbOrg.OwnerAddress),
		Token:        organization.Token,
	}, nil
}

func (s *RPC) DeleteOrganization(ctx context.Context, organizationID string) (bool, error) {
	organizationUuid, err := uuid.FromBytes([]byte(organizationID))
	if err != nil {
		s.Log.Err(err).Msg("Invalid UUID provided")
		return false, proto.WrapError(proto.ErrInvalidArgument, err, "Invalid UUID provided")
	}
	isAdmin, err := s.checkOrgAdmin(ctx, organizationUuid)
	if err != nil {
		return false, err
	}

	if !isAdmin {
		return false, proto.WrapError(proto.ErrPermissionDenied, errors.New("not an admin"), "not an admin")
	}

	err = data.DB.DeleteOrganization(ctx, organizationUuid)
	if err != nil {
		return false, proto.WrapError(proto.ErrInternal, err, "Could not delete the organization.")
	}

	return true, nil

}

func (s *RPC) GetAllOrganizations(ctx context.Context) ([]*proto.Organization, error) {
	user, ok := ctx.Value(middleware.UserCtxKey).(sqlc.Accounts)
	// wont fail cause we ensure in middleware, ok
	if !ok {
		s.Log.Err(errors.New("User does not exist")).Msg("Could not get user.")
		return nil, proto.WrapError(proto.ErrPermissionDenied, errors.New("User does not exist"), "Could not get user")
	}

	orgs, err := data.DB.GetAllOrganizations(ctx, prototyp.HashFromString(user.Address).String())
	if err != nil {
		s.Log.Err(err).Msg("Could not get organizations")
		return nil, proto.WrapError(proto.ErrInternal, err, "Could not get organizations")
	}
	fmt.Println("orgs", orgs, prototyp.HashFromString(user.Address))
	// 0x307865306339383238646565333431316132386363623462623832613138643061616432343438396530
	resultOrgs := make([]*proto.Organization, len(orgs))
	for i, org := range orgs {
		tokenType, _ := GetTokeTypeFromAddress(org.Token)
		resultOrgs[i] = &proto.Organization{
			Id:           org.ID.String(),
			Name:         org.Name,
			CreatedAt:    &org.CreatedAt,
			OwnerAddress: prototyp.Hash(org.OwnerAddress),
			Token:        &tokenType,
		}
	}
	return resultOrgs, nil
}

func (s *RPC) checkOrgAdmin(ctx context.Context, organizationID uuid.UUID) (bool, error) {
	orgMember, err := data.DB.GetOrganizationMember(ctx, sqlc.GetOrganizationMemberParams{
		OrganizationID: organizationID,
		MemberAddress:  ctx.Value(middleware.WalletCtxKey).(string),
	})
	if err != nil {
		return false, proto.WrapError(proto.ErrInternal, err, "Could not get org member")
	}

	if !orgMember.IsAdmin.Bool {
		return false, nil
	}

	return true, nil
}
