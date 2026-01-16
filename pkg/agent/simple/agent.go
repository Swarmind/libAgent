package simple

import (
	"context"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

type Agent struct {
	LLM *openai.LLM
}

func (a *Agent) Run(
	ctx context.Context,
	state []llms.MessageContent,
	opts ...llms.CallOption,
) (llms.MessageContent, error) {
	response, err := a.LLM.GenerateContent(
		ctx, state, opts...,
	)
	if err != nil {
		return llms.MessageContent{}, err
	}

	content := response.Choices[0].Content

	return llms.TextParts(llms.ChatMessageTypeAI, content), nil
}

func (a *Agent) SimpleRun(
	ctx context.Context,
	input string,
	opts ...llms.CallOption,
) (string, error) {
	response, err := a.LLM.GenerateContent(ctx,
		[]llms.MessageContent{
			llms.TextParts(llms.ChatMessageTypeHuman,
				input,
			)},
		opts...,
	)
	if err != nil {
		return "", err
	}

	content := response.Choices[0].Content

	return content, nil
}
