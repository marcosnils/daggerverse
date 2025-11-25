package main

import (
	"context"
)

type ChecksTest struct{}

type CheckStatus struct{}

// +check
func (m *ChecksTest) CheckMatias(ctx context.Context) *CheckStatus {
	// lala
	dag.Scaleout().Work(ctx)
	return &CheckStatus{}
}
