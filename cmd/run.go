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
	m, err := InitModel(ctx, cfg.Evaluator)
	if err != nil {
		return err
	}

	i, err := InitInterrogator(ctx, cfg.Target)
	if err != nil {
		return err
	}

	// service
	j := judge.New(m, i)

	// do it
	fmt.Println("Attacking agent via", cfg.Evaluator.Provider, "...")

	results := j.Judge(ctx, cfg.Evaluator.AttackCategories, cfg.Evaluator.Policy)

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

		fmt.Printf("\n[%d] %s %s\n", i+1, icon, result.Reasoning)
	}

	return nil
}

func InitModel(ctx context.Context, cfg *v1alpha1.EvaluatorConfig) (llms.Model, error) {
	var model llms.Model
	var err error

	switch cfg.Provider {
	case "vertex":
		projectId, ok := cfg.Params["project_id"].(string)
		if !ok {
			return nil, fmt.Errorf("failed to parse vertex project_id as string")
		}
		model, err = vertex.New(
			ctx,
			googleai.WithCloudProject(projectId),
			googleai.WithCloudLocation("us-central1"),
			googleai.WithDefaultModel("gemini-2.0-flash-001"),
		)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", cfg.Provider)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to init %s llm: %w", cfg.Provider, err)
	}

	return model, nil
}

func InitInterrogator(_ context.Context, cfg *v1alpha1.TargetConfig) (interrogator.Interrogator, error) {
	return http.NewInterrogator(
		interrogator.WithTarget(cfg.URL),
	), nil
}
