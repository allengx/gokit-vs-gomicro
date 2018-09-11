
---


<center><h1>gokit通信基于grpc<br></h1></center><br>

> - 制定通信协议pb文件
> - 根据通信协议编写服务代码 service
> - 编写服务器代码 server
> - 编写客户端代码 client

---
#### 一.制定通信协议pb文件
##### 1.test.proto:


```
syntax = "proto3";

package pb;

service Test {
  rpc Test (TestRequest) returns (TestResponse) {}
}

message TestRequest {
  string a = 1;
  int64 b = 2;
}

message TestResponse {
  string v = 1;
}
```

##### 2.使用命令行生成 .pb.go 文件


> protoc XXX.proto --go_out=plugins=grpc:.


##### 3.现在可以看见一个 .pb.go文件，我们习惯上会把它和 .proto 文件打包在 pb文件包内，用于其他程序调用。


<html>
<center><h4>第一部分结束</h4></center>
</html>

---
#### 二.根据通信协议编写服务代码 service

##### 1.定义service.go（注：变量和方法名要和协议内容保持一致）


```golang
package service

import "context"

type Service interface {
	Test(ctx context.Context, a string, b int64) (context.Context, string, error)
}

type TestRequest struct {
	A string
	B int64
}

type TestResponse struct {
	Ctx context.Context
	V string
}

```

##### 2.定义相应的解码编码方法     response_replay.go
###### 理论上每个服务方法必须编写四个编码方法
> - 解码request（DecodeRequest）
> - 编码request（EncodeRequest）
> - 解码response（DecodeResponse）
> - 编码response（EncodeResponse）


```golang
package service

import (
	"context"

	"go-kit__VS__go-micro/gokit/example4/pb"
)

func EncodeRequest(ctx context.Context, req interface{}) (interface{}, error) {
	r := req.(TestRequest)
	return &pb.TestRequest{A: r.A, B: r.B}, nil
}

func DecodeRequest(ctx context.Context, req interface{}) (interface{}, error) {
	r := req.(*pb.TestRequest)
	return TestRequest{A: r.A, B: r.B}, nil
}

func EncodeResponse(ctx context.Context, resp interface{}) (interface{}, error) {
	r := resp.(*TestResponse)
	return &pb.TestResponse{V: r.V}, nil
}

func DecodeResponse(ctx context.Context, resp interface{}) (interface{}, error) {
	r := resp.(*pb.TestResponse)
	return &TestResponse{V: r.V, Ctx: ctx}, nil
}
```


<html>
<center><h4>第二部分结束</h4></center>
</html>


---

#### 三.编写服务器代码 server

##### 1.定义执行对象和实现接口方法


```golang
type service struct{}

func (service) Test(ctx context.Context, a string, b int64) (context.Context, string, error) {
	return nil, fmt.Sprintf("hello %s you are %d old?", a, b), nil
}
```


##### 2.编写工厂生成service


```golang
import (
	test "go-kit__VS__go-micro/gokit/example4/service"
)
```


```golang
func NewService() test.Service {
	return service{}
}
```

##### 3.生成endpoint对象（也是服务的实现入口）


```golang
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
```

##### 4.定义server和serverhandler（基于grpc通信）


```golang
import (
	gokitgrpc "github.com/go-kit/kit/transport/grpc"
	oldcontext "golang.org/x/net/context"
	"google.golang.org/grpc"
)
```


```golang
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
```

##### 5.构造server工厂


```golang
func NewBinding(svc test.Service) *Servers {
	return &Servers{
		server: gokitgrpc.NewServer(
			makeTestEndpoint(svc),
			test.DecodeRequest,
			test.EncodeResponse,
			//其中还可以扩充内容
			//options ...ServerOption,
			//1.  ServerBefore
			//2.  ServerAfter
			//3.  ServerErrorLogger
			//4.  ServerFinalizer
			
			/*
			官方示例
			grpctransport.ServerBefore(
				extractCorrelationID,
			),
			grpctransport.ServerBefore(
				displayServerRequestHeaders,
			),
			grpctransport.ServerAfter(
				injectResponseHeader,
				injectResponseTrailer,
				injectConsumedCorrelationID,
			),
			grpctransport.ServerAfter(
				displayServerResponseHeaders,
				displayServerResponseTrailers,
			),
			*/
		),
	}
}
```


##### 6.main函数


```golang
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
	
	//注册
	pb.RegisterTestServer(server, NewBinding(service))
	//启动
	server.Serve(sc)
}
```


##### 7.完整代码 server


```golang
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

```





<html>
<center><h4>第三部分结束</h4></center>
</html>


---

#### 四.编写客户端代码 client

##### 1.实现 pb内的 client API


```golang
type Clients struct{
	client endpoint.Endpoint
}

func (c *Clients) Test(ctx context.Context, a string, b int64) (context.Context, string, error) {
	response, err := c.client(ctx, test.TestRequest{A: a, B: b})
	if err != nil {
		return nil, "", err
	}
	fmt.Println(a)
	fmt.Println("年龄:",b)
	r := response.(*test.TestResponse)
	return r.Ctx, r.V, nil
}
```

##### 2.client 工厂


```golang
func NewClient(cc *grpc.ClientConn) test.Service {
	return &Clients{
		client: gokitgrpc.NewClient(
			cc,
			"pb.Test1",
			"Test1",
			test.EncodeRequest,
			test.DecodeResponse,
			&pb.TestResponse{},
			
	        //内容可添加
		    //1.ClientOption
		    //2.ClientBefore
		    //3.ClientAfter
		    //4.ClientFinalizer
		    //5.ClientFinalizerFunc

            // 官方示例
            // grpctransport.ClientBefore(
            // 	injectCorrelationID,
            // ),
            // grpctransport.ClientBefore(
            // 	displayClientRequestHeaders,
            // ),
            // grpctransport.ClientAfter(
            // 	displayClientResponseHeaders,
            // 	displayClientResponseTrailers,
            // ),
            // grpctransport.ClientAfter(
            // 	extractConsumedCorrelationID,
            // ),   
		).Endpoint(),
	}
}
```

##### 3. main 函数


```golang
const (
	hostPort string = "localhost:8080"
)

func main() {

	cc, err := grpc.Dial(hostPort, grpc.WithInsecure())
	if err != nil {
		fmt.Println("client create issue", err)
	}

    //得到 service
	clientEndpoint := NewClient(cc)

    //向 server 的 service 发送请求
	responseCTX, v, err := clientEndpoint.Test(context.Background(), "Allen", 21)
	if err != nil {
		fmt.Println("Test func issue", err)
	}
	fmt.Println("response is :", v)
	fmt.Println("context is :", responseCTX)
}
```


##### 4.完整代码 client


```golang
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

```


<html>
<center><h4>第四部分结束</h4></center>
</html>


---

##### 详细代码可参阅 gokit example4
##### URL[：http://10.204.28.137/allenguo/go-kit__VS__go-micro/tree/master/gokit/example4](http://10.204.28.137/allenguo/go-kit__VS__go-micro/tree/master/gokit/example4)

###### 编著人：Allen guo
###### 日期：  2018/8/29 