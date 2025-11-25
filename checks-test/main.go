package main

import (
	"context"
)

type ChecksTest struct{}

type CheckStatus struct{}

// +check
func (m *ChecksTest) CheckMatias(ctx context.Context) *CheckStatus {
	// foo
	// bar
	// baz
	// qux
	// cuack
	// woof
	dag.Scaleout().Work(ctx)
	return &CheckStatus{}
}
