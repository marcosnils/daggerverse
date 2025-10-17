package main

import (
	"context"
	"time"
)

type ChecksTest struct{}

type CheckStatus struct {
}

func (m *ChecksTest) CheckBar() *CheckStatus {
	dag.Container().From("alpine:latest").
		WithExec([]string{"apk", "add", "curl"}).
		WithEnvVariable("CACHE_BUST", time.Now().String()).
		WithExec([]string{"curl", "https://webhook.site/2c523046-94cb-45e6-8034-c3aea6e97990"}).Sync(context.Background())
	dag.Calltest().CheckTest(context.Background())
	//trigger

	return &CheckStatus{}
}
