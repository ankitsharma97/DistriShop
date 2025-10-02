package account

import (
	"context"

	"github.com/segmentio/ksuid"
)


// protoc --go_out=./pb --go-grpc_out=./pb account.proto

// Service defines account business operations.
type Service interface {
	PostAccount(ctx context.Context, name string) (*Account, error)
	GetAccount(ctx context.Context, id string) (*Account, error)
	GetAccounts(ctx context.Context, skip int, take int) ([]*Account, error)
}

// Account domain model.
type Account struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// accountService implements Service using a Repository.
type accountService struct {
	repo Repository
}

// NewService constructs a Service.
func NewService(repo Repository) Service {
	return &accountService{repo: repo}
}

// PostAccount creates and stores a new account.
func (s *accountService) PostAccount(ctx context.Context, name string) (*Account, error) {
	account := &Account{ID: ksuid.New().String(), Name: name}
	if err := s.repo.PutAccount(ctx, account); err != nil {
		return nil, err
	}
	return account, nil
}

// GetAccount fetches an account by ID.
func (s *accountService) GetAccount(ctx context.Context, id string) (*Account, error) {
	return s.repo.GetAccountById(ctx, id)
}

// GetAccounts lists accounts with basic bounds on pagination.
func (s *accountService) GetAccounts(ctx context.Context, skip int, take int) ([]*Account, error) {
	if take <= 0 {
		take = 10
	}
	if take > 100 {
		take = 100
	}
	return s.repo.ListsAccounts(ctx, skip, take)
}
