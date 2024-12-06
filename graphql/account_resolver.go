package main

import "context"

type accountResolver struct {
	server *Server
}

func (r *accountResolver) Orders(ctx context.Context, oj *Account) ([]*Order, error) {
	orders := []*Order{}
	return orders, nil
}
