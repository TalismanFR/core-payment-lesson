package main

import (
	"auth/internal/api/v1"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func main() {

	con, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return
	}

	defer func(con *grpc.ClientConn) {
		err := con.Close()
		if err != nil {
			log.Println(err)
		}
	}(con)

	cl := v1.NewIdentityProviderClient(con)

	ctx := context.TODO()

	res, err := cl.SignUp(ctx, &v1.SignUpRequest{Email: "tim@ya.ru", Password: "qwerty12345", Role: "user"})
	if err != nil {
		log.Fatal(err)
	}

	log.Println(res)

	res2, err := cl.Login(ctx, &v1.LoginRequest{Email: "tim@ya.ru", Password: "qwerty12345"})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(res2)
}
