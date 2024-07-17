package main

import "dagger/windows/internal/dagger"

type Windows struct {
	// +private
	Ctr *dagger.Container
}

func New() *Windows {
	return &Windows{Ctr: dag.Container().From("dockurr/windows")}
}

func (m *Windows) Service(
	// +optional
	// +default="win11"
	version string,
	// +optional
	// +default="4G"
	ram string,
	// +optional
	// +default="2"
	cpu string,
	// +optional
	// +default="64G"
	disk string,
) *dagger.Service {
	return m.Ctr.
		WithExposedPort(8006).
		WithEnvVariable("VERSION", version).
		WithEnvVariable("RAM_SIZE", ram).
		WithEnvVariable("CPU_CORES", cpu).
		WithEnvVariable("DISK_SIZE", disk).
		WithMountedTemp("/storage").
		WithExec([]string{"/usr/bin/tini", "-s", "/run/entry.sh"}, dagger.ContainerWithExecOpts{InsecureRootCapabilities: true}).
		AsService()
}
