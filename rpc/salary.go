package rpc

import (
	"context"
	"errors"
	"fmt"

	"github.com/0xsequence/go-sequence/lib/prototyp"
	"github.com/crypto-pe/backend/data"
	"github.com/crypto-pe/backend/data/sqlc"
	"github.com/crypto-pe/backend/proto"
	"github.com/google/uuid"
)

func (s *RPC) CreateSalaryPayment(ctx context.Context, organizationID string, memberAddress string,
	transactionHash string, amount uint64,
	tokenType *proto.TokenType) (bool, *proto.Payment, error) {

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
	tokenAddr, _ := GetTokenAddressFromTokenType(tokenType)
	payment, err := data.DB.CreateSalaryPayment(ctx, sqlc.CreateSalaryPaymentParams{
		OrganizationID:  organizationUuid,
		MemberAddress:   []byte(memberAddress),
		TransactionHash: transactionHash,
		Amount:          fmt.Sprintf("%d", amount),
		Token:           []byte(tokenAddr),
	})
	if err != nil {
		s.Log.Err(err).Msg("Could not create salary payment")
		return false, nil, proto.WrapError(proto.ErrInternal, err, "Could not create salary payment")
	}

	return false, &proto.Payment{
		OrganizationID:  organizationUuid.String(),
		MemberAddress:   prototyp.HashFromBytes(payment.MemberAddress),
		TransactionHash: payment.TransactionHash,
		Amount:          payment.Amount,
		Token:           tokenType,
		Date:            &payment.Date,
	}, nil
}
