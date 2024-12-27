package main

import (
	"context"
	"log"
	"time"
)

type accountResolver struct {
	server *Server
}

func (r *accountResolver) Orders(ctx context.Context, obj *Account) ([]*Order, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	orderList, err := r.server.orderClient.GetOrdersForAccount(ctx, *&obj.Id)
	if err != nil {
		log.Println(err)
		return []*Order{}, err
	}
	var orders []*Order
	for _, o := range orderList {
		var oProducts []*OrderedProduct
		for _, op := range o.Products {
			oProducts = append(oProducts, &OrderedProduct{
				ID:          op.ID,
				Description: op.Description,
				Name:        op.Name,
				Price:       op.Price,
				Quantity:    int(op.Quantity),
			})
		}
		orders = append(orders, &Order{
			ID:         o.ID,
			TotalPrice: o.TotalPrice,
			Products:   oProducts,
		})
	}
	return orders, nil
}
