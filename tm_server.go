// tm_server.go
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

type tmServer struct{
	transaction.UnimplementedTransactionManagerServer
}

var port_number [3]string = [3]string{"50052", "50053", "50054"}

func (s *tmServer) Prepare(ctx context.Context, req *transaction.PrepareRequest) (*transaction.PrepareResponse, error) {
	// Perform any necessary checks and validations here
	fmt.Printf("Received Prepare request for transaction ID: %s\n", req.TransactionId)
	for index, port_num := range port_number {
		fmt.Printf("Send Prepare to service # %d, which is on port %s\n", index, port_num)
		conn, err := grpc.Dial("localhost:"+port_num, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("could not connect: %v", err)
		}
		defer conn.Close()

		client := transaction.NewTransactionManagerClient(conn)

		prepareResp, err := client.Prepare(context.Background(), &transaction.PrepareRequest{TransactionId: port_num})
		if err != nil {
			log.Fatalf("Prepare failed for service %d: %v", index, err)
		}
		if !prepareResp.Prepared {
			log.Printf("Prepare phase failed for service %d. Aborting...", index)
			return &transaction.PrepareResponse{Prepared: false}, nil
		}
		fmt.Printf("Prepare phase succeeded for service %d. Proceeding to Commit Phase.\n", index)
	}
	return &transaction.PrepareResponse{Prepared: true}, nil
}

func (s *tmServer) Commit(ctx context.Context, req *transaction.CommitRequest) (*transaction.CommitResponse, error) {
	// Perform the Commit phase here
	fmt.Printf("Received Commit request for transaction ID: %s\n", req.TransactionId)
	for index, port_num := range port_number {
		fmt.Printf("Send Commit to service # %d, which is on port %s\n", index, port_num)
		conn, err := grpc.Dial("localhost:"+port_num, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("could not connect: %v", err)
		}
		defer conn.Close()

		client := transaction.NewTransactionManagerClient(conn)

		commitResp, err := client.Commit(context.Background(), &transaction.CommitRequest{TransactionId: port_num})
		if err != nil {
			log.Fatalf("Commit failed for service %d: %v", index, err)
		}
		if !commitResp.Committed {
			log.Printf("Commit phase failed for service %d. Initiating Abort...", index)
			client.Abort(context.Background(), &transaction.AbortRequest{TransactionId: port_num})
			return &transaction.CommitResponse{Committed: false}, nil
		}
		fmt.Printf("Commit phase succeeded for service %d. Transaction completed successfully.\n", index)
	}
	return &transaction.CommitResponse{Committed: true}, nil
}

func (s *tmServer) Abort(ctx context.Context, req *transaction.AbortRequest) (*transaction.AbortResponse, error) {
	// Perform the Abort phase here
	fmt.Printf("Received Abort request for transaction ID: %s\n", req.TransactionId)
	for index, port_num := range port_number {
		fmt.Printf("Send Prepare to service # %d, which is on port %s\n", index, port_num)
		conn, err := grpc.Dial("localhost:"+port_num, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("could not connect: %v", err)
		}
		defer conn.Close()

		client := transaction.NewTransactionManagerClient(conn)

		client.Abort(context.Background(), &transaction.AbortRequest{TransactionId: port_num})
		fmt.Printf("Abord phase succeeded for service %d. Proceeding to Commit Phase.\n", index)
	}
	return &transaction.AbortResponse{Aborted: true}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	transaction.RegisterTransactionManagerServer(s, &tmServer{})
	reflection.Register(s)
	fmt.Println("Transaction Manager Server started...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func get_service_client(port string) () {
	
}
