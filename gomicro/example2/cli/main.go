package main

import (
	"fmt"
	"time"

	"context"
	pb "gomicro/example2/pb"

	"github.com/micro/go-log"
	"github.com/micro/go-micro"
)

// send events using the publisher
func sendMsg(cont string, p micro.Publisher) {
	t := time.NewTicker(time.Second)
	var i int64=0
	for _ = range t.C {
		ev := &pb.Msg{
			Id:      i,
			Context: fmt.Sprintf("context:%s", cont),
		}
		log.Logf("publishing message: %+v\n", ev)
		// publish an event
		if err := p.Publish(context.Background(), ev); err != nil {
			log.Logf("error publishing %v", err)
		}
	}
}

func main() {
	// create a service
	service := micro.NewService(
		micro.Name("pubsub"),
	)
	// parse command line
	service.Init()
	//as sanme as the RegisterSubscriber
	pub2 := micro.NewPublisher("publish--1--:", service.Client())
	go sendMsg("address--1--", pub2)
	time.Sleep(time.Duration(2) * time.Second)
}
