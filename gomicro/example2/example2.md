Example2（plugins使用）

proto文件

    syntax = "proto3";

    package pb;
    // Example message
    message Msg {
        	// unique id
       	 int64 id=1;
       	 // message
       	 string context = 2;
    }

cli

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

server

        import (
            pb "gomicro/example2/pb"
    
            "github.com/micro/go-log"
            "github.com/micro/go-micro"
            "github.com/micro/go-micro/metadata"
    
            "context"
            _ "github.com/micro/go-plugins/broker/rabbitmq"
            _ "github.com/micro/go-plugins/registry/kubernetes"
            _ "github.com/micro/go-plugins/transport/nats"
        )
    
        // All methods of Sub will be executed when
        // a message is received
        type Sub struct{}

        //Method can be of any name
        func (s *Sub) Process(ctx context.Context, msg *pb.Msg) error {
            md, _ := metadata.FromContext(ctx)
            log.Logf("get %+v with %+v\n", msg, md)
            return nil
        }
    
        func main() {
            // create a service
            service := micro.NewService(
                micro.Name("pubsub"),
            )
            // parse command line
            service.Init()
            //register subscriber                                       and run sub full func
            micro.RegisterSubscriber("publish--1--:", service.Server(), new(Sub))
            // one publish msg and all subscriber can get
            if err := service.Run(); err != nil {
                log.Fatal(err)
            }
        }

调用插件两种方法

    MICRO_REGISTRY=mdns MICRO_BROKER=http go run main.go
    
 	MICRO_REGISTRY=mdns \
    > MICRO_BROKER=http \
    > go run main.go