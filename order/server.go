package order

import (
	"context"
	"fmt"
	"net"

	"microservice/account"
	"microservice/catalog"
	pb "microservice/order/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	pb.UnimplementedOrderServiceServer
	service       Service
	accountClient *account.Client
	catalogClient *catalog.Client
}

// ListenGRPC starts a gRPC server for the Order service.
func ListenGRPC(service Service, accountURL, catalogURL string, port int) error {
	accountClient, err := account.NewClient(accountURL)
	if err != nil {
		return err
	}
	catalogClient, err := catalog.NewClient(catalogURL)
	if err != nil {
		accountClient.Close()
		return err
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		accountClient.Close()
		catalogClient.Close()
		return err
	}
	s := grpc.NewServer()
	pb.RegisterOrderServiceServer(s, &grpcServer{service: service, accountClient: accountClient, catalogClient: catalogClient})
	reflection.Register(s)
	return s.Serve(lis)
}

func (s *grpcServer) PostOrder(ctx context.Context, req *pb.PostOrderRequest) (*pb.PostOrderResponse, error) {
	// Verify account exists
	_, err := s.accountClient.GetAccount(ctx, req.AccountId)
	if err != nil {
		return nil, err
	}

	var products []*OrderedProduct
	var total float64 = 0.0

	for _, item := range req.Products {
		product, err := s.catalogClient.GetProduct(ctx, item.ProductId)
		if err != nil {
			return nil, err
		}

		orderedProduct := &OrderedProduct{
			ProductID: product.ID,
			Quantity:  int(item.Quantity),
		}
		products = append(products, orderedProduct)
		total += product.Price * float64(item.Quantity)
	}

	order, err := s.service.PostOrder(ctx, req.AccountId, products)
	if err != nil {
		return nil, err
	}

	// Update the total calculated from catalog service
	order.Total = total

	// Get product details for response
	productIDs := []string{}
	for _, p := range order.Products {
		productIDs = append(productIDs, p.ProductID)
	}

	catalogProducts, err := s.catalogClient.GetProducts(ctx, productIDs, "", 0, 0)
	if err != nil {
		return nil, err
	}

	orderProto := &pb.Order{
		Id:        order.ID,
		AccountId: order.AccountID,
		Total:     order.Total,
		Products:  []*pb.OrderedProduct{},
	}
	orderProto.CreatedAt, _ = order.CreatedAt.MarshalBinary()

	for _, p := range order.Products {
		for _, catalogProduct := range catalogProducts {
			if catalogProduct.ID == p.ProductID {
				orderProto.Products = append(orderProto.Products, &pb.OrderedProduct{
					Id:          catalogProduct.ID,
					Name:        catalogProduct.Name,
					Description: catalogProduct.Description,
					Price:       catalogProduct.Price,
					Quantity:    uint32(p.Quantity),
				})
				break
			}
		}
	}

	return &pb.PostOrderResponse{Order: orderProto}, nil
}

func (s *grpcServer) GetOrderForAccount(ctx context.Context, req *pb.GetOrderForAccountRequest) (*pb.GetOrderForAccountResponse, error) {
	accountOrders, err := s.service.GetOrdersForAccount(ctx, req.AccountId)
	if err != nil {
		return nil, err
	}

	productIDs := map[string]bool{}
	for _, order := range accountOrders {
		for _, product := range order.Products {
			productIDs[product.ProductID] = true
		}
	}
	ids := []string{}
	for id := range productIDs {
		ids = append(ids, id)
	}

	catalogProducts, err := s.catalogClient.GetProducts(ctx, ids, "", 0, 0)
	if err != nil {
		return nil, err
	}

	productMap := map[string]*catalog.Product{}
	for _, p := range catalogProducts {
		productMap[p.ID] = p
	}

	var orders []*pb.Order
	for _, order := range accountOrders {
		orderProto := &pb.Order{
			Id:        order.ID,
			AccountId: order.AccountID,
			Total:     order.Total,
			Products:  []*pb.OrderedProduct{},
		}
		orderProto.CreatedAt, _ = order.CreatedAt.MarshalBinary()

		for _, orderedProduct := range order.Products {
			if catalogProduct, exists := productMap[orderedProduct.ProductID]; exists {
				orderProto.Products = append(orderProto.Products, &pb.OrderedProduct{
					Id:          catalogProduct.ID,
					Name:        catalogProduct.Name,
					Description: catalogProduct.Description,
					Price:       catalogProduct.Price,
					Quantity:    uint32(orderedProduct.Quantity),
				})
			}
		}
		orders = append(orders, orderProto)
	}

	return &pb.GetOrderForAccountResponse{Orders: orders}, nil
}

func (s *grpcServer) GetOrders(ctx context.Context, req *pb.GetOrdersRequest) (*pb.GetOrdersResponse, error) {
	// This method could return all orders or orders based on some criteria
	// For now, return the order from the request
	return &pb.GetOrdersResponse{Orders: []*pb.Order{req.Order}}, nil
}
