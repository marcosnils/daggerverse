package main

import (
	"dagger/go/internal/dagger"
)

type Go struct{}

func (m *Go) Utils_WithEnvVariables(envs *dagger.Secret) *dagger.Container {
	return dag.Utils().WithEnvVariables(dag.Container().From("alpine"), envs)
}
