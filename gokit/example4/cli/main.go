package main

import (
	"context"
	"fmt"
	"go-kit__VS__go-micro/gokit/example4/pb"
	test "go-kit__VS__go-micro/gokit/example4/service"

	"github.com/go-kit/kit/endpoint"
	gokitgrpc "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
)

type Clients struct {
	client endpoint.Endpoint
}

func (c *Clients) Test(ctx context.Context, a string, b int64) (context.Context, string, error) {
	response, err := c.client(ctx, test.TestRequest{A: a, B: b})
	if err != nil {
		return nil, "", err
	}
	fmt.Println(a)
	fmt.Println("年龄:", b)
	r := response.(*test.TestResponse)
	return r.Ctx, r.V, nil
}

func NewClient(cc *grpc.ClientConn) test.Service {
	return &Clients{
		client: gokitgrpc.NewClient(
			cc,
			"pb.Test",
			"Test",
			test.EncodeRequest,
			test.DecodeResponse,
			&pb.TestResponse{},
		).Endpoint(),
	}
}

const (
	hostPort string = "localhost:8080"
)

func main() {

	cc, err := grpc.Dial(hostPort, grpc.WithInsecure())
	if err != nil {
		fmt.Println("client create issue", err)
	}

	clientEndpoint := NewClient(cc)

	responseCTX, v, err := clientEndpoint.Test(context.Background(), "Allen", 21)
	if err != nil {
		fmt.Println("Test func issue", err)
	}
	fmt.Println("response is :", v)
	fmt.Println("context is :", responseCTX)
}
