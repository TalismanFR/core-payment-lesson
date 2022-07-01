package bepaid

import (
	"context"
	"diLesson/application"
	"diLesson/application/domain"
	"diLesson/payment/contract/dto"
	"fmt"
	sdkapi "github.com/TalismanFR/bepaid/api"
	sdkservice "github.com/TalismanFR/bepaid/service"
	sdkvo "github.com/TalismanFR/bepaid/service/vo"
	"github.com/golobby/container/v3"
	"go.opentelemetry.io/otel"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.opentelemetry.io/otel/trace"
	"net/http"
)

const (
	instrumentation        = "payment.bepaid.charge"
	instrumentationVersion = "v0.0.1"
)

var (
	tracer = otel.Tracer(
		instrumentation,
		trace.WithSchemaURL(semconv.SchemaURL),
		trace.WithInstrumentationVersion(instrumentationVersion),
	)
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
	request := sdkvo.NewAuthorizationRequest(sdkvo.Amount(pay.Amount()), sdkvo.Currency(pay.Currency().String()), pay.Description(), pay.InvoiceId(), true, *cc)

	return request
}

func (c Charge) Charge(ctx context.Context, pay *domain.Pay) (*dto.VendorChargeResult, error) {

	ctx, span := tracer.Start(ctx, "Charge")
	defer span.End()

	// Get args for api service
	var terminals application.TerminalRepo
	err := container.Resolve(&terminals)
	if err != nil {
		return nil, err
	}

	//TODO: remove FindByUuid call
	terminal, err := terminals.FindByUuid(ctx, pay.Terminal().Uuid())
	if err != nil {
		return nil, fmt.Errorf("cannot extract shop credentials: %w", err)
	}

	shopId, secret, url, err := readTerminalSecrets(terminal.AdditionalParams())
	if err != nil {
		return nil, err
	}

	client := sdkservice.NewApiService(sdkapi.NewApi(http.DefaultClient, url, shopId, secret))

	// Create auth request
	ar := authorizationRequestFromPay(pay)
	_, spanAuth := tracer.Start(ctx, "Authorizations")
	resp, err := client.Authorizations(context.Background(), *ar)
	spanAuth.End()
	if err != nil {
		return nil, err
	}

	uid := resp.Uid()
	if resp.IsError() {
		return nil, fmt.Errorf(resp.Response.Message)
	}

	// Create capture request
	cr := sdkvo.NewCaptureRequest(uid, sdkvo.Amount(pay.Amount()))
	_, spanCapt := tracer.Start(ctx, "Capture")
	resp, err = client.Capture(context.Background(), *cr)
	spanCapt.End()
	if err != nil {
		return nil, err
	}

	uid = resp.Uid()

	// Create VendorChargeResult
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

func readTerminalSecrets(sm map[string]interface{}) (shopId string, secret string, url string, err error) {
	v1, ok := sm["shop_id"]
	if !ok {
		return "", "", "", fmt.Errorf("key shop_id is absent in a map")
	}

	shopId, ok = v1.(string)
	if !ok {
		return "", "", "", fmt.Errorf("shopId type isn't a string. type: %T", v1)
	}

	v2, ok := sm["secret"]
	if !ok {
		return "", "", "", fmt.Errorf("key 'secret' is absent in a map")
	}

	secret, ok = v2.(string)
	if !ok {
		return "", "", "", fmt.Errorf("secret type isn't a string. type: %T", v2)
	}

	v3, ok := sm["url"]
	if !ok {
		return "", "", "", fmt.Errorf("key 'url' is absent in a map")
	}

	url, ok = v3.(string)
	if !ok {
		return "", "", "", fmt.Errorf("url type isn't a string. type: %T", v2)
	}

	return
}
