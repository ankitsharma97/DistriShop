package order

import (
	"context"
	"time"

	pb "microservice/order/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.OrderServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.Dial(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	service := pb.NewOrderServiceClient(conn)
	return &Client{conn: conn, service: service}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) PostOrder(ctx context.Context, accountID string, products []*OrderedProduct) (*Order, error) {
	protoProducts := []*pb.PostOrderRequest_OrderProduct{}
	for _, p := range products {
		protoProducts = append(protoProducts, &pb.PostOrderRequest_OrderProduct{
			ProductId: p.ProductID,
			Quantity:  uint32(p.Quantity),
		})
	}

	resp, err := c.service.PostOrder(ctx, &pb.PostOrderRequest{
		AccountId: accountID,
		Products:  protoProducts,
	})
	if err != nil {
		return nil, err
	}

	newOrder := resp.Order
	newOrderCreatedAt := time.Time{}
	newOrderCreatedAt.UnmarshalBinary(resp.Order.CreatedAt)

	// Convert protobuf products back to OrderedProducts
	var orderProducts []*OrderedProduct
	for _, p := range newOrder.Products {
		orderProducts = append(orderProducts, &OrderedProduct{
			ProductID: p.Id,
			Quantity:  int(p.Quantity),
		})
	}

	return &Order{
		ID:        resp.Order.Id,
		AccountID: newOrder.AccountId,
		Products:  orderProducts,
		Total:     newOrder.Total,
		CreatedAt: newOrderCreatedAt,
	}, nil
}

func (c *Client) GetOrderForAccount(ctx context.Context, accountID string) ([]*Order, error) {
	resp, err := c.service.GetOrderForAccount(ctx, &pb.GetOrderForAccountRequest{
		AccountId: accountID,
	})
	if err != nil {
		return nil, err
	}

	orders := []*Order{}
	for _, o := range resp.Orders {
		newOrder := &Order{
			ID:        o.Id,
			AccountID: o.AccountId,
			Total:     o.Total,
		}
		newOrderCreatedAt := time.Time{}
		newOrderCreatedAt.UnmarshalBinary(o.CreatedAt)
		newOrder.CreatedAt = newOrderCreatedAt

		var products []*OrderedProduct
		for _, p := range o.Products {
			products = append(products, &OrderedProduct{
				ProductID: p.Id,
				Quantity:  int(p.Quantity),
			})
		}
		newOrder.Products = products
		orders = append(orders, newOrder)
	}

	return orders, nil
}
