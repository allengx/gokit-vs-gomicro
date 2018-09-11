
---

<center><h1>gokit实现etcd服务发现和负载均衡<br></h1></center><br>

> - 制定通信协议pb文件
> - 根据通信协议编写服务代码 service
> - 编写服务器代码 server
> - 编写客户端代码 client


---

### 1.制定通信协议pb文件
### 2.根据通信协议编写服务代码 service

#####  前两部分内容参考 gokit example4
#####  URL[:http://10.204.28.137/allenguo/go-kit__VS__go-micro/blob/master/gokit/example4/README.md](http://10.204.28.137/allenguo/go-kit__VS__go-micro/blob/master/gokit/example4/README.md)

<html>
<center><h4>第一、二部分结束</h4></center>
</html>


---

### 3.编写服务器代码 server

> ##### 1.定义执行对象和实现接口方法
> ##### 2.编写工厂生成service
> ##### 3.生成endpoint对象（也是服务的实现入口）
> ##### 4.定义server和serverhandler（基于grpc通信）
> ##### 5.构造server工厂

##### 前面内容与 gokit example4一致

##### 6.完成etcd服务注册（main 函数）


```golang
import(
    "github.com/go-kit/kit/sd/etcdv3"
)
```



```golang
func main() {

	//----------------------------------------------------------

	var (
		//etcd 服务地址
		etcdServer = "10.204.29.77:2379"
		// 服务的信息目录
		prefix = "/services/book/"
		// 当前启动服务实例的地址
		instance = "10.204.29.77:50051"
		// 服务实例注册的路径
		key = prefix + instance
		// 服务实例注册的 val
		value = instance
		ctx   = context.Background()
		// 服务监听地址
		serviceAddress = ":50051"
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
```

##### 7. 完整代码


```golang
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
		etcdServer = "10.204.29.77:2379"
		// 服务的信息目录
		prefix = "/services/book/"
		// 当前启动服务实例的地址
		instance = "10.204.29.77:50051"
		// 服务实例注册的路径
		key = prefix + instance
		// 服务实例注册的 val
		value = instance
		ctx   = context.Background()
		// 服务监听地址
		serviceAddress = ":50051"
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

```

<html>
<center><h4>第三部分结束</h4></center>
</html>


---

### 4.编写客户端代码 client

> ###### 1.实现 pb内的 client API
> ###### 2.client 工厂
##### 前面内容与 gokit example4一致

##### 3.请求工厂 reqFactory （把 gokit example4 main的代码打包了一下）


```golang
// 通过传入的 实例地址  创建对应的请求 endPoint
func reqFactory(instanceAddr string) (endpoint.Endpoint, io.Closer, error) {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		fmt.Println("请求服务:", instanceAddr)
		cc, err := grpc.Dial(instanceAddr, grpc.WithInsecure())
		if err != nil {
			fmt.Println("client create issue", err)
		}
		defer cc.Close()
		// 得到 service
		clientEndpoint := NewClient(cc)
		// 向 server 的 service 发送请求
		responseCTX, v, err := clientEndpoint.Test(context.Background(), "Allen", 21)
		if err != nil {
			fmt.Println("Test func issue", err)
		}
		fmt.Println("response is :", v)
		fmt.Println("context is :", responseCTX)
		return nil, nil
	}, nil, nil
}
```


##### 4. main 函数（完成etcd注册）



```golang
import(
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/etcdv3"
	"github.com/go-kit/kit/sd/lb"
)
```



```golang
func main() {

	//--------------------------------

	var (
		// 注册中心地址
		etcdServer = "10.204.29.77:2379"
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
    // 现在我们可以通过 endPoint 发起请求了
	req := struct{}{}
	for i := 1; i <= 8; i++ {
		if _, err = reqEndPoint(ctx, req); err != nil {
			panic(err)
		}
	}
}

```

##### 5.完整代码


```golang
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
		return nil, nil
	}, nil, nil
}

func main() {

	//--------------------------------

	var (
		// 注册中心地址
		etcdServer = "10.204.29.77:2379"
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

```


<html>
<center><h4>第四部分结束</h4></center>
</html>


---

> 如何使用该服务

#### 一、 我们需要安装etcdv3服务
##### 1.download etcdv3二进制文件

> https://github.com/coreos/etcd/releases/		（etcd-v3.3.9-linux-amd64.tar.gz）

##### 2.解压得到两个文件
> - etcd
> - etcdctl

##### 3. 复制到 /user/local/bin 目录下

##### 4. 配置环境变量

> ETCDCTL_API=3

##### 5.在命令行测试 etcd 
> etcd

##### 6.在命令行测试 etcdctl
> etcdctl

#### 二、 我们需要部署etcdv3服务集群
##### 1.通用配置信息（假设实现三台电脑的etcdv3集群）
###### 每台电脑都需要配置

```
TOKEN=token-07
CLUSTER_STATE=new
NAME_1=machine-1
NAME_2=machine-2
NAME_3=machine-3
HOST_1=10.204.29.77
HOST_2=10.204.29.70
HOST_3=10.204.29.73
CLUSTER=${NAME_1}=http://${HOST_1}:2380,${NAME_2}=http://${HOST_2}:2380,${NAME_3}=http://${HOST_3}:2380
```

##### 2.为各电脑启动 etcdv3
###### 电脑一（IP:10.204.29.77）

```
THIS_NAME=${NAME_1}
THIS_IP=${HOST_1}
etcd --data-dir=data.etcd --name ${THIS_NAME} --initial-advertise-peer-urls http://${THIS_IP}:2380 --listen-peer-urls http://${THIS_IP}:2380 --advertise-client-urls http://${THIS_IP}:2379 --listen-client-urls http://${THIS_IP}:2379 --initial-cluster ${CLUSTER} --initial-cluster-state ${CLUSTER_STATE} --initial-cluster-token ${TOKEN}
```

###### 电脑二（IP:10.204.29.70）


```
THIS_NAME=${NAME_2}
THIS_IP=${HOST_2}
etcd --data-dir=data.etcd --name ${THIS_NAME} --initial-advertise-peer-urls http://${THIS_IP}:2380 --listen-peer-urls http://${THIS_IP}:2380 --advertise-client-urls http://${THIS_IP}:2379 --listen-client-urls http://${THIS_IP}:2379 --initial-cluster ${CLUSTER} --initial-cluster-state ${CLUSTER_STATE} --initial-cluster-token ${TOKEN}
```


###### 电脑三（IP:10.204.29.73）



```
THIS_NAME=${NAME_3}
THIS_IP=${HOST_3}
etcd --data-dir=data.etcd --name ${THIS_NAME} --initial-advertise-peer-urls http://${THIS_IP}:2380 --listen-peer-urls http://${THIS_IP}:2380 --advertise-client-urls http://${THIS_IP}:2379 --listen-client-urls http://${THIS_IP}:2379 --initial-cluster ${CLUSTER} --initial-cluster-state ${CLUSTER_STATE} --initial-cluster-token ${TOKEN}
```

##### 3.查看集群列表（任意一台主机都可以）出现集群列表表示服务成功启动


```
HOST_1=10.204.29.77
HOST_2=10.204.29.70
HOST_3=10.204.29.73
ENDPOINTS=$HOST_1:2379,$HOST_2:2379,$HOST_3:2379

etcdctl --endpoints=$ENDPOINTS member list
```

##### 4.在三台机器上运行统一服务（记得修改IP）

###### server 的 main(修改IP)


```golang
	var (
		//etcd 服务地址
		etcdServer = "10.204.29.77"
		// 服务的信息目录
		prefix = "/services/book/"
		// 当前启动服务实例的地址
		instance = "10.204.29.77:50051"
		// 服务实例注册的路径
		key = prefix + instance
		// 服务实例注册的 val
		value = instance
		ctx   = context.Background()
		// 服务监听地址
		serviceAddress = ":50051"
	)
```

###### 运行服务

> go run main.go

###### 至此各个程序已经处于监听状态

##### 5. 开启客户端 发送请求服务

###### 电脑一（IP 10.204.29.77）client 注册代码



```golang
	var (
		// 注册中心地址
		etcdServer = "10.204.29.77:2379"
		// 监听的服务前缀
		prefix = "/services/book/"
		ctx    = context.Background()
	)
	options := etcdv3.ClientOptions{
		DialTimeout:   time.Second * 3,
		DialKeepAlive: time.Second * 3,
	}
```

###### 启动程序 发送请求

> go run main.go


##### 6. 最终各个程序的执行由电脑一（leader）发送给另外两个follower执行，完成负载均衡


---



##### 详细代码可参阅 gokit example5
##### URL[：http://10.204.28.137/allenguo/go-kit__VS__go-micro/tree/master/gokit/example5](http://10.204.28.137/allenguo/go-kit__VS__go-micro/tree/master/gokit/example5)

###### 编著人：Allen guo
###### 日期：  2018/8/29 