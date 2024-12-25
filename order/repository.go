package order

import (
	"context"
	"database/sql"
)

type Repository interface {
	Close()
	CreateOrder(ctx context.Context, order Order) (Order, error)
	GetOrdersForAccount(ctx context.Context, id string) ([]Order, error)
}

type orderRepository struct {
	db *sql.DB
}

func NewOrderRepository(url string) (Repository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &orderRepository{db}, nil
}

type OrderDocument struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

func (r *orderRepository) Close() {
	r.db.Close()
}

func (r *orderRepository) CreateOrder(ctx context.Context, order Order) (Order, error) {
	var newOrder Order
	err := r.db.QueryRowContext(ctx, "INSERT INTO orders (name) VALUES ($1) RETURNING id", order.AccountID).Scan(&newOrder)
	if err != nil {
		return Order{}, err
	}
	return newOrder, nil
}
func (r *orderRepository) GetOrdersForAccount(ctx context.Context, id string) ([]Order, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, name FROM accounts ORDER BY id DESC OFFSET $1 LIMIT $2")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var order Order
		if err := rows.Scan(&order); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}
