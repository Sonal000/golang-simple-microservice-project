package account

import (
	"context"
	"fmt"
	"net"

	pb "github.com/Sonal000/golang-simple-microservice-project/account/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	pb.UnimplementedAccountServiceServer
	accountService AccountService
}

func ListenGRPC(s AccountService, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	serv := grpc.NewServer()
	pb.RegisterAccountServiceServer(serv, &grpcServer{
		UnimplementedAccountServiceServer: pb.UnimplementedAccountServiceServer{},
		accountService:                    s,
	})
	reflection.Register(serv)
	return serv.Serve(lis)
}

func (s *grpcServer) PostAccount(ctx context.Context, r *pb.PostAccountRequest) (*pb.PostAccountResponse, error) {
	a, err := s.accountService.CreateAccount(ctx, r.Id)
	if err != nil {
		return nil, err
	}
	return &pb.PostAccountResponse{
		Account: &pb.Account{
			Id:   a.ID,
			Name: a.Name,
		},
	}, nil
}

func (s *grpcServer) GetAccount(ctx context.Context, r *pb.GetAccountRequest) (*pb.GetAccountResponse, error) {
	a, err := s.accountService.GetAccount(ctx, r.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetAccountResponse{
		Account: &pb.Account{
			Id:   a.ID,
			Name: a.Name,
		},
	}, nil

}

func (s *grpcServer) GetAccounts(ctx context.Context, r *pb.GetAccountsRequest) (*pb.GetAccountsResponse, error) {
	skip := r.Skip
	take := r.Take
	res, err := s.accountService.GetAccounts(ctx, skip, take)
	if err != nil {
		return nil, err
	}

	accounts := []*pb.Account{}
	for _, acc := range res {
		accounts = append(accounts, &pb.Account{
			Id:   acc.ID,
			Name: acc.Name,
		})
	}
	return &pb.GetAccountsResponse{
		Accounts: accounts,
	}, nil
}
