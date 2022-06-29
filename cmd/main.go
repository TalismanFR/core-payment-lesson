package main

import (
	"diLesson/application/contract"
	"diLesson/config"
	"diLesson/server"
	"github.com/golobby/container/v3"
	"google.golang.org/grpc"
	"log"
	"net"
	"path/filepath"
)

func main() {

	log.Println("config: start")

	p, err := filepath.Abs("configs/main.yaml")
	if err != nil {
		log.Fatal(err)
	}

	conf, err := config.Parse(p)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("config: end")

	// TODO:sometimes test fail without time.Sleep
	//time.Sleep(5 * time.Second)

	log.Println("building dependencies: start")
	err = config.BuildDI(conf)
	if err != nil {
		log.Fatal(err)
	}

	var service contract.Charge
	err = container.Resolve(&service)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("building dependencies: end")

	s := server.NewServer(service)
	ls, err := net.Listen("tcp", ":50051")

	gs := grpc.NewServer()
	server.RegisterPayServiceServer(gs, s)

	log.Println("server: start")

	if err = gs.Serve(ls); err != nil {
		log.Fatal(err)
	}

	log.Println("server: end")
}
