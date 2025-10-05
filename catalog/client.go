package catalog

import (
	"context"
	pb "microservice/catalog/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.CatalogServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.Dial(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	service := pb.NewCatalogServiceClient(conn)
	return &Client{conn: conn, service: service}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) PostProduct(ctx context.Context, name, description string, price float64) (*Product, error) {
	req := &pb.PostProductRequest{Name: name, Description: description, Price: price}
	resp, err := c.service.PostProduct(ctx, req)
	if err != nil {
		return nil, err
	}
	return &Product{ID: resp.Product.Id, Name: resp.Product.Name, Description: resp.Product.Description, Price: resp.Product.Price}, nil
}

func (c *Client) GetProduct(ctx context.Context, id string) (*Product, error) {
	req := &pb.GetProductRequest{Id: id}
	resp, err := c.service.GetProduct(ctx, req)
	if err != nil {
		return nil, err
	}
	return &Product{ID: resp.Product.Id, Name: resp.Product.Name, Description: resp.Product.Description, Price: resp.Product.Price}, nil
}

func (c *Client) GetProducts(ctx context.Context, ids []string, query string, skip int, take int) ([]*Product, error) {
	req := &pb.GetProductsRequest{Ids: ids, Query: query, Skip: uint64(skip), Take: uint64(take)}
	resp, err := c.service.GetProducts(ctx, req)
	if err != nil {
		return nil, err
	}
	return convertProducts(resp.Products), nil
}

func convertProducts(products []*pb.Product) []*Product {
	var result []*Product
	for _, p := range products {
		result = append(result, &Product{ID: p.Id, Name: p.Name, Description: p.Description, Price: p.Price})
	}
	return result
}
