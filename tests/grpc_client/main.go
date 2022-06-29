package main

import (
	"context"
	"diLesson/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	cl := server.NewPayServiceClient(conn)

	cc := server.ChargeRequestMessage_CreditCard{
		Number:            "4200000000000000",
		VerificationValue: "123",
		Holder:            "tim",
		ExpMonth:          server.ChargeRequestMessage_CreditCard_JAN,
		ExpYear:           "2024",
	}
	req := server.ChargeRequestMessage{
		Amount:      1000,
		Currency:    "RUB",
		TerminalId:  "4a92cc1c-d381-4d16-92f6-55925490614a",
		InvoiceId:   "invoiceID1",
		Description: "my_description",
		CreditCard:  &cc,
	}

	resp, err := cl.Charge(context.Background(), &req)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(resp)
}
