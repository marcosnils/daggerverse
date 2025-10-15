package main

import "context"

type ChecksTest struct{}

func (m *ChecksTest) CheckFoo() {
	dag.Container().From("alpine:latest").WithExec([]string{"echo", "hi check here"}).Sync(context.Background())
}
