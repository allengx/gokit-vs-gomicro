package main

import (
	"context"
	"fmt"
	"go-kit__VS__go-micro/gokit/example4/pb"
	test "go-kit__VS__go-micro/gokit/example4/service"
	"net"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd/etcdv3"
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

func main() {

	//----------------------------------------------------------

	var (
		//etcd 服务地址
		etcdServer = "127.0.0.1:2379"
		// 服务的信息目录
		prefix = "/services/book/"
		// 当前启动服务实例的地址
		instance = "127.0.0.1:50053"
		// 服务实例注册的路径
		key = prefix + instance
		// 服务实例注册的 val
		value = instance
		ctx   = context.Background()
		// 服务监听地址
		serviceAddress = ":50053"
	)

	//etcd 的连接参数
	options := etcdv3.ClientOptions{
		DialTimeout:   time.Second * 3,
		DialKeepAlive: time.Second * 3,
	}

	// 创建 etcd 连接
	client, err := etcdv3.NewClient(ctx, []string{etcdServer}, options)
	if err != nil {
		panic(err)
	}

	// 创建注册器
	registrar := etcdv3.NewRegistrar(client, etcdv3.Service{
		Key:   key,
		Value: value,
	}, log.NewNopLogger())

	// 注册器启动注册
	registrar.Register()

	//----------------------------------------------------------

	sc, err := net.Listen("tcp", serviceAddress)
	if err != nil {
		fmt.Println("client create issue", err)
	}

	service := NewService()
	server := grpc.NewServer()
	defer server.GracefulStop()
	pb.RegisterTestServer(server, NewBinding(service))

	server.Serve(sc)
}
