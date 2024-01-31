package main

type Cami struct{}

// example usage: "dagger call hello"
func (m *Cami) Hello() *Container {
	return dag.Container().From("alpine:latest").WithExec([]string{"hello", "cami"})
}
