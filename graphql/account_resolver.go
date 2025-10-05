package main

import (
	"context"
	"time"
)

type accountResolver struct {
	server *Server
}

func (r *accountResolver) Orders(ctx context.Context, obj *Account) ([]*Order, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	orderList, err := r.server.orderClient.GetOrderForAccount(ctx, obj.ID)
	if err != nil {
		return nil, err
	}

	var orders []*Order
	for _, o := range orderList {
		var products []*OrderedProduct
		for _, p := range o.Products {
			products = append(products, &OrderedProduct{
				ID:       p.ProductID,
				Quantity: p.Quantity,
				// Note: order.OrderedProduct doesn't have Name, Description, Price
				// These would need to be fetched separately from catalog service
			})
		}
		orders = append(orders, &Order{
			ID:          o.ID,
			TotalAmount: o.Total,
			CreatedAt:   o.CreatedAt,
			Products:    products,
		})
	}

	return orders, nil
}
