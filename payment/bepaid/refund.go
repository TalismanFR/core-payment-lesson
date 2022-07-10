package bepaid

import (
	"context"
	"diLesson/application/domain"
	"diLesson/payment/contract/dto"
	"github.com/dantedenis/bepaid/api/contracts"
	"github.com/dantedenis/bepaid/service"
	"github.com/dantedenis/bepaid/service/vo"
	"github.com/satori/go.uuid"
)

type Refund struct {
	api        contracts.Api
	card       vo.CreditCard
	result     string
	parentUUID uuid.UUID
}

func NewRefund(api contracts.Api, card vo.CreditCard, uuid uuid.UUID) *Refund {
	return &Refund{
		api:        api,
		parentUUID: uuid,
		card:       card,
	}
}

func (r Refund) Refund(ctx context.Context, pay *domain.Pay, reason string) (*dto.VendorRefundResult, error) {
	serv := ApiServiceRefund{service.NewApiService(r.api)}
	tokenAuth := vo.NewAuthorizationRequest(int64(pay.Amount()), pay.Currency(), "", pay.Transaction(), false, r.card)

	authorizationRequest, err := serv.Authorizations(context.Background(), *tokenAuth)
	if err != nil {
		return nil, err
	}

	if !authorizationRequest.IsAuthorization() {
		return dto.NewVendorRefundResult(
			pay.Terminal().Alias(),
			authorizationRequest.Response.Message,
			reason,
			true,
		), nil
	}
	resultRefund, err := serv.Refund(ctx, vo.NewRefundRequest(r.parentUUID.String(), int64(pay.Amount()), reason))
	if err != nil {
		return nil, err
	}

	return dto.NewVendorRefundResult(pay.Terminal().Alias(), resultRefund.Response.Message, reason, resultRefund.IsFailed()), nil
}

// ApiServiceRefund Так как нет реализации в bepaid, создаем композицию и добавляем нужный метод
type ApiServiceRefund struct {
	*service.ApiService
}

func (a *ApiServiceRefund) Refund(ctx context.Context, request *vo.RefundRequest) (vo.TransactionResponse, error) {
	return vo.TransactionResponse{}, nil
}
