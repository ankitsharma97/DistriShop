package main

import (
	"context"
	"time"
)

type queryResolver struct {
	server *Server
}

func (r *queryResolver) Accounts(ctx context.Context, pagination *PaginationInput, id *string) ([]*Account, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if id != nil && *id != "" {
		account, err := r.server.accountClient.GetAccount(ctx, *id)
		if err != nil {
			return nil, err
		}
		return []*Account{&Account{ID: account.ID, Username: account.Name}}, nil
	}

	skip, take := 0, 10
	if pagination != nil {
		if pagination.Skip != nil {
			skip = *pagination.Skip
		}
		if pagination.Take != nil {
			take = *pagination.Take
		}
	}

	accounts, err := r.server.accountClient.GetAccounts(ctx, skip, take)
	if err != nil {
		return nil, err
	}

	var result []*Account
	for _, a := range accounts {
		result = append(result, &Account{ID: a.ID, Username: a.Name})
	}
	return result, nil
}

func (r *queryResolver) Products(ctx context.Context, pagination *PaginationInput, query *string, id *string) ([]*Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if id != nil && *id != "" {
		product, err := r.server.catalogClient.GetProduct(ctx, *id)
		if err != nil {
			return nil, err
		}
		return []*Product{&Product{
			ID:          product.ID,
			Name:        product.Name,
			Description: &product.Description,
			Price:       product.Price,
		}}, nil
	}

	skip, take := 0, 10
	if pagination != nil {
		if pagination.Skip != nil {
			skip = *pagination.Skip
		}
		if pagination.Take != nil {
			take = *pagination.Take
		}
	}

	searchQuery := ""
	if query != nil {
		searchQuery = *query
	}

	products, err := r.server.catalogClient.GetProducts(ctx, []string{}, searchQuery, skip, take)
	if err != nil {
		return nil, err
	}

	var result []*Product
	for _, p := range products {
		result = append(result, &Product{
			ID:          p.ID,
			Name:        p.Name,
			Description: &p.Description,
			Price:       p.Price,
		})
	}
	return result, nil
}
