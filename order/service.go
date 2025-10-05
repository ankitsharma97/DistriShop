package order

import (
	"context"
	"time"

	"github.com/segmentio/ksuid"
)

type Service interface {
	PostOrder(ctx context.Context, accountID string, products []*OrderedProduct) (*Order, error)
	GetOrdersForAccount(ctx context.Context, accountID string) ([]*Order, error)
}

type Order struct {
	ID        string            `json:"id"`
	AccountID string            `json:"account_id"`
	CreatedAt time.Time         `json:"created_at"`
	Total     float64           `json:"total"`
	Products  []*OrderedProduct `json:"products"`
}

type OrderedProduct struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type orderService struct {
	repo Repository
}

func NewOrderService(repo Repository) Service {
	return &orderService{repo: repo}
}

func (s *orderService) PostOrder(ctx context.Context, accountID string, products []*OrderedProduct) (*Order, error) {
	order := &Order{
		ID:        ksuid.New().String(),
		AccountID: accountID,
		CreatedAt: time.Now().UTC(),
		Products:  products,
	}
	// Note: Total calculation should be done at the server layer where we have product prices
	order.Total = 0.0

	if err := s.repo.PutOrder(ctx, order); err != nil {
		return nil, err
	}
	return order, nil
}

func (s *orderService) GetOrdersForAccount(ctx context.Context, accountID string) ([]*Order, error) {
	return s.repo.GetOrdersForAccount(ctx, accountID)
}
