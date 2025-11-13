// A generated module for Scaleout functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

package main

import (
	"context"
)

type Scaleout struct{}

func (m *Scaleout) Work(ctx context.Context) {
	dag.Container().From("alpine").WithExec([]string{"echo", "doing some very tough work!!"}).
		WithExec([]string{"sleep", "10"}).Stdout(ctx)
	return
}
