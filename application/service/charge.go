package service

import (
	"context"
	"diLesson/application"
	"diLesson/application/contract/dto"
	"diLesson/application/domain"
	"diLesson/application/domain/currency"
	"diLesson/payment/contract"
	"fmt"
	"github.com/golobby/container/v3"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	chargeSuffix = "_charge"
)

const (
	instrumentation        = "application.service.charge"
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
	payRepo      application.PayRepository `container:"type"`
	terminalRepo application.TerminalRepo  `container:"type"`
}

func (c Charge) Charge(ctx context.Context, request dto.ChargeRequest) (*dto.ChargeResult, error) {

	ctx, span := tracer.Start(ctx, "Charge")
	defer span.End()

	if err := request.Valid(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	cur, err := currency.FromString(request.Currency)
	if err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	uuidTerminal, err := uuid.Parse(request.TerminalId)
	if err != nil {
		return nil, err
	}

	terminal, err := c.terminalRepo.FindByUuid(ctx, uuidTerminal)
	if err != nil {
		return nil, err
	}

	pay, err := domain.NewPay(uuid.New(), domain.Amount(request.Amount), cur, request.Description, request.InvoiceId, terminal, &request.CreditCard)
	if err != nil {
		return nil, err
	}

	var vendor contract.VendorCharge
	err = container.NamedResolve(&vendor, terminal.Alias()+chargeSuffix)
	if err != nil {
		return nil, err
	}

	result, err := vendor.Charge(ctx, pay)
	if err != nil {
		return nil, err
	}

	pay.HandleChargeResult(result)

	err = c.payRepo.Save(ctx, pay)
	if err != nil {
		return nil, err
	}

	threeDS := &dto.ThreeDs{Status: dto.UnknownThreeDsStatus, RedirectUrl: result.ThreeDs().RedirectUrl}

	r := dto.NewChargeResult(pay.Status().Code(), pay.Status().Description(), pay.Uuid().String(), result.ReceiptUrl(), result.Message(), threeDS)
	return r, nil
}
