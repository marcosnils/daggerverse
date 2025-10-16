package main

import "context"

type ChecksTest struct{}

func (m *ChecksTest) CheckFoo() {
	dag.Container().From("alpine:latest").
		WithExec([]string{"apk", "add", "curl"}).
		WithExec([]string{"curl", "https://webhook.site/5b523f5e-1d50-40e3-9f69-6df44197011a"}).Sync(context.Background())

	//foo
}
