package grpcclient

import (
	"context"
	"log"

	"api-gateway/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClientInterface interface {
	RegisterPhone(ctx context.Context, req *proto.RegisterPhoneRequest) (*proto.RegisterPhoneResponse, error)
	UnregisterPhone(ctx context.Context, req *proto.UnregisterPhoneRequest) (*proto.ActionResponse, error)
	ListPhones(ctx context.Context, req *proto.ListPhonesRequest) (*proto.PhoneListResponse, error)
	Call(ctx context.Context, req *proto.CallRequest) (*proto.ActionResponse, error)
	AnswerCall(ctx context.Context, req *proto.CallActionRequest) (*proto.ActionResponse, error)
	HangupCall(ctx context.Context, req *proto.CallActionRequest) (*proto.ActionResponse, error)
	HoldCall(ctx context.Context, req *proto.CallActionRequest) (*proto.ActionResponse, error)
	ResumeCall(ctx context.Context, req *proto.CallActionRequest) (*proto.ActionResponse, error)
}

type GRPCClient struct {
	conn   *grpc.ClientConn
	client proto.SoftPhoneServiceClient
}

func NewGRPCClient(address string) (*GRPCClient, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := proto.NewSoftPhoneServiceClient(conn)
	return &GRPCClient{conn: conn, client: client}, nil
}

func (c *GRPCClient) Close() {
	if err := c.conn.Close(); err != nil {
		log.Printf("Failed to close gRPC connection: %v", err)
	}
}

func (c *GRPCClient) RegisterPhone(ctx context.Context, req *proto.RegisterPhoneRequest) (*proto.RegisterPhoneResponse, error) {
	return c.client.RegisterPhone(ctx, req)
}

func (c *GRPCClient) ListPhones(ctx context.Context, req *proto.ListPhonesRequest) (*proto.PhoneListResponse, error) {
	return c.client.ListPhones(ctx, req)
}

func (c *GRPCClient) Call(ctx context.Context, req *proto.CallRequest) (*proto.ActionResponse, error) {
	return c.client.Call(ctx, req)
}

func (c *GRPCClient) AnswerCall(ctx context.Context, req *proto.CallActionRequest) (*proto.ActionResponse, error) {
	return c.client.AnswerCall(ctx, req)
}

func (c *GRPCClient) HangupCall(ctx context.Context, req *proto.CallActionRequest) (*proto.ActionResponse, error) {
	return c.client.HangupCall(ctx, req)
}

func (c *GRPCClient) HoldCall(ctx context.Context, req *proto.CallActionRequest) (*proto.ActionResponse, error) {
	return c.client.HoldCall(ctx, req)
}

func (c *GRPCClient) ResumeCall(ctx context.Context, req *proto.CallActionRequest) (*proto.ActionResponse, error) {
	return c.client.ResumeCall(ctx, req)
}

func (c *GRPCClient) UnregisterPhone(ctx context.Context, req *proto.UnregisterPhoneRequest) (*proto.ActionResponse, error) {
	return c.client.UnregisterPhone(ctx, req)
}
