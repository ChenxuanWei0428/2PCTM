package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	transaction "proto_gen"
)

type serviceAServer struct{
	transaction.UnimplementedTransactionManagerServer
}

func (s *serviceAServer) Prepare(ctx context.Context, req *transaction.PrepareRequest) (*transaction.PrepareResponse, error) {
	fmt.Printf("Service C received Prepare request for transaction ID: %s\n", req.TransactionId)
	return &transaction.PrepareResponse{Prepared: true}, nil
}

func (s *serviceAServer) Commit(ctx context.Context, req *transaction.CommitRequest) (*transaction.CommitResponse, error) {
	fmt.Printf("Service C received Commit request for transaction ID: %s\n", req.TransactionId)
	return &transaction.CommitResponse{Committed: true}, nil
}

func (s *serviceAServer) Abort(ctx context.Context, req *transaction.AbortRequest) (*transaction.AbortResponse, error) {
	fmt.Printf("Service C received Abort request for transaction ID: %s\n", req.TransactionId)
	return &transaction.AbortResponse{Aborted: true}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50054")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	transaction.RegisterTransactionManagerServer(s, &serviceAServer{})
	reflection.Register(s)
	fmt.Println("Service C Server started...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
