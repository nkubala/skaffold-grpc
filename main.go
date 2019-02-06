package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/event/proto"

	"google.golang.org/grpc"
)

var (
	addr        = ":50051"
	MAX_ERRORS  = 5
	MAX_RETRIES = 10
)

func main() {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		fmt.Printf("error opening connection: %s\n", err.Error())
		os.Exit(1)
	}
	defer conn.Close()
	client := proto.NewSkaffoldServiceClient(conn)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// r, err := client.GetState(ctx, &empty.Empty{})
	// if err != nil {
	// 	fmt.Printf("error retrieving state: %v\n", err)
	// }
	// fmt.Printf("response from server: %+v\n", r)
	retries := 0
	var stream proto.SkaffoldService_EventLogClient
	for {
		stream, err = client.EventLog(ctx)
		if err == nil {
			break
		} else if retries < MAX_RETRIES {
			retries = retries + 1
			fmt.Println("waiting for connection...")
			time.Sleep(3 * time.Second)
			continue
		}
		fmt.Printf("error retrieving event log: %v\n", err)
		os.Exit(1)
	}

	errors := 0
	for {
		entry, err := stream.Recv()
		if err != nil {
			errors = errors + 1
			fmt.Printf("[%d] error receiving message from stream: %v\n", errors, err)
			if errors == MAX_ERRORS {
				fmt.Printf("%d errors encountered: quitting", MAX_ERRORS)
				os.Exit(1)
			}
			time.Sleep(1 * time.Second)
		} else {
			fmt.Printf("%+v\n", entry)
		}
	}
	// waitc := make(chan struct{})
	// go func() {
	// 	for {
	// 		entry, err := stream.Recv()
	// 		if err == io.EOF {
	// 			close(waitc)
	// 			return
	// 			// break
	// 		}
	// 		if err != nil {
	// 			fmt.Printf("error receiving message from stream: %v\n", err)
	// 		}
	// 		fmt.Printf("%+v\n", entry)
	// 	}
	// }()
	// // stream.CloseSend()
	// <-waitc

}
