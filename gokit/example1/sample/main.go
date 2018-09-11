package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

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

func main() {

	//workflow
	// 	1.	define service struct
	ys := youServiceStruct{}

	//	2.	bound func
	//	decodeYouServiceFuncNameRequest		run   get a endpoint.Endpoint
	printVisitorNameHandler := httptransport.NewServer(
		makeYouServiceFuncNameEndpoint(ys),
		decodeYouServiceFuncNameRequest,
		encodeYouServiceFuncNameResponse,
	)

	// 	3.	set router
	http.Handle("/", printVisitorNameHandler)

	// 	4.	set listen port
	log.Fatal(http.ListenAndServe(":8080", nil))

	//  5.  makeYouServiceFuncNameEndpoint      run

	//  6.  request ------->router
	//  7.  decodeYouServiceFuncNameRequest     run
	//  8.  endpoint.Endpoint       run
	//  9.  YouServiceFuncName  run
	//  10. encodeYouServiceFuncNameResponse    run
}
