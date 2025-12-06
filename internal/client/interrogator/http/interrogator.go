package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/w-h-a/interrogo/api/chat_response/v1alpha1"
	"github.com/w-h-a/interrogo/internal/client/interrogator"
)

type httpInterrogator struct {
	options interrogator.Options
	client  *http.Client
}

func (i *httpInterrogator) Interrogate(ctx context.Context, prompt string) (response string, toolCalls []string, err error) {
	body, _ := json.Marshal(map[string]string{"message": prompt})

	req, err := http.NewRequestWithContext(ctx, "POST", i.options.Target, bytes.NewBuffer(body))
	if err != nil {
		return "", nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	rsp, err := i.client.Do(req)
	if err != nil {
		return "", nil, err
	}
	defer rsp.Body.Close()

	bs, _ := io.ReadAll(rsp.Body)
	if rsp.StatusCode != http.StatusOK {
		return "", nil, fmt.Errorf("target error %d: %s", rsp.StatusCode, string(bs))
	}

	var chatRsp v1alpha1.ChatResponse
	if err := json.Unmarshal(bs, &chatRsp); err != nil {
		return "", nil, err
	}

	return chatRsp.Response, chatRsp.ToolCalls, nil
}

func NewInterrogator(opts ...interrogator.Option) interrogator.Interrogator {
	options := interrogator.NewOptions(opts...)

	// TODO: validate options for httpInterrogator

	i := &httpInterrogator{
		options: options,
		client:  &http.Client{},
	}

	return i
}
