package account

import (
	"context"
	"log"
)

type accountService struct {
	repository Repository
}

type AccountService interface {
	CreateAccount(ctx context.Context, name string) (Account, error)
	GetAccount(ctx context.Context, id string) (Account, error)
	GetAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error)
}

type Account struct {
	ID   string
	Name string
}

func NewAccountService() AccountService {
	repo, err := NewAccountRepository("db_url")
	if err != nil {
		log.Fatal(err)
	}

	return &accountService{
		repository: repo,
	}
}

func (s *accountService) CreateAccount(ctx context.Context, name string) (Account, error) {
	account, err := s.repository.PutAccount(ctx, name)
	if err != nil {
		return Account{}, err
	}
	return account, nil
}

func (s *accountService) GetAccount(ctx context.Context, id string) (Account, error) {
	var a = Account{}
	a, err := s.repository.GetAccountById(ctx, id)
	if err != nil {
		return Account{}, err
	}
	return a, nil
}

func (s *accountService) GetAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error) {
	if (take > 100) || skip == 0 || take == 0 {
		take = 100
	}
	return s.repository.ListAccounts(ctx, skip, take)
}
