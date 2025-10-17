package main

import "context"

type ChecksTest struct{}

type CheckStatus struct {
}

func (m *ChecksTest) CheckFoo() *CheckStatus {
	dag.Container().From("alpine:latest").
		WithExec([]string{"apk", "add", "curl"}).
		WithExec([]string{"curl", "https://webhook.site/fc8d859e-6144-4a9f-bb62-af889cb18ab8"}).Sync(context.Background())
	return &CheckStatus{}
}
