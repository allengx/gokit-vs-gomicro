package main

import (
	"fmt"
	test "gomicro/example1/pb"
	"github.com/micro/go-micro"
	"context"
	"github.com/micro/go-micro/metadata"
)

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
