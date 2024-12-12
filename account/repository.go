package account

import (
	"context"
	"database/sql"
)

type Repository interface {
	Close() error
	PutAccount(ctx context.Context, name string) (Account, error)
	GetAccountById(ctx context.Context, id string) (Account, error)
	ListAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error)
}

type postgresqlRepository struct {
	db *sql.DB
}

func NewAccountRepository(url string) (Repository, error) {

	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &postgresqlRepository{db: db}, nil
}

func (r *postgresqlRepository) close() {
	r.db.Close()
}

func (r *postgresqlRepository) PutAccount(ctx context.Context, name string) (Account, error) {
	var account Account
	err := r.db.QueryRowContext(ctx, "INSERT INTO accounts (name) VALUES ($1) RETURNING id", name).Scan(&account.ID)
	if err != nil {
		return Account{}, err
	}
	account.Name = name
	return account, nil
}
func (r *postgresqlRepository) GetAccountById(ctx context.Context, id string) (Account, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id, name FROM accounts WHERE id = $1", id)
	var a = Account{}
	if err := row.Scan(&a.ID, &a.Name); err != nil {
		return Account{}, err
	}
	return a, nil
}

func (r *postgresqlRepository) ListAccounts(ctx context.Context, take uint64, skip uint64) ([]Account, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, name FROM accounts ORDER BY id DESC OFFSET $1 LIMIT $2", skip, take)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []Account
	for rows.Next() {
		var a Account
		if err := rows.Scan(&a.ID, &a.Name); err != nil {
			return nil, err
		}
		accounts = append(accounts, a)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}

func (r *postgresqlRepository) Close() error {
	return r.db.Close()
}
