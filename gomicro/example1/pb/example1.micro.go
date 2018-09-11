// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: example1.proto

/*
Package pb is a generated protocol buffer package.

It is generated from these files:
	example1.proto

It has these top-level messages:
	Request
	Response
*/
package pb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	client "github.com/micro/go-micro/client"
	server "github.com/micro/go-micro/server"
	context "context"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for Test service

type TestService interface {
	Hello(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error)
}

type testService struct {
	c    client.Client
	name string
}

func NewTestService(name string, c client.Client) TestService {
	if c == nil {
		c = client.NewClient()
	}
	if len(name) == 0 {
		name = "pb"
	}
	return &testService{
		c:    c,
		name: name,
	}
}

func (c *testService) Hello(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.name, "Test.Hello", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Test service

type TestHandler interface {
	Hello(context.Context, *Request, *Response) error
}

func RegisterTestHandler(s server.Server, hdlr TestHandler, opts ...server.HandlerOption) error {
	type test interface {
		Hello(ctx context.Context, in *Request, out *Response) error
	}
	type Test struct {
		test
	}
	h := &testHandler{hdlr}
	return s.Handle(s.NewHandler(&Test{h}, opts...))
}

type testHandler struct {
	TestHandler
}

func (h *testHandler) Hello(ctx context.Context, in *Request, out *Response) error {
	return h.TestHandler.Hello(ctx, in, out)
}