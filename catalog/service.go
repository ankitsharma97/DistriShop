package catalog

import (
	"context"
	"errors"

	"github.com/segmentio/ksuid"
)

type Service interface {
	PostProduct(ctx context.Context, name string, description string, price float64) (*Product, error)
	GetProduct(ctx context.Context, id string) (*Product, error)
	GetProductByIDs(ctx context.Context, ids []string) ([]*Product, error)
	GetProducts(ctx context.Context, skip int, take int) ([]*Product, error)
	SearchProducts(ctx context.Context, query string, skip int, take int) ([]*Product, error)
}

type Product struct {
	ID          string	`json:"id"`
	Name        string	`json:"name"`
	Description string	`json:"description"`
	Price       float64	`json:"price"`
}


type CatalogService struct {
	repo Repository
}

func NewCatalogService(repo Repository) *CatalogService {
	return &CatalogService{repo: repo}
}

func (s *CatalogService) PostProduct(ctx context.Context, name string, description string, price float64) (*Product, error) {
	product := &Product{
		ID:          ksuid.New().String(),
		Name:        name,
		Description: description,
		Price:       price,
	}
	err := s.repo.PutProduct(ctx, product)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *CatalogService) GetProduct(ctx context.Context, id string) (*Product, error) {
	return s.repo.GetProductById(ctx, id)
}

func (s *CatalogService) GetProductByIDs(ctx context.Context, ids []string) ([]*Product, error) {
	return s.repo.ListsProductsWithIDs(ctx, ids)
}

func (s *CatalogService) GetProducts(ctx context.Context, skip int, take int) ([]*Product, error) {

	if skip < 0 || take <= 0 {
		return nil, errors.New("invalid pagination parameters")
	}
	// Adjust the take parameter to a maximum value
	if take > 100 {
		take = 100
	}
	return s.repo.ListsProducts(ctx, skip, take)
}

func (s *CatalogService) SearchProducts(ctx context.Context, query string, skip int, take int) ([]*Product, error) {

	return s.repo.SearchProducts(ctx, query, skip, take)
}


