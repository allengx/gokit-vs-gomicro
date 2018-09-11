Example1（简单通信grpc）

生成 .pb.go  &&  .micro.go 文件

编写 .proto 文件

    syntax = "proto3";
    package pb;
    service Test {
       	 rpc Hello(Request) returns (Response) {}
    }
    message Request {
        string name = 1;
    }
    message Response {
      	  string msg = 1;
    }

protoc --proto_path=$GOPATH/src:. --micro_out=. --go_out=. xxxxx.proto（得到.pb.go&&.micro.go 文件）

根据生成文件编写代码

client

    func main() {
       	 service := micro.NewService(micro.Name("aaa"))
        	service.Init()
       	 // use the generated client stub
        	cl := test.NewTestService("test", service.Client())
       	 // Set arbitrary headers in context
        	ctx := metadata.NewContext(context.Background(), make(map[string]string))
       	 rsp, err := cl.Hello(ctx, &test.Request{
       	     Name: "aohn",
       	 })
       	 if err != nil {
        	    fmt.Println(err)
       	     return
       	 }
        	fmt.Println(rsp.Msg)
    }

server

    type TestHello struct{}
    func (s *TestHello) Hello(ctx context.Context, req *test.Request, rsp *test.Response) error {
        	log.Print("Received Say.Hello request")
        	rsp.Msg = "Hello " + req.Name
       	 return nil
    }
    func main() {
        	service := micro.NewService(
           	 	micro.Name("test"),
            	micro.RegisterTTL(time.Second*30),
            	micro.RegisterInterval(time.Second*10),
        	)
        	// optionally setup command line usage
        	service.Init()
       	 // Register Handlers
        	test.RegisterTestHandler(service.Server(), new(TestHello))
        	// Run server
        	if err := service.Run(); err != nil {
            	log.Fatal(err)
       	 }
    }

运行程序（MICRO_REGISTRY=mdns go run main.go）
