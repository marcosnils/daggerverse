package main

import "context"

type Examples struct{}

func (m *Examples) GPToolsYtChat(ctx context.Context, openaiApiKey *Secret, url, question string) (string, error) {
	return dag.Gptools().YtChat(ctx, openaiApiKey, url, question)
}
