package catalog

import (
	"context"
	"fmt"
	"net"

	pb "microservice/catalog/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	pb.UnimplementedCatalogServiceServer
	service Service
}

// ListenGRPC starts a gRPC server for the Catalog service.
func ListenGRPC(service Service, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	s := grpc.NewServer()
	pb.RegisterCatalogServiceServer(s, &grpcServer{service: service})
	reflection.Register(s)
	return s.Serve(lis)
}

func (s *grpcServer) PostProduct(ctx context.Context, req *pb.PostProductRequest) (*pb.PostProductResponse, error) {
	product, err := s.service.PostProduct(ctx, req.Name, req.Description, req.Price)
	if err != nil {
		return nil, err
	}
	return &pb.PostProductResponse{Product: &pb.Product{Id: product.ID, Name: product.Name, Description: product.Description, Price: product.Price}}, nil
}

func (s *grpcServer) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	product, err := s.service.GetProduct(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetProductResponse{Product: &pb.Product{Id: product.ID, Name: product.Name, Description: product.Description, Price: product.Price}}, nil
}

func (s *grpcServer) GetProducts(ctx context.Context, req *pb.GetProductsRequest) (*pb.GetProductsResponse, error) {
	var products []*Product
	var err error

	if req.Query != "" { // search
		products, err = s.service.SearchProducts(ctx, req.Query, int(req.Skip), int(req.Take))
	} else if len(req.Ids) > 0 { // by IDs
		products, err = s.service.GetProductByIDs(ctx, req.Ids)
	} else { // pagination only
		products, err = s.service.GetProducts(ctx, int(req.Skip), int(req.Take))
	}
	if err != nil {
		return nil, err
	}

	resp := make([]*pb.Product, 0, len(products))
	for _, p := range products {
		resp = append(resp, &pb.Product{Id: p.ID, Name: p.Name, Description: p.Description, Price: p.Price})
	}
	return &pb.GetProductsResponse{Products: resp}, nil
}
