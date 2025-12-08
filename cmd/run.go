package cmd

import (
	"context"
	"fmt"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"
	"github.com/tmc/langchaingo/llms/googleai/vertex"
	"github.com/urfave/cli/v2"
	"github.com/w-h-a/interrogo/api/config/v1alpha1"
	"github.com/w-h-a/interrogo/internal/client/interrogator"
	"github.com/w-h-a/interrogo/internal/client/interrogator/http"
	"github.com/w-h-a/interrogo/internal/config"
	"github.com/w-h-a/interrogo/internal/service/judge"
)

func Judge(c *cli.Context) error {
	// context
	ctx := c.Context

	// config
	configPath := c.String("config")
	if len(configPath) == 0 {
		return fmt.Errorf("--config is required")
	}
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("config error: %w", err)
	}

	// clients
	// TODO: check config for init options
	httpInterrogator := http.NewInterrogator(
		interrogator.WithTarget(cfg.Target.URL),
	)

	evaluatorClient, err := InitEvaluator(ctx, cfg)
	if err != nil {
		return err
	}

	// service
	j := judge.New(evaluatorClient, httpInterrogator)

	// do it
	fmt.Println("Attacking agent via", cfg.Evaluator.Provider, "...")

	results := j.Judge(ctx, cfg.Evaluator.Policy)

	// report
	for i, result := range results {
		if len(result.Error) > 0 {
			fmt.Printf("[%d] ⚠️  Execution Error: %s\n", i+1, result.Error)
			continue
		}

		icon := "✅"
		if !result.Passed {
			icon = "🚫"
		}

		fmt.Printf("[%d] %s %s\n", i+1, icon, result.Reasoning)
	}

	return nil
}

func InitEvaluator(ctx context.Context, cfg *v1alpha1.Config) (llms.Model, error) {
	var model llms.Model
	var err error

	switch cfg.Evaluator.Provider {
	case "vertex":
		model, err = vertex.New(
			ctx,
			googleai.WithCloudProject(cfg.Evaluator.Params["project_id"]),
			googleai.WithCloudLocation("us-central1"),
			googleai.WithDefaultModel("gemini-2.0-flash-001"),
		)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", cfg.Evaluator.Provider)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to init %s llm: %w", cfg.Evaluator.Provider, err)
	}

	return model, nil
}
