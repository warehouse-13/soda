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

	// This for emulates a long-running process. It will only exit when the whole
	// program exits, which in this case will only happen when a SIGINT (ctrl-c) is
	// sent. This is fine for the purpose of this repro, as it is not meant to be
	// exemplary code.
	for {
		// Here we are opening the connection, creating the client, and calling the service
		// all from within the long running process.
		// This means that the defer conn.Close() is ineffective and the bug will not be
		// solved.
		// A defer is called just before a function exits. The function in this case
		// is main, and it does not exit until this for loop is finished. As there is
		// no way for this loop to finish (it would require a ctrl-c to cancel
		// it and exit the program), the defer can never be called and the connections
		// on the service will stay open.
		conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			fmt.Printf("could not open connection to server at %s. err: %s\n", address, err)
		}
		defer conn.Close()

		client := proto.NewSodaServiceClient(conn)

		resp, err := client.RandomNumber(context.Background(), &proto.RandomNumberRequest{})
		if err != nil {
			fmt.Printf("could not make call %s\n", err)
			continue
		}

		fmt.Println(resp.Result)
	}
}
