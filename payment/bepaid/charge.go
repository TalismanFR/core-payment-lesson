package bepaid

import (
	"context"
	"diLesson/application/domain"
	"diLesson/payment/contract/dto"
	"github.com/TalismanFR/bepaid/api"
	"github.com/TalismanFR/bepaid/service"
	sdkvo "github.com/TalismanFR/bepaid/service/vo"
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

	client := service.NewApiService(api.NewApi(http.DefaultClient, "", "", ""))

	ar := authorizationRequestFromPay(pay)
	resp, err := client.Authorizations(context.Background(), *ar)
	if err != nil {
		return nil, err
	}

	uid := resp.Uid()

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
