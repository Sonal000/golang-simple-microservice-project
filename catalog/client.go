package catalog

import (
	"context"

	pb "github.com/Sonal000/golang-simple-microservice-project/catalog/pb"
	"google.golang.org/grpc"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.CatalogServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	c := pb.NewCatalogServiceClient(conn)
	return &Client{conn, c}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) PostProduct(ctx context.Context, name, description string, price float64) (*Product, error) {
	product, err := c.service.PostProduct(ctx, &pb.PostProductRequest{
		Name:        name,
		Description: description,
		Price:       price,
	})
	if err != nil {
		return &Product{}, err
	}
	return &Product{
		Id:          product.Product.Id,
		Name:        product.Product.Name,
		Description: product.Product.Description,
		Price:       product.Product.Price,
	}, nil
}

func (c *Client) GetProduct(ctx context.Context, id string) (*Product, error) {
	product, err := c.service.GetProduct(ctx, &pb.GetProductRequest{Id: id})
	if err != nil {
		return &Product{}, nil
	}
	return &Product{
		Id:          product.Product.Id,
		Name:        product.Product.Name,
		Description: product.Product.Description,
		Price:       product.Product.Price,
	}, nil
}

func (c *Client) GetProducts(ctx context.Context, take, skip uint64) ([]Product, error) {
	prods, err := c.service.GetProducts(ctx, &pb.GetProductsRequest{
		Take: take,
		Skip: skip,
	})
	if err != nil {
		return []Product{}, err
	}
	var products []Product
	for _, product := range prods.Products {
		products = append(products, Product{
			Id:          product.Id,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
		})
	}
	return products, err
}
