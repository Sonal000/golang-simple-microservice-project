package account

import (
	"context"

	pb "github.com/Sonal000/golang-simple-microservice-project/account/pb"
	"google.golang.org/grpc"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.AccountServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	c := pb.NewAccountServiceClient(conn)
	return &Client{conn, c}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) PostAccount(ctx context.Context, name string) (*Account, error) {
	acc, err := c.service.PostAccount(ctx, &pb.PostAccountRequest{Id: name})
	if err != nil {
		return &Account{}, err
	}
	return &Account{
		ID:   acc.Account.Id,
		Name: acc.Account.Name,
	}, nil
}

func (c *Client) GetAccount(ctx context.Context, id string) (*Account, error) {
	acc, err := c.service.GetAccount(ctx, &pb.GetAccountRequest{Id: id})
	if err != nil {
		return &Account{}, nil
	}
	return &Account{
		ID:   acc.Account.Id,
		Name: acc.Account.Name,
	}, nil
}

func (c *Client) GetAccounts(ctx context.Context, take, skip uint64) ([]Account, error) {
	accs, err := c.service.GetAccounts(ctx, &pb.GetAccountsRequest{
		Take: take,
		Skip: skip,
	})
	if err != nil {
		return []Account{}, err
	}
	var accounts []Account
	for _, a := range accs.Accounts {
		accounts = append(accounts, Account{
			ID:   a.Id,
			Name: a.Name,
		})
	}
	return accounts, err
}
