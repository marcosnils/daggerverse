package main

import (
	"context"
)

type ChecksTest struct{}

type CheckStatus struct {
}

func (m *ChecksTest) CheckMatias(ctx context.Context) *CheckStatus {
	dag.Scaleout().Work(ctx)
	return &CheckStatus{}
}
