Example3（wrappers中间件使用）

proto文件

    syntax = "proto3";

    package pb;

    service Wrappers {
        rpc Test(Request) returns (Response) {}
    }

    message Request {
        string name = 1;
    }

    message Response {
        string msg = 1;
    }

cli

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

server

    type TestWrappers struct{}

    func (tw *TestWrappers) Test(ctx context.Context, req *pb.Request, rep *pb.Response) error {
        log.Println("This is a test")
        log.Println("client:my name is ", req.Name)
        rep.Msg=fmt.Sprintf("\nserver: welcome %s\nserver: %s Bye~\n" ,req.Name,req.Name)
        return nil
    }

    func logWrapper(fn server.HandlerFunc) server.HandlerFunc {
        return func(ctx context.Context, req server.Request, rsp interface{}) error {
            log.Printf("request method: %v\n", req.Method())
            err := fn(ctx, req, rsp)
            return err
        }
    }

    func main() {
        service := micro.NewService(
            micro.Name("wrap"),
            micro.WrapHandler(logWrapper),
        )
        // optionally setup command line usage
        service.Init()
        // Register Handlers
        pb.RegisterWrappersHandler(service.Server(), new(TestWrappers))
        // Run server
        if err := service.Run(); err != nil {
            log.Fatal(err)
        }
    }


