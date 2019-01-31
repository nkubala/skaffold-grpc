package main

import (
	"context"
	"fmt"
	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/event/proto"
	// empty "github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	// "io"
	"os"
	// "time"
)

var (
	addr = ":50051"
	// request = &event.Request{Name: "Request"}
	// response = new(event.Response)
	// stateResponse = &event.StateResponse{}
)

func main() {

	// client, err := rpc.Dial("tcp", addr)
	// client, err := rpc.DialHTTP("tcp", addr)
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		fmt.Printf("error opening connection: %s\n", err.Error())
		os.Exit(1)
	}
	defer conn.Close()
	client := proto.NewSkaffoldServiceClient(conn)

	// ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// r, err := client.GetState(ctx, &empty.Empty{})
	// if err != nil {
	// 	fmt.Printf("error retrieving state: %v\n", err)
	// }
	// fmt.Printf("response from server: %+v\n", r)

	stream, err := client.EventLog(ctx)
	if err != nil {
		fmt.Printf("error retrieving event log: %v\n", err)
	}
	for {
		entry, err := stream.Recv()
		if err != nil {
			fmt.Printf("error receiving message from stream: %v\n", err)
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
