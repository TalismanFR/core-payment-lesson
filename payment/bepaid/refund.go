package bepaid

import (
	"context"
	"diLesson/application/domain"
	"diLesson/payment/contract/dto"
	"github.com/dantedenis/bepaid/api/contracts"
	"github.com/dantedenis/bepaid/service"
	"github.com/dantedenis/bepaid/service/vo"
	uuid "github.com/satori/go.uuid"
)

type Refund struct {
	api        contracts.Api
	Card       vo.CreditCard
	result     string
	reason     string
	parentUUID uuid.UUID
}

func NewRefund(uuid uuid.UUID, reason string) *Refund {
	return &Refund{
		reason:     reason,
		parentUUID: uuid,
	}
}

func (r Refund) Refund(pay *domain.Pay) (*dto.VendorRefundResult, error) {
	serv := service.NewApiService(r.api)
	tokenAuth := vo.NewAuthorizationRequest(int64(pay.Amount()), pay.Currency().String(), "", pay.TransactionId(), false, r.Card)

	authorizationRequest, err := serv.Authorizations(context.Background(), *tokenAuth)
	if err != nil {
		return nil, err
	}

	if !authorizationRequest.IsAuthorization() {
		return dto.NewVendorRefundResult(
			pay.Terminal().Alias(),
			authorizationRequest.Response.Message,
			r.reason,
			true,
		), nil
	}
	resultRefund, err := serv.Refund(context.Background(), vo.NewRefundRequest(r.parentUUID.String(), int64(pay.Amount()), r.reason))
	if err != nil {
		return nil, err
	}

	return dto.NewVendorRefundResult(pay.Terminal().Alias(), resultRefund.Response.Message, r.reason, resultRefund.IsFailed()), nil
}
