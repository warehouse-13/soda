package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.io/warehouse-13/soda/proto"
	"google.golang.org/grpc"
)

func main() {
	var port string

	flag.StringVar(&port, "port", "1430", "port to start server on")
	flag.Parse()

	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Printf("Failed to listen on localhost:%s, %s", port, err)
		os.Exit(1)
	}

	s := newServer()

	grpcServer := grpc.NewServer()
	proto.RegisterSodaServiceServer(grpcServer, &s)

	errChan := make(chan error)
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, syscall.SIGTERM, syscall.SIGINT)

	fmt.Println("starting Soda Service")
	go func() {
		if err := grpcServer.Serve(l); err != nil {
			fmt.Printf("failed to start service: %s", err)
			errChan <- err
		}
	}()

	defer func() {
		grpcServer.GracefulStop()
	}()

	fmt.Println("running on port " + port)

	select {
	case err := <-errChan:
		fmt.Printf("server received fatal error: %s", err)
		os.Exit(1)
	case <-stopChan:
		fmt.Println("service stopped, finishing last request")
	}
}

type server struct {
	proto.UnimplementedSodaServiceServer
}

func newServer() server {
	return server{}
}

func (s server) RandomNumber(ctx context.Context, req *proto.RandomNumberRequest) (*proto.RandomNumberResponse, error) {
	rand.Seed(time.Now().UnixNano())
	// artificially slow things down a bit
	time.Sleep(time.Millisecond * 10)
	return &proto.RandomNumberResponse{
		Result: rand.Uint32(),
	}, nil
}
