package bepaid

import (
	"context"
	"diLesson/application/domain"
	"diLesson/payment/contract/dto"
	"github.com/dantedenis/bepaid/api/contracts"
	"github.com/dantedenis/bepaid/service"
	"github.com/dantedenis/bepaid/service/vo"
)

type Charge struct {
	api  contracts.Api
	Card vo.CreditCard
}

func (c Charge) Charge(pay *domain.Pay) (*dto.VendorChargeResult, error) {
	serv := service.NewApiService(c.api)
	// not descriptions
	tokenAuth := vo.NewAuthorizationRequest(int64(pay.Amount()), pay.Currency().String(), "", pay.TransactionId(), false, c.Card)
	authorizationRequest, err := serv.Authorizations(context.Background(), *tokenAuth)
	if err != nil {
		return nil, err
	}
	if !authorizationRequest.IsAuthorization() {
		return dto.NewVendorChargeResult(pay.Terminal().Alias(), authorizationRequest.IsFailed(), authorizationRequest.Response.Message), nil
	}
	resultPayment, err := serv.Capture(context.Background(), *transToCapture(&authorizationRequest))
	if err != nil {
		return nil, err
	}
	return dto.NewVendorChargeResult(pay.Terminal().Alias(), resultPayment.IsFailed(), resultPayment.Response.Message), nil
}

func transToCapture(response *vo.TransactionResponse) *vo.CaptureRequest {
	return vo.NewCaptureRequest(response.Transaction.Amount, response.Transaction.ParentUid)
}
