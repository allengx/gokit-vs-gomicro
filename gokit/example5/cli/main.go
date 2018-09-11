package main

import (
	"context"
	"fmt"
	"go-kit__VS__go-micro/gokit/example4/pb"
	test "go-kit__VS__go-micro/gokit/example4/service"
	"io"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/etcdv3"
	"github.com/go-kit/kit/sd/lb"
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

// 通过传入的 实例地址  创建对应的请求 endPoint
func reqFactory(instanceAddr string) (endpoint.Endpoint, io.Closer, error) {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		fmt.Println("请求服务:", instanceAddr)
		cc, err := grpc.Dial(instanceAddr, grpc.WithInsecure())
		if err != nil {
			fmt.Println("client create issue", err)
		}
		defer cc.Close()
		clientEndpoint := NewClient(cc)
		responseCTX, v, err := clientEndpoint.Test(context.Background(), "Allen", 21)
		if err != nil {
			fmt.Println("Test func issue", err)
		}
		fmt.Println("response is :", v)
		fmt.Println("context is :", responseCTX)
		//conn, err := grpc.Dial(instanceAddr, grpc.WithInsecure())
		// if err != nil {
		// 	fmt.Println(err)
		// 	panic("connect error")
		// }
		//defer conn.Close()
		// bookClient := book.NewBookServiceClient(conn)
		// bi,_:=bookClient.GetBookInfo(context.Background(),&book.BookInfoParams{BookId:1})
		// fmt.Println("获取书籍详情")
		// fmt.Println("bookId: 1", "=>", "bookName:", bi.BookName)

		// bl,_ := bookClient.GetBookList(context.Background(), &book.BookListParams{Page:1, Limit:10})
		// fmt.Println("获取书籍列表")
		// for _,b := range bl.BookList {
		// 	fmt.Println("bookId:", b.BookId, "=>", "bookName:", b.BookName)
		// }
		return nil, nil
	}, nil, nil
}

func main() {

	//--------------------------------

	var (
		// 注册中心地址
		etcdServer = "10.204.29.73:2379"
		// 监听的服务前缀
		prefix = "/services/book/"
		ctx    = context.Background()
	)
	options := etcdv3.ClientOptions{
		DialTimeout:   time.Second * 3,
		DialKeepAlive: time.Second * 3,
	}
	// 连接注册中心
	client, err := etcdv3.NewClient(ctx, []string{etcdServer}, options)
	if err != nil {
		panic(err)
	}
	logger := log.NewNopLogger()
	// 创建实例管理器, 此管理器会 Watch 监听 etc 中 prefix 的目录变化更新缓存的服务实例数据
	instancer, err := etcdv3.NewInstancer(client, prefix, logger)
	if err != nil {
		panic(err)
	}
	// 创建端点管理器， 此管理器根据 Factory 和监听的到实例创建 endPoint 并订阅 instancer 的变化动态更新 Factory 创建的 endPoint
	endpointer := sd.NewEndpointer(instancer, reqFactory, logger)
	// 创建负载均衡器
	balancer := lb.NewRoundRobin(endpointer)

	/**
	我们可以通过负载均衡器直接获取请求的 endPoint，发起请求
	reqEndPoint,_ := balancer.Endpoint()
	*/

	/**
	也可以通过 retry 定义尝试次数进行请求
	*/
	reqEndPoint := lb.Retry(3, 3*time.Second, balancer)
	
	//--------------------------------

	fmt.Println("aa")
	req := struct{}{}
	for i := 1; i <= 8; i++ {
		if _, err = reqEndPoint(ctx, req); err != nil {
			panic(err)
		}
	}


}
