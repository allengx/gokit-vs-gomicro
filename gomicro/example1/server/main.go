package main

import (
	"log"
	"time"
	test "gomicro/example1/pb"
	"context"
	"github.com/micro/go-micro"
)

type TestHello struct{}

func (s *TestHello) Hello(ctx context.Context, req *test.Request, rsp *test.Response) error {
	log.Print("Received Say.Hello request")
	rsp.Msg = "Hello " + req.Name
	return nil
}

func main() {
	service := micro.NewService(
		micro.Name("test"),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
	)
	// optionally setup command line usage
	service.Init()
	// Register Handlers
	test.RegisterTestHandler(service.Server(), new(TestHello))
	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
