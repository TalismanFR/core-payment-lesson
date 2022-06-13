package bepaid

import (
	"context"
	"diLesson/application/domain"
	"diLesson/payment/contract/dto"
	"github.com/dantedenis/bepaid/api/contracts"
	"github.com/dantedenis/bepaid/service"
	"github.com/dantedenis/bepaid/service/vo"
)

type Refund struct {
	api contracts.Api
}

func NewRefund(pay domain.Pay, reason string)  {
	
}

func (c Charge) Refund(pay *domain.Pay) (*dto.VendorRefundResult, error) {
	a := service.NewApiService(c.api)
	tokenAuth := vo.NewAuthorizationRequest(int64(pay.Amount()), pay.Currency().String(), "", pay.TransactionId(), false, c.Card)

	authorizationRequest, err := a.Authorizations(context.Background(), *tokenAuth)
	if err != nil {
		return nil, err
	}
	c.api.Refund()

	return dto.NewVendorRefundResult(pay.Terminal().Alias()), nil
}

func ()