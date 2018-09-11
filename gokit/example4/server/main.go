package main

import (
	"context"
	"fmt"
	"go-kit__VS__go-micro/gokit/example4/pb"
	test "go-kit__VS__go-micro/gokit/example4/service"
	"net"
	"github.com/go-kit/kit/endpoint"
	gokitgrpc "github.com/go-kit/kit/transport/grpc"
	oldcontext "golang.org/x/net/context"
	"google.golang.org/grpc"
)

type service struct{}

func (service) Test(ctx context.Context, a string, b int64) (context.Context, string, error) {
	return nil, fmt.Sprintf("hello %s you are %d old?", a, b), nil
}

func NewService() test.Service {
	return service{}
}

func makeTestEndpoint(svc test.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(test.TestRequest)
		newCtx, v, err := svc.Test(ctx, req.A, req.B)
		return &test.TestResponse{
			V:   v,
			Ctx: newCtx,
		}, err
	}
}

type Servers struct {
	server gokitgrpc.Handler
}

func (b *Servers) Test(ctx oldcontext.Context, req *pb.TestRequest) (*pb.TestResponse, error) {
	_, response, err := b.server.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return response.(*pb.TestResponse), nil
}

func NewBinding(svc test.Service) *Servers {
	return &Servers{
		server: gokitgrpc.NewServer(
			makeTestEndpoint(svc),
			test.DecodeRequest,
			test.EncodeResponse,
		),
	}
}

const (
	hostPort string = "localhost:8080"
)

func main() {

	sc, err := net.Listen("tcp", hostPort)
	if err != nil {
		fmt.Println("client create issue", err)
	}

	service := NewService()
	server := grpc.NewServer()
	defer server.GracefulStop()
	pb.RegisterTestServer(server, NewBinding(service))

	server.Serve(sc)
}
