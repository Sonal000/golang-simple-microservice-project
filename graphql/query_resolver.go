package main

import (
	"context"
	"log"
	"time"
)

type queryResolver struct {
	server *Server
}

func (r *queryResolver) Accounts(ctx context.Context, pagination *PaginationInput, id *string) ([]*Account, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	if id != nil {
		r, err := r.server.accountClient.GetAccount(ctx, *id)
		if err != nil {
			log.Println(err)
			return []*Account{}, err
		}
		return []*Account{{
			Id:   r.ID,
			Name: r.Name,
		}}, nil
	}
	accounts := []*Account{}
	return accounts, nil
}
func (r *queryResolver) Products(ctx context.Context, pagination *PaginationInput, query *string, id *string) ([]*Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	if id != nil {
		r, err := r.server.catalogClient.GetProduct(ctx, *id)
		if err != nil {
			log.Println(err)
			return []*Product{}, err
		}
		return []*Product{{
			ID:          r.Id,
			Name:        r.Name,
			Price:       r.Price,
			Description: r.Description,
		}}, nil
	}
	products := []*Product{}
	return products, nil
}

// func (p PaginationInput) bounds() (uint64, uint64) {
// 	skipValue := uint64(0)
// 	takeValue := uint64(100)
// 	// if p.Skip != nil {
// 	// 	skipValue = uint64(*p.Skip)
// 	// }
// 	// if p.Take != nil {
// 	// 	takeValue = uint64(*p.Take)
// 	// }
// 	return skipValue, takeValue
// }
