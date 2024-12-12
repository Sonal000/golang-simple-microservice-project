package catalog

import (
	"context"
	"fmt"
	"net"

	pb "github.com/Sonal000/golang-simple-microservice-project/catalog/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	pb.UnimplementedAccountServiceServer
	catalogService CatalogService
}

func ListenGRPC(s CatalogService, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	serv := grpc.NewServer()
	pb.RegisterAccountServiceServer(serv, &grpcServer{
		UnimplementedAccountServiceServer: pb.UnimplementedAccountServiceServer{},
		catalogService:                    s,
	})
	reflection.Register(serv)
	return serv.Serve(lis)
}

func (s *grpcServer) PostProduct(ctx context.Context, r *pb.PostProductRequest) (*pb.PostProductResponse, error) {
	a, err := s.catalogService.PostProduct(ctx, r.name, r.description, r.price)
	if err != nil {
		return nil, err
	}
	return &pb.PostProductResponse{
		Product: &pb.Product{
			Id:          a.Id,
			Name:        a.Name,
			Description: a.Description,
			Price:       a.Price,
		},
	}, nil
}

func (s *grpcServer) GetAccount(ctx context.Context, r *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	a, err := s.catalogService.GetProduct(ctx, r.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetProductResponse{
		Account: &pb.Product{
			Id:          a.Id,
			Name:        a.Name,
			Description: a.Description,
			Price:       a.Price,
		},
	}, nil

}

func (s *grpcServer) GetProducts(ctx context.Context, r *pb.GetProductsRequest) (*pb.GetProductsResponse, error) {
	skip := r.Skip
	take := r.Take
	res, err := s.catalogService.GetProducts(ctx, skip, take)
	if err != nil {
		return nil, err
	}

	products := []*pb.Product{}
	for _, p := range res {
		products = append(products, &pb.Product{
			Id:          p.Id,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		})
	}
	return &pb.GetProductsResponse{
		Products: products,
	}, nil
}
