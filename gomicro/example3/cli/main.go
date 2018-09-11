package main

import (
	"context"
	"log"
	pb "go-kit__VS__go-micro/gomicro/example3/pb"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	_ "github.com/micro/go-plugins/registry/etcd"
)

// log wrapper logs every time a request is made
type logWrapper struct {
	client.Client
}

func (l *logWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	log.Printf("req.serv: %s method: %s\n", req.Service(), req.Method())
	return l.Client.Call(ctx, req, rsp)
}

// Implements client.Wrapper as logWrapper
func logWrap(c client.Client) client.Client {
	return &logWrapper{c}
}

func main() {
	service := micro.NewService(
		micro.Name("Wrappers:"),
		// wrap the client
		micro.WrapClient(logWrap),
	)
	service.Init()
	wrappers := pb.NewWrappersService("wrap", service.Client())
	rsp, err := wrappers.Test(context.TODO(), &pb.Request{Name: "John"})
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(rsp.Msg)
}
