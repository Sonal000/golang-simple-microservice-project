package catalog

import (
	"context"

	"github.com/segmentio/ksuid"
)

type Product struct {
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type Service interface {
	PostProduct(ctx context.Context, name, description string, price float64) (*Product, error)
	GetProduct(ctx context.Context, id string) (*Product, error)
	GetProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error)
	GetProductsByIDs(ctx context.Context, ids []string) ([]Product, error)
	SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]Product, error)
}

type catalogService struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &catalogService{r}
}

func (s *catalogService) PostProduct(ctx context.Context, name, description string, price float64) (*Product, error) {
	p := &Product{
		Name:        name,
		Description: description,
		Price:       price,
		Id:          ksuid.New().String(),
	}
	product, err := s.repository.CreateProduct(ctx, *p)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *catalogService) GetProduct(ctx context.Context, id string) (*Product, error) {
	return s.repository.GetProductById(ctx, id)
}

func (s *catalogService) GetProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error) {
	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}
	return s.repository.ListProducts(ctx, skip, take)
}

func (s *catalogService) GetProductsByIDs(ctx context.Context, ids []string) ([]Product, error) {
	return s.repository.ListProductsWithIds(ctx, ids)
}

func (s *catalogService) SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]Product, error) {
	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}
	return s.repository.SearchProducts(ctx, query, skip, take)
}
