package server

import (
	"context"
	"diLesson/application/contract"
	"diLesson/application/contract/dto"
	"diLesson/application/domain/credit_card"
	"go.opentelemetry.io/otel"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	instrumentation        = "server"
	instrumentationVersion = "v0.0.1"
)

var (
	tracer = otel.Tracer(
		instrumentation,
		trace.WithSchemaURL(semconv.SchemaURL),
		trace.WithInstrumentationVersion(instrumentationVersion),
	)
)

type Server struct {
	UnimplementedPayServiceServer
	chargeService contract.Charge
}

func (s Server) Charge(ctx context.Context, message *ChargeRequestMessage) (*ChargeResponseMessage, error) {

	// TODO: replace with ctx
	ctx, span := tracer.Start(context.Background(), "charge request")
	defer span.End()

	messageCC := message.GetCreditCard()

	expMonth, _ := credit_card.FromInt(int(messageCC.GetExpMonth().Number() + 1))

	cc := credit_card.NewCreditCard(
		messageCC.GetNumber(),
		messageCC.GetVerificationValue(),
		messageCC.GetHolder(),
		expMonth, messageCC.ExpYear,
	).WithSkip3DSVerification(messageCC.SkipThreeDSecureVerification)

	chargeRequest := *dto.NewChargeRequest(message.GetAmount(), message.GetCurrency(), message.GetTerminalId(), message.GetInvoiceId(), message.GetDescription(), *cc)

	chargeResult, err := s.chargeService.Charge(ctx, chargeRequest)

	if err != nil {
		return nil, err
	}

	chargeResult.ThreeDS()

	chargeResponseMsg := &ChargeResponseMessage{
		StatusCode: int32(chargeResult.Status()),
		StatusName: chargeResult.StatusName(),
		Uuid:       chargeResult.Uuid(),
		ReceiptUrl: chargeResult.ReceiptUrl(),
		Message:    chargeResult.Message(),
		ThreeDs: &ChargeResponseMessage_ThreeDs{
			Status:      string(chargeResult.ThreeDS().Status),
			RedirectUrl: chargeResult.ThreeDS().RedirectUrl,
		},
	}

	return chargeResponseMsg, nil
}

func NewServer(chargeService contract.Charge) *Server {
	return &Server{chargeService: chargeService}
}
