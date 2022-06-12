package bepaid

import (
	"context"
	"diLesson/application"
	"diLesson/application/domain"
	"diLesson/payment/contract/dto"
	"fmt"
	"github.com/TalismanFR/bepaid/api"
	"github.com/TalismanFR/bepaid/service"
	sdkvo "github.com/TalismanFR/bepaid/service/vo"
	"github.com/golobby/container/v3"
	"net/http"
)

type Charge struct {
	vendor string
}

func NewCharge(vendor string) *Charge {
	return &Charge{vendor: vendor}
}

func authorizationRequestFromPay(pay *domain.Pay) *sdkvo.AuthorizationRequest {
	ccPay := pay.CreditCard()
	cc := sdkvo.NewCreditCard(ccPay.Number(), ccPay.VerificationValue(), ccPay.Holder(), sdkvo.ExpMonth(ccPay.ExpMonth().String()), ccPay.ExpYear())
	request := sdkvo.NewAuthorizationRequest(sdkvo.Amount(pay.Amount()), sdkvo.Currency(pay.Currency().String()), pay.Description(), pay.InvoiceId(), false, *cc)

	return request
}

func (c Charge) Charge(pay *domain.Pay) (*dto.VendorChargeResult, error) {

	// TODO: remove additionalParams, add Url field
	url, ok := pay.Terminal().AdditionalParams()["url"]
	if !ok {
		return nil, fmt.Errorf("terminal doesn't contain url")
	}

	if url == "" {
		return nil, fmt.Errorf("terminal url is empty")
	}

	var secrets application.SecretsService
	err := container.Resolve(&secrets)
	if err != nil {
		return nil, err
	}

	pair, err := secrets.Get(context.Background(), pay.Terminal().Uuid())
	if err != nil {
		return nil, fmt.Errorf("cannot extract shop credentials: %w", err)
	}

	client := service.NewApiService(api.NewApi(http.DefaultClient, url, pair.ShopId, pair.Secret))

	ar := authorizationRequestFromPay(pay)
	resp, err := client.Authorizations(context.Background(), *ar)
	if err != nil {
		return nil, err
	}

	uid := resp.Uid()
	if resp.IsError() {
		return nil, fmt.Errorf(resp.Response.Message)
	}

	cr := sdkvo.NewCaptureRequest(uid, sdkvo.Amount(pay.Amount()))
	resp, err = client.Capture(context.Background(), *cr)
	if err != nil {
		return nil, err
	}

	uid = resp.Uid()

	vendorStatus := dto.UnknownVendorChargeStatus

	if resp.IsFailed() {
		vendorStatus = dto.FailedVendorChargeStatus
	}
	if resp.IsSuccess() {
		vendorStatus = dto.SuccessfulVendorChargeStatus
	}
	if resp.IsIncomplete() {
		vendorStatus = dto.Need3DSVendorChargeStatus
	}

	vendor3ds := &dto.VendorThreeDs{Status: dto.UnknownThreeDsVendorStatus, RedirectUrl: "example.com/redirect"}

	return dto.NewVendorChargeResult(c.vendor, uid, resp.Transaction.Message, vendorStatus, resp.Transaction.ReceiptUrl, vendor3ds), nil
}
