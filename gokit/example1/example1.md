example1

*定义一个服务 interface
*设置一些服务需要具备的方法
*定义一个结构体
*实现对应的方法

    type YouServiceName interface {
        YouServiceFuncName( /* you need params */ ) /* (you need renturn value) */
        //....
    }
    type youServiceStruct struct {
        youServiceAttributeName string
        //...
    }
    // you need to follow the interface (you define service)
    func (ys youServiceStruct) YouServiceFuncName() {
        // func content
        //...
    }



*定义一个factory去生产endpoint
*其中可以调用需要的方法
*定义一个解码器，一个编码器用于拆包和装包 request和response
*根据需要还可定义request和response的Json样式
*并在解码和编码器中进行校验

    func makeYouServiceFuncNameEndpoint(ys YouServiceName) endpoint.Endpoint {
        return func(ctx context.Context, request interface{}) (interface{}, error) {
            // you can use you func  YouServiceFuncName
            ys.YouServiceFuncName()
            //return  you need
            return nil, nil
        }
    }

    func decodeYouServiceFuncNameRequest(_ context.Context, r *http.Request) (interface{}, error) {
        // deal request
        fmt.Println("request method  :", r.Method)
        return nil, nil
    }
        
    func encodeYouServiceFuncNameResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
        // set response
        fmt.Println("Hello Visitor")
        return nil
    }
    //request struct
    //response struct



*创建一个服务器结构对象
*构造一个endpoint的对象
*设置endpoint对应的路由
*开启监听

    func main() {
        //workflow
        //  1.  define service struct
        ys := youServiceStruct{}
        //  2.  bound func
        //  decodeYouServiceFuncNameRequest     run   get a endpoint.Endpoint
        printVisitorNameHandler := httptransport.NewServer(
            makeYouServiceFuncNameEndpoint(ys),
            decodeYouServiceFuncNameRequest,
            encodeYouServiceFuncNameResponse,
        )
        //  3.  set router
        http.Handle("/", printVisitorNameHandler)
        //  4.  set listen port
        log.Fatal(http.ListenAndServe(":8080", nil))
        	//	5.	makeYouServiceFuncNameEndpoint    	run
        	
        	// 	6.	request ------->router
        	// 	7.	decodeYouServiceFuncNameRequest   	run
        	//	8.	endpoint.Endpoint    	run
        	// 	9.	YouServiceFuncName	run
        	// 	10.	encodeYouServiceFuncNameResponse	run
    }



