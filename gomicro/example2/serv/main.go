package main

import (
	pb "gomicro/example2/pb"

	"github.com/micro/go-log"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/metadata"

	"context"
	_ "github.com/micro/go-plugins/broker/rabbitmq"
	_ "github.com/micro/go-plugins/registry/kubernetes"
	_ "github.com/micro/go-plugins/transport/nats"
)

// All methods of Sub will be executed when
// a message is received
type Sub struct{}

//Method can be of any name
func (s *Sub) Process(ctx context.Context, msg *pb.Msg) error {
	md, _ := metadata.FromContext(ctx)
	log.Logf("get %+v with %+v\n", msg, md)
	return nil
}

func main() {
	// create a service
	service := micro.NewService(
		micro.Name("pubsub"),
	)
	// parse command line
	service.Init()
	//register subscriber										and run sub full func
	micro.RegisterSubscriber("publish--1--:", service.Server(), new(Sub))
	// one publish msg and all subscriber can get
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
