package config

import "github.com/tmc/langchaingo/llms"

func ConifgToCallOptions(cfg DefaultCallOptions) []llms.CallOption {
	opts := []llms.CallOption{}

	if cfg.Model != nil {
		opts = append(opts, llms.WithModel(*cfg.Model))
	}
	if cfg.CandidateCount != nil {
		opts = append(opts, llms.WithCandidateCount(*cfg.CandidateCount))
	}
	if cfg.MaxTokens != nil {
		opts = append(opts, llms.WithMaxTokens(*cfg.MaxTokens))
	}
	if cfg.Temperature != nil {
		opts = append(opts, llms.WithTemperature(*cfg.Temperature))
	}
	if cfg.StopWords != nil {
		opts = append(opts, llms.WithStopWords(*cfg.StopWords))
	}
	if cfg.TopK != nil {
		opts = append(opts, llms.WithTopK(*cfg.TopK))
	}
	if cfg.TopP != nil {
		opts = append(opts, llms.WithTopP(*cfg.TopP))
	}
	if cfg.Seed != nil {
		opts = append(opts, llms.WithSeed(*cfg.Seed))
	}
	if cfg.MinLength != nil {
		opts = append(opts, llms.WithMinLength(*cfg.MinLength))
	}
	if cfg.MaxLength != nil {
		opts = append(opts, llms.WithMaxLength(*cfg.MaxLength))
	}
	if cfg.N != nil {
		opts = append(opts, llms.WithN(*cfg.N))
	}
	if cfg.RepetitionPenalty != nil {
		opts = append(opts, llms.WithRepetitionPenalty(*cfg.RepetitionPenalty))
	}
	if cfg.FrequencyPenalty != nil {
		opts = append(opts, llms.WithFrequencyPenalty(*cfg.FrequencyPenalty))
	}
	if cfg.PresencePenalty != nil {
		opts = append(opts, llms.WithPresencePenalty(*cfg.PresencePenalty))
	}
	if cfg.JSONMode != nil && *cfg.JSONMode {
		opts = append(opts, llms.WithJSONMode())
	}
	if cfg.ResponseMIMEType != nil {
		opts = append(opts, llms.WithResponseMIMEType(*cfg.ResponseMIMEType))
	}
	if len(cfg.Metadata) > 0 {
		opts = append(opts, llms.WithMetadata(cfg.Metadata))
	}

	return opts
}
