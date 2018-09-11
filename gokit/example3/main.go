package main

import (
	"context"
	"flag"
	"net/http"
	"os"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
)

func main() {

	//set listen and proxy
	// listen default value is :8080
	// proxy default value is ""
	var (
		listen = flag.String("listen", ":8080", "HTTP listen address")
		proxy  = flag.String("proxy", "", "Optional comma-separated list of URLs to proxy uppercase requests")
	)
	flag.Parse()

	//set log css
	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "listen", *listen, "caller", log.DefaultCaller)

	// create a stringService and add a middlewares to Service discovery and load balancing
	var svc StringService
	svc = stringService{}
	// add middlewares
	svc = proxyingMiddleware(context.Background(), *proxy, logger)(svc)

	// get a Handler and add a uppercase Func
	uppercaseHandler := httptransport.NewServer(
		makeUppercaseEndpoint(svc),
		decodeUppercaseRequest,
		encodeResponse,
	)

	// set a router 
	http.Handle("/uppercase", uppercaseHandler)

	// log listen port and err
	logger.Log("msg", "HTTP", "addr", *listen)
	logger.Log("err", http.ListenAndServe(*listen, nil))
}
