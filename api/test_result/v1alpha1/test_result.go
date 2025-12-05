package v1alpha1

type TestResult struct {
	Passed       bool
	Score        int
	Reasoning    string
	Conversation []Message
	Error        string
}

type Message struct {
	Role    string
	Content string
}
