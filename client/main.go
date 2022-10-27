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
		resp, err := call(address)
		if err != nil {
			fmt.Printf("could not make call %s\n", err)
			continue
		}

		fmt.Println(resp.Result)
	}
}

// call will create a new conection to the service, create a client with
// that connection, then calls an endpoint on the service.
// What we are emulating here is a long running process which repeats this process
// on every reconcile/action. If the connection were established once before the
// long running process began (ie before the loop above), then we would not encounter
// the error.
// Similarly, if this connection/create/call process was not abstracted into a func,
// and was instead called directly in the for loop, the fix of adding a defer conn.Close()
// would not have the desired effect. A defer is only called upon a function's exit.
// If this was all happening directly inside the for loop, then the defer would never
// be called because main has not extited, and the connections would remain open.
//
// Checkout to branch no-call to see what I mean here.
func call(address string) (*proto.RandomNumberResponse, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("could not open connection to server at %s. err: %s\n", address, err)
	}
	// defer conn.Close()

	client := proto.NewSodaServiceClient(conn)

	return client.RandomNumber(context.Background(), &proto.RandomNumberRequest{})
}
