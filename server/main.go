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
	_ "net/http/pprof"
)

func main() {
	var (
		servicePort string
		pprofPort   string
	)

	flag.StringVar(&servicePort, "port", "1430", "port to start service on")
	flag.StringVar(&pprofPort, "pprof", "1431", "port to start pprof on")
	flag.Parse()

	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, syscall.SIGTERM, syscall.SIGINT)
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	fmt.Println("starting Soda Service")
	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
		}()

		if err := runSoda(ctx, servicePort); err != nil {
			fmt.Printf("failed to start service %s", err)
		}
	}()
	fmt.Println("running on port " + servicePort)

	fmt.Println("starting pprof server")
	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
		}()

		if err := runPProf(ctx, pprofPort); err != nil {
			fmt.Printf("failed to start pprof: %s", err)
		}
	}()
	fmt.Println("running on port " + pprofPort)

	<-stopChan
	cancel()
	wg.Wait()
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

func runSoda(ctx context.Context, port string) error {
	s := newServer()
	grpcServer := grpc.NewServer()
	defer func() {
		fmt.Println("service stopped, finishing last request")
		grpcServer.GracefulStop()
	}()

	proto.RegisterSodaServiceServer(grpcServer, &s)

	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	go func() {
		grpcServer.Serve(l)
	}()

	<-ctx.Done()

	return nil
}

func runPProf(ctx context.Context, port string) error {
	srv := &http.Server{
		Addr:    "localhost:" + port,
		Handler: http.DefaultServeMux,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("could not start pprof %s", err)
		}
	}()

	<-ctx.Done()
	fmt.Println("pprof server stopped")

	if err := srv.Shutdown(context.Background()); err != nil {
		fmt.Printf("pprof shutdown failed %s", err)
	}

	return nil
}
