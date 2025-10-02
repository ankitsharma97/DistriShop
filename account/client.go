package account

import (
	"context"
	pb "microservice/account/pb"

	"google.golang.org/grpc"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.AccountServiceClient
}

func NewClient(address string) (*Client, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	service := pb.NewAccountServiceClient(conn)
	return &Client{conn: conn, service: service}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) PostAccount(ctx context.Context, name string) (*Account, error) {
	req := &pb.PostAccountRequest{Name: name}
	resp, err := c.service.PostAccount(ctx, req)
	if err != nil {
		return nil, err
	}
	return &Account{ID: resp.Account.Id, Name: resp.Account.Name}, nil
}

func (c *Client) GetAccount(ctx context.Context, id string) (*Account, error) {
	req := &pb.GetAccountRequest{Id: id}
	resp, err := c.service.GetAccount(ctx, req)
	if err != nil {
		return nil, err
	}
	return &Account{ID: resp.Account.Id, Name: resp.Account.Name}, nil
}

func (c *Client) GetAccounts(ctx context.Context, skip int, take int) ([]*Account, error) {
	req := &pb.GetAccountsRequest{Skip: uint64(skip), Take: uint64(take)}
	resp, err := c.service.GetAccounts(ctx, req)
	if err != nil {
		return nil, err
	}
	return convertAccounts(resp.Accounts), nil
}

func convertAccounts(accounts []*pb.Account) []*Account {
	var result []*Account
	for _, acc := range accounts {
		result = append(result, &Account{ID: acc.Id, Name: acc.Name})
	}
	return result
}
