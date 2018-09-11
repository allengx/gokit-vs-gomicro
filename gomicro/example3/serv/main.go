package main

import (
	"context"
	"go-kit__VS__go-micro/gomicro/example3/pb"
	"log"
"fmt"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/server"
	_ "github.com/micro/go-plugins/registry/etcd"
	"github.com/micro/go-micro/transport"
)

type TestWrappers struct{}

func (tw *TestWrappers) Test(ctx context.Context, req *pb.Request, rep *pb.Response) error {
	log.Println("This is a test")
	log.Println("client:my name is ", req.Name)
	rep.Msg=fmt.Sprintf("\nserver: welcome %s\nserver: %s Bye~\n" ,req.Name,req.Name)
	return nil
}

func logWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		log.Printf("request method: %v\n", req.Method())
		err := fn(ctx, req, rsp)
		return err
	}
}


func main() {
	service := micro.NewService(
		micro.Name("wrap"),
		micro.WrapHandler(logWrapper),
		micro.Transport(transport.NewTransport()),
	)
	// optionally setup command line usage
	service.Init()
	// Register Handlers
	pb.RegisterWrappersHandler(service.Server(), new(TestWrappers))
	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
