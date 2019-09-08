// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/dice.proto

package dice

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	client "github.com/micro/go-micro/client"
	server "github.com/micro/go-micro/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for Dice service

type DiceService interface {
	Roll(ctx context.Context, in *RollRequest, opts ...client.CallOption) (*RollResponse, error)
}

type diceService struct {
	c    client.Client
	name string
}

func NewDiceService(name string, c client.Client) DiceService {
	if c == nil {
		c = client.NewClient()
	}
	if len(name) == 0 {
		name = "dice"
	}
	return &diceService{
		c:    c,
		name: name,
	}
}

func (c *diceService) Roll(ctx context.Context, in *RollRequest, opts ...client.CallOption) (*RollResponse, error) {
	req := c.c.NewRequest(c.name, "Dice.Roll", in)
	out := new(RollResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Dice service

type DiceHandler interface {
	Roll(context.Context, *RollRequest, *RollResponse) error
}

func RegisterDiceHandler(s server.Server, hdlr DiceHandler, opts ...server.HandlerOption) error {
	type dice interface {
		Roll(ctx context.Context, in *RollRequest, out *RollResponse) error
	}
	type Dice struct {
		dice
	}
	h := &diceHandler{hdlr}
	return s.Handle(s.NewHandler(&Dice{h}, opts...))
}

type diceHandler struct {
	DiceHandler
}

func (h *diceHandler) Roll(ctx context.Context, in *RollRequest, out *RollResponse) error {
	return h.DiceHandler.Roll(ctx, in, out)
}
