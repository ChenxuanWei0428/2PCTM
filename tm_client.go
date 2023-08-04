// tm_client.go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	transaction "proto_gen"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()


	client := transaction.NewTransactionManagerClient(conn)

	// Example 2PC Transaction
	transactionID := "sample-transaction-id"

	// Prepare Phase
	prepareResp, err := client.Prepare(context.Background(), &transaction.PrepareRequest{TransactionId: transactionID})
	if err != nil {
		log.Fatalf("Prepare failed: %v", err)
	}
	if !prepareResp.Prepared {
		log.Println("Prepare phase failed. Aborting...")
		return
	}
	fmt.Println("Prepare phase succeeded. Proceeding to Commit Phase.")

	// Wait for a short duration to simulate some processing time
	time.Sleep(2 * time.Second)

	// Commit Phase
	commitResp, err := client.Commit(context.Background(), &transaction.CommitRequest{TransactionId: transactionID})
	if err != nil {
		log.Fatalf("Commit failed: %v", err)
	}
	if !commitResp.Committed {
		log.Println("Commit phase failed. Initiating Abort...")
		client.Abort(context.Background(), &transaction.AbortRequest{TransactionId: transactionID})
		return
	}
	fmt.Println("Commit phase succeeded. Transaction completed successfully.")
}
