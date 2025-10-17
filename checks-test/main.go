package main

import "context"

type ChecksTest struct{}

type CheckStatus struct {
}

func (m *ChecksTest) CheckFoo() *CheckStatus {
	dag.Container().From("alpine:latest").
		WithExec([]string{"apk", "add", "curl"}).
		WithExec([]string{"curl", "https://webhook.site/2c523046-94cb-45e6-8034-c3aea6e97990"}).Sync(context.Background())
	dag.Calltest().CheckTest(context.Background())

	return &CheckStatus{}
}
