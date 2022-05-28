package rpc

import (
	"context"

	"github.com/crypto-pe/backend/proto"
)

func (s *RPC) CreateOrganizationMember(ctx context.Context, organization_uuid string, member_address string, role string, isAdmin bool, salary int) (bool, *proto.OrganizationMember, error) {

}

func (s *RPC) GetOrganizationMember(ctx context.Context, organization_uuid string, member_address string) (*proto.OrganizationMember, error) {

}

func (s *RPC) GetAllOrganizationMembers(ctx context.Context, organization_uuid string) ([]*proto.OrganizationMember, error) {

}

func (s *RPC) UpdateOrganizationMember(ctx context.Context, organizationMember *proto.OrganizationMember) (bool, *proto.OrganizationMember, error) {

}

func (s *RPC) DeleteOrganizationMember(ctx context.Context, organization_uuid, member_address string) (bool, error) {

}
