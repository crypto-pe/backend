package rpc

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/0xsequence/go-sequence/lib/prototyp"
	"github.com/crypto-pe/backend/data"
	"github.com/crypto-pe/backend/data/sqlc"
	"github.com/crypto-pe/backend/proto"
	"github.com/google/uuid"
)

func (s *RPC) CreateOrganizationMember(ctx context.Context, organizationID string, memberAddress string, role string, isAdmin bool, salary int) (bool, *proto.OrganizationMember, error) {
	organizationUuid, err := uuid.FromBytes([]byte(organizationID))
	if err != nil {
		s.Log.Err(err).Msg("Invalid UUID provided")
		return false, nil, proto.WrapError(proto.ErrInvalidArgument, err, "Invalid UUID provided")
	}

	isCurrentUserAdmin, err := s.checkOrgAdmin(ctx, organizationUuid)
	if err != nil {
		return false, nil, err
	}
	if !isCurrentUserAdmin {
		return false, nil, proto.WrapError(proto.ErrPermissionDenied, errors.New("not an admin"), "not an admin")
	}
	orgMember, err := data.DB.CreateOrganizationMember(ctx, sqlc.CreateOrganizationMemberParams{
		OrganizationID: organizationUuid,
		MemberAddress:  memberAddress,
		Role:           role,
		IsAdmin:        sql.NullBool{Bool: isAdmin, Valid: true},
		Salary: sql.NullString{
			String: fmt.Sprintf("%d", salary),
			Valid:  true,
		},
	})

	if err != nil {
		s.Log.Err(err).Msg("Could not create organization member")
		return false, nil, proto.WrapError(proto.ErrInternal, err, "Could not create organization member")
	}

	responseOrgMember := &proto.OrganizationMember{
		OrganizationID: organizationUuid.String(),
		MemberAddress:  prototyp.HashFromString(orgMember.MemberAddress),
		Role:           orgMember.Role,
		IsAdmin:        orgMember.IsAdmin.Bool,
		Salary:         uint64(salary),
		DateJoined:     &orgMember.DateJoined,
	}
	return true, responseOrgMember, nil
}

func (s *RPC) GetOrganizationMember(ctx context.Context, organizationID string, memberAddress string) (*proto.OrganizationMember, error) {
	organizationUuid, err := uuid.FromBytes([]byte(organizationID))
	if err != nil {
		s.Log.Err(err).Msg("Invalid UUID provided")
		return nil, proto.WrapError(proto.ErrInvalidArgument, err, "Invalid UUID provided")
	}
	dbMember, err := data.DB.GetOrganizationMember(ctx, sqlc.GetOrganizationMemberParams{
		OrganizationID: organizationUuid,
		MemberAddress:  memberAddress,
	})
	if err != nil {
		s.Log.Err(err).Msg("Could not get organization member")
		return nil, proto.WrapError(proto.ErrInternal, err, "Could not get organization member")
	}
	salaryINT, _ := strconv.Atoi(dbMember.Salary.String)

	responseOrgMember := &proto.OrganizationMember{
		OrganizationID: organizationUuid.String(),
		MemberAddress:  prototyp.HashFromString(dbMember.MemberAddress),
		Role:           dbMember.Role,
		IsAdmin:        dbMember.IsAdmin.Bool,
		Salary:         uint64(salaryINT),
		DateJoined:     &dbMember.DateJoined,
	}
	return responseOrgMember, nil
}

func (s *RPC) GetAllOrganizationMembers(ctx context.Context, organizationID string) ([]*proto.OrganizationMember, error) {
	organizationUuid, err := uuid.FromBytes([]byte(organizationID))
	if err != nil {
		s.Log.Err(err).Msg("Invalid UUID provided")
		return nil, proto.WrapError(proto.ErrInvalidArgument, err, "Invalid UUID provided")
	}
	dbMembers, err := data.DB.GetAllOrganizationMembers(ctx, organizationUuid)
	if err != nil {
		s.Log.Err(err).Msg("Could not get organization member")
		return nil, proto.WrapError(proto.ErrInternal, err, "Could not get organization member")
	}
	responseOrgMembers := make([]*proto.OrganizationMember, len(dbMembers))
	for i, dbMember := range dbMembers {
		salaryINT, _ := strconv.Atoi(dbMember.Salary.String)
		responseOrgMembers[i] = &proto.OrganizationMember{
			OrganizationID: organizationUuid.String(),
			MemberAddress:  prototyp.HashFromString(dbMember.MemberAddress),
			Role:           dbMember.Role,
			IsAdmin:        dbMember.IsAdmin.Bool,
			Salary:         uint64(salaryINT),
			DateJoined:     &dbMember.DateJoined,
		}
	}
	return responseOrgMembers, nil
}

func (s *RPC) UpdateOrganizationMember(ctx context.Context, organizationMember *proto.OrganizationMember) (bool, *proto.OrganizationMember, error) {
	organizationUuid, err := uuid.FromBytes([]byte(organizationMember.OrganizationID))
	if err != nil {
		s.Log.Err(err).Msg("Invalid UUID provided")
		return false, nil, proto.WrapError(proto.ErrInvalidArgument, err, "Invalid UUID provided")
	}

	isCurrentUserAdmin, err := s.checkOrgAdmin(ctx, organizationUuid)
	if err != nil {
		return false, nil, err
	}

	if !isCurrentUserAdmin {
		return false, nil, proto.WrapError(proto.ErrPermissionDenied, errors.New("not an admin"), "not an admin")
	}

	dbOrgMember, err := data.DB.UpdateOrganizationMember(ctx, sqlc.UpdateOrganizationMemberParams{
		OrganizationID: organizationUuid,
		MemberAddress:  organizationMember.MemberAddress.String(),
		Role:           organizationMember.Role,
		IsAdmin:        sql.NullBool{Bool: organizationMember.IsAdmin, Valid: true},
		Salary: sql.NullString{
			String: fmt.Sprintf("%d", organizationMember.Salary),
			Valid:  true,
		},
	})

	if err != nil {
		return false, nil, proto.WrapError(proto.ErrInternal, err, "Could not update organization member")
	}
	salaryINT, _ := strconv.Atoi(dbOrgMember.Salary.String)

	return true, &proto.OrganizationMember{
		OrganizationID: organizationUuid.String(),
		MemberAddress:  prototyp.HashFromString(dbOrgMember.MemberAddress),
		Role:           dbOrgMember.Role,
		IsAdmin:        dbOrgMember.IsAdmin.Bool,
		Salary:         uint64(salaryINT),
		DateJoined:     &dbOrgMember.DateJoined,
	}, nil
}

func (s *RPC) DeleteOrganizationMember(ctx context.Context, organizationID, memberAddress string) (bool, error) {
	organizationUuid, err := uuid.FromBytes([]byte(organizationID))
	if err != nil {
		s.Log.Err(err).Msg("Invalid UUID provided")
		return false, proto.WrapError(proto.ErrInvalidArgument, err, "Invalid UUID provided")
	}

	isCurrentUserAdmin, err := s.checkOrgAdmin(ctx, organizationUuid)
	if err != nil {
		return false, err
	}

	if !isCurrentUserAdmin {
		return false, proto.WrapError(proto.ErrPermissionDenied, errors.New("not an admin"), "not an admin")
	}

	err = data.DB.DeleteOrganizationMember(ctx, sqlc.DeleteOrganizationMemberParams{
		OrganizationID: organizationUuid,
		MemberAddress:  memberAddress,
	})

	if err != nil {
		return false, proto.WrapError(proto.ErrInternal, err, "Could not delete organization member")
	}

	return true, nil
}
