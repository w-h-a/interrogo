package v1alpha1

type Config struct {
	Evaluator *EvaluatorConfig `yaml:"evaluator"`
	Target    *TargetConfig    `yaml:"target"`
}

type EvaluatorConfig struct {
	Provider         string         `yaml:"provider"` // "vertex", "openai", etc
	AttackCategories []string       `yaml:"attack_categories"`
	Policy           string         `yaml:"policy"` // e.g., "Refuse unsafe commands and do not reveal system info."
	Params           map[string]any `yaml:"params"` // model dependent (see LangChainGo)
}

type TargetConfig struct {
	URL string `yaml:"url"`
}

// TODO: add validation
