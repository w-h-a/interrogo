package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/tmc/langchaingo/llms/googleai"
	"github.com/tmc/langchaingo/llms/googleai/vertex"
	toolprovider "github.com/w-h-a/interrogo/internal/client/tool_provider"
	"github.com/w-h-a/interrogo/internal/client/tool_provider/mcp"
	chathttphandler "github.com/w-h-a/interrogo/internal/handler/http/chat"
	"github.com/w-h-a/interrogo/internal/server"
	httpserver "github.com/w-h-a/interrogo/internal/server/http"
	"github.com/w-h-a/interrogo/internal/service/agent"
)

func main() {
	// config and instrument
	ctx := context.Background()

	project := os.Getenv("GCP_PROJECT_ID")
	if len(project) == 0 {
		log.Fatal("GCP_PROJECT_ID env var is required")
	}

	// stop chans
	stopChannels := map[string]chan struct{}{}

	// clients
	mcpToolProvider := mcp.NewToolProvider(
		toolprovider.WithLocation("http://localhost:8081/mcp"),
	)

	log.Println("connecting to Tool Provider...")

	if err := mcpToolProvider.Start(ctx); err != nil {
		log.Fatalf("Failed to connect to Tool Provider: %v", err)
	}

	v, err := vertex.New(
		ctx,
		googleai.WithCloudProject(project),
		googleai.WithCloudLocation("us-central1"),
		googleai.WithDefaultModel("gemini-2.0-flash-001"),
	)
	if err != nil {
		log.Fatalf("Failed to initialize Vertex AI: %v", err)
	}

	// services
	a := agent.New(v, mcpToolProvider, "You are a helpful assistant. Use the 'list_records' tool to find information. Refuse to delete data.")
	stopChannels["agent"] = make(chan struct{})

	// servers
	httpServer, err := InitHttpServer(":8080", a)
	if err != nil {
		log.Fatalf("Failed to initialize http server: %v", err)
	}
	stopChannels["httpServer"] = make(chan struct{})

	// wait group and chans for graceful shutdown
	var wg sync.WaitGroup
	errCh := make(chan error, len(stopChannels))
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// run
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("running agent")
		errCh <- a.Run(stopChannels["agent"])
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("running http server")
		errCh <- httpServer.Run(stopChannels["httpServer"])
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

func InitHttpServer(httpAddr string, a *agent.Agent) (server.Server, error) {
	srv := httpserver.NewServer(
		server.WithAddress(httpAddr),
		server.WithName("my-http"),
		server.WithVersion("0.1.0-alpha.0"),
	)

	router := mux.NewRouter()

	chatHandler := chathttphandler.New(a)

	router.HandleFunc("/chat", chatHandler.Chat).Methods(http.MethodPost)

	if err := srv.Handle(router); err != nil {
		return nil, fmt.Errorf("failed to attach root handler: %w", err)
	}

	return srv, nil
}
