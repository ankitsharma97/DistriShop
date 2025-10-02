package account

import (
	"context"
	"fmt"
	"net"

	pb "microservice/account/pb" // generated via: protoc --go_out=. --go-grpc_out=. account/account.proto

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// grpcServer implements the generated gRPC AccountServiceServer.
type grpcServer struct {
	pb.UnimplementedAccountServiceServer
	service Service
}

// ListenGRPC starts a gRPC server for the Account service.
func ListenGRPC(service Service, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	s := grpc.NewServer()
	pb.RegisterAccountServiceServer(s, &grpcServer{service: service})
	reflection.Register(s)
	return s.Serve(lis)
}

func (s *grpcServer) PostAccount(ctx context.Context, req *pb.PostAccountRequest) (*pb.PostAccountResponse, error) {
	acc, err := s.service.PostAccount(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	return &pb.PostAccountResponse{Account: &pb.Account{Id: acc.ID, Name: acc.Name}}, nil
}

func (s *grpcServer) GetAccount(ctx context.Context, req *pb.GetAccountRequest) (*pb.GetAccountResponse, error) {
	acc, err := s.service.GetAccount(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetAccountResponse{Account: &pb.Account{Id: acc.ID, Name: acc.Name}}, nil
}

func (s *grpcServer) GetAccounts(ctx context.Context, req *pb.GetAccountsRequest) (*pb.GetAccountsResponse, error) {
	accounts, err := s.service.GetAccounts(ctx, int(req.Skip), int(req.Take))
	if err != nil {
		return nil, err
	}
	resp := make([]*pb.Account, 0, len(accounts))
	for _, a := range accounts {
		resp = append(resp, &pb.Account{Id: a.ID, Name: a.Name})
	}
	return &pb.GetAccountsResponse{Accounts: resp}, nil
}
