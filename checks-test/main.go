package main

import (
	"context"
	"time"
)

type ChecksTest struct{}

type CheckStatus struct {
}

func (m *ChecksTest) CheckMatias() *CheckStatus {
	dag.Container().From("alpine:latest").
		WithExec([]string{"apk", "add", "curl"}).
		WithEnvVariable("CACHE_BUST", time.Now().String()).
		WithExec([]string{"curl", "https://webhook.site/12f90a42-6a19-4864-86bf-f81f841ccc96"}).Sync(context.Background())
	dag.Calltest().CheckTest(context.Background())
	return &CheckStatus{}
}
