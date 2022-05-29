package rpc

import (
	"context"
	"errors"
	"strconv"

	"github.com/0xsequence/go-sequence/lib/prototyp"
	"github.com/crypto-pe/backend/data"
	"github.com/crypto-pe/backend/data/sqlc"
	"github.com/crypto-pe/backend/proto"
	"github.com/google/uuid"
)

func (s *RPC) CreateSalaryPayments(ctx context.Context, organizationID string, memberAddressesAmountMap map[string]uint64, transactionHash string, tokenType *proto.TokenType) (bool, []*proto.Payment, error) {

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

	var protoPayments []*proto.Payment

	for address, amount := range memberAddressesAmountMap {

		payment, err := data.DB.CreateSalaryPayments(ctx, sqlc.CreateSalaryPaymentsParams{
			OrganizationID:  organizationUuid,
			MemberAddress:   address,
			TransactionHash: transactionHash,
			Amount:          strconv.Itoa(int(amount)),
			Token:           tokenAddr,
		})

		if err != nil {
			s.Log.Err(err).Msg("Could not create salary payment")
			return false, nil, proto.WrapError(proto.ErrInternal, err, "Could not create salary payment")
		}

		protoPayments = append(protoPayments, &proto.Payment{
			PaymentID:       payment.PaymentID.String(),
			OrganizationID:  organizationID,
			MemberAddress:   prototyp.Hash(payment.MemberAddress),
			TransactionHash: payment.TransactionHash,
			Amount:          amount,
			Token:           tokenType,
			Date:            &payment.Date,
		})
	}

	return true, protoPayments, nil

}

func (s *RPC) GetOrgMemberSalaryPaymentsHistory(ctx context.Context, organizationID string, memberAddress string) ([]*proto.Payment, error) {
	organizationUuid, err := uuid.FromBytes([]byte(organizationID))
	if err != nil {
		s.Log.Err(err).Msg("Invalid UUID provided")
		return nil, proto.WrapError(proto.ErrInvalidArgument, err, "Invalid UUID provided")
	}

	payments, err := data.DB.GetOrgMemberSalaryPaymentsHistory(ctx, sqlc.GetOrgMemberSalaryPaymentsHistoryParams{
		MemberAddress:  memberAddress,
		OrganizationID: organizationUuid,
	})

	var protoPayemnts []*proto.Payment

	for _, payment := range payments {

		token, err := GetTokeTypeFromAddress(string(payment.Token))
		if err != nil {
			return nil, proto.WrapError(proto.ErrInternal, err, "Invalid token type")
		}

		paymentAmount, _ := strconv.Atoi(payment.Amount)
		protoPayemnts = append(protoPayemnts, &proto.Payment{
			PaymentID:       payment.PaymentID.String(),
			OrganizationID:  payment.OrganizationID.String(),
			MemberAddress:   prototyp.Hash(payment.MemberAddress),
			TransactionHash: payment.TransactionHash,
			Amount:          uint64(paymentAmount),
			Token:           &token,
			Date:            &payment.Date,
		})
	}

	return protoPayemnts, nil

}

func (s *RPC) GetMemberOverallSalaryHistory(ctx context.Context, memberAddress string) ([]*proto.Payment, error) {
	payments, err := data.DB.GetMemberOverallSalaryHistory(ctx, memberAddress)
	if err != nil {
		s.Log.Err(err).Msg("Could not find any member salary payment")
		return nil, proto.WrapError(proto.ErrNotFound, err, "Could not find any member salary payment")
	}

	var protoPayemnts []*proto.Payment

	for _, payment := range payments {

		token, err := GetTokeTypeFromAddress(string(payment.Token))
		if err != nil {
			return nil, proto.WrapError(proto.ErrInternal, err, "Invalid token type")
		}

		paymentAmount, _ := strconv.Atoi(payment.Amount)
		protoPayemnts = append(protoPayemnts, &proto.Payment{
			PaymentID:       payment.PaymentID.String(),
			OrganizationID:  payment.OrganizationID.String(),
			MemberAddress:   prototyp.Hash(payment.MemberAddress),
			TransactionHash: payment.TransactionHash,
			Amount:          uint64(paymentAmount),
			Token:           &token,
			Date:            &payment.Date,
		})
	}

	return protoPayemnts, nil
}
