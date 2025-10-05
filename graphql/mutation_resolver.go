package main

import (
	"context"
	"errors"
	"microservice/order"
	"time"
)

var (
	ErrValidParameters = errors.New("valid parameters are required")
)

type mutationResolver struct {
	server *Server
}

func (r *mutationResolver) CreateAccount(ctx context.Context, input AccountInput) (*Account, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	if input.Username == "" {
		return nil, ErrValidParameters
	}
	account, err := r.server.accountClient.PostAccount(ctx, input.Username)
	if err != nil {
		return nil, err
	}
	return &Account{ID: account.ID, Username: account.Name}, nil
}

func (r *mutationResolver) CreateProduct(ctx context.Context, input ProductInput) (*Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	var description string
	if input.Description != nil {
		description = *input.Description
	}
	if input.Name == "" || input.Price <= 0 {
		return nil, ErrValidParameters
	}
	product, err := r.server.catalogClient.PostProduct(ctx, input.Name, description, input.Price)
	if err != nil {
		return nil, err
	}
	return &Product{
		ID:          product.ID,
		Name:        product.Name,
		Description: &product.Description,
		Price:       product.Price,
	}, nil
}

func (r *mutationResolver) CreateOrder(ctx context.Context, input OrderInput) (*Order, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	if input.AccountID == "" || len(input.Products) == 0 {
		return nil, ErrValidParameters
	}
	var products []*order.OrderedProduct
	for _, p := range input.Products {
		if p.ID == "" || p.Quantity <= 0 {
			return nil, ErrValidParameters
		}
		products = append(products, &order.OrderedProduct{
			ProductID: p.ID,
			Quantity:  p.Quantity,
		})
	}

	orderResult, err := r.server.orderClient.PostOrder(ctx, input.AccountID, products)
	if err != nil {
		return nil, err
	}

	// Convert order.Order to GraphQL Order type
	var orderedProducts []*OrderedProduct
	for _, p := range orderResult.Products {
		orderedProducts = append(orderedProducts, &OrderedProduct{
			ID:       p.ProductID,
			Quantity: p.Quantity,
		})
	}

	return &Order{
		ID:          orderResult.ID,
		CreatedAt:   orderResult.CreatedAt,
		TotalAmount: orderResult.Total,
		Products:    orderedProducts,
	}, nil
}
