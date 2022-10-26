package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.io/warehouse-13/soda/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	var address string
	flag.StringVar(&address, "address", "", "server address and port")
	flag.Parse()

	if address == "" {
		fmt.Println("required: --address")
		os.Exit(1)
	}

	for {
		resp, err := call(address)
		if err != nil {
			fmt.Printf("could not make call %s\n", err)
			continue
		}

		fmt.Println(resp.Result)
	}
}

func call(address string) (*proto.RandomNumberResponse, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("could not open connection to server at %s. err: %s\n", address, err)
	}
	// defer conn.Close()

	client := proto.NewSodaServiceClient(conn)

	return client.RandomNumber(context.Background(), &proto.RandomNumberRequest{})
}
