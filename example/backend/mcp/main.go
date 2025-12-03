package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	recordmcphandler "github.com/w-h-a/interrogo/internal/handler/mcp/record"
	"github.com/w-h-a/interrogo/internal/server"
	mcpserver "github.com/w-h-a/interrogo/internal/server/mcp"
)

func main() {
	// config and instrument

	// stop chans
	stopChannels := map[string]chan struct{}{}

	// create clients

	// create services

	// create servers
	mcpSrv, err := InitMcpServer(":8081")
	if err != nil {
		panic(err)
	}
	stopChannels["mcpServer"] = make(chan struct{})

	// wait group and chans for graceful shutdown
	var wg sync.WaitGroup
	errCh := make(chan error, len(stopChannels))
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// run
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("running mcp server")
		errCh <- mcpSrv.Run(stopChannels["mcpServer"])
	}()

	// block until shutdown
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-sigChan:
		for _, stop := range stopChannels {
			close(stop)
		}
	}

	wg.Wait()

	log.Println("successfully shutdown")
}

func InitMcpServer(mcpAddr string) (server.Server, error) {
	srv := mcpserver.NewServer(
		server.WithAddress(mcpAddr),
		server.WithName("my-mcp"),
		server.WithVersion("0.1.0-alpha.0"),
	)

	recordHandler := recordmcphandler.New()

	if err := srv.Handle(recordHandler.ListRecordsTool()); err != nil {
		return nil, fmt.Errorf("failed to register list tool: %w", err)
	}

	return srv, nil
}
