package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.io/warehouse-13/soda/proto"
	"google.golang.org/grpc"

	"net/http"

	// This is all you need to add a debug endpoint on whatever http
	// server you have running.
	// See here for more details https://pkg.go.dev/net/http/pprof
	_ "net/http/pprof"
)

// Note that none of this is production code. It is all a quick hack to demonstrate
// the bug and pprof.
func main() {
	var (
		servicePort string
		pprofPort   string
	)

	flag.StringVar(&servicePort, "port", "1430", "port to start service on")
	flag.StringVar(&pprofPort, "pprof", "1431", "port to start pprof on")
	flag.Parse()

	// Set up a channel to watch and wait for user stop signals.
	// Once such a signal is received, running services are shut down.
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, syscall.SIGTERM, syscall.SIGINT)

	// Create a context with a cancel func. The cancel will be called later to
	// trigger shutdown of services, via the ctx that each will receive.
	ctx, cancel := context.WithCancel(context.Background())

	// Create a wait group to synchronize jobs and ensure all are complete before
	// exiting.
	wg := &sync.WaitGroup{}

	// Add a count to the wait group to represent the gRPC service job
	wg.Add(1)
	fmt.Println("starting Soda Service")
	go func() {
		defer func() {
			// When the gRPC service exits the job is marked as done and removed from
			// the wait group.
			wg.Done()
		}()

		// Start the gRPC service
		if err := runSoda(ctx, servicePort); err != nil {
			fmt.Printf("failed to start service %s\n", err)
		}
	}()
	fmt.Println("gRPC service running on port " + servicePort)

	// Add a count to the wait group to represent the pprof http job
	wg.Add(1)
	fmt.Println("starting pprof server")
	go func() {
		defer func() {
			// When the gRPC service exits the job is marked as done and removed from
			// the wait group.
			wg.Done()
		}()

		// Start the http service
		// If we already had an http server in this app we would not need to create
		// a new one. The import of the 'net/http/pprof' package simply adds the endpoint
		// to any existing http service.
		// In our case we don't have an http service already, just a gRPC one, so we
		// start a dummy http server which will only serve the /debug/pprof endpoint.
		if err := runPProf(ctx, pprofPort); err != nil {
			fmt.Printf("failed to start pprof: %s\n", err)
		}
	}()
	fmt.Println("pprof running on port " + pprofPort)

	// The program will hold at this line until a stop signal is received on that
	// channel.
	<-stopChan

	// A stop has been received, so we send a cancel to the services via their contexts.
	cancel()

	// We wait for each service to be done. Once all are reported the program can exit.
	wg.Wait()
}

func runSoda(ctx context.Context, port string) error {
	// Create a new instance of the dummy SodaService implementation
	s := newServer()

	// Create a new gRPC server
	grpcServer := grpc.NewServer()
	// Set up a graceful shutdown to happen on function exit
	defer func() {
		grpcServer.GracefulStop()
		fmt.Println("grpc service stopped")
	}()

	// Register the dummy service on that server
	proto.RegisterSodaServiceServer(grpcServer, &s)

	// Add a listener to it on the given port
	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	go func() {
		// Start the service
		if err := grpcServer.Serve(l); err != nil {
			fmt.Printf("could not start grpc service %s\n", err)
		}
	}()

	// Hold here until a cancel is sent on the context
	<-ctx.Done()

	fmt.Println("cancel signal received, stopping grpc service")

	// Now the cancel has been received we can move to exit the function, calling
	// the deferred shutdown that we set up at the start.
	return nil
}

func runPProf(ctx context.Context, port string) error {
	// Create a new http server
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: http.DefaultServeMux,
	}

	go func() {
		// Start the server on the given port
		// The debug endpoint will be localhost:1431/debug/pprof/
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("could not start pprof http service %s\n", err)
		}
	}()

	// Hold here until a cancel is sent on the context
	<-ctx.Done()

	fmt.Println("cancel signal received, stopping pprof server")

	// Now the cancel has been received shut down the server
	if err := srv.Shutdown(context.Background()); err != nil {
		fmt.Printf("pprof shutdown failed %s\n", err)
	}

	fmt.Println("pprof server stopped")

	return nil
}

// server is an implementation of the dummy gRPC service I created for this repro
type server struct {
	proto.UnimplementedSodaServiceServer
}

func newServer() server {
	return server{}
}

// RandomNumber is a dummy func which could return any old thing. It it there to
// show that the server is responding to client requests during the demo.
func (s server) RandomNumber(ctx context.Context, req *proto.RandomNumberRequest) (*proto.RandomNumberResponse, error) {
	rand.Seed(time.Now().UnixNano())
	// Artificially slow things down a bit, otherwise connections get created so quickly
	// you barely have time to appreciate what the bug is doing.
	time.Sleep(time.Millisecond * 10)
	return &proto.RandomNumberResponse{
		Result: rand.Uint32(),
	}, nil
}
