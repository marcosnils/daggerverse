package main

import (
	"context"
	"strconv"
)

type Dapr struct{}

func (m *Dapr) Dapr(
	ctx context.Context,
	appId string,
	// +optional
	appPort *int,
	// +optional
	appChannelAddress *string,
	// +optional
	componentsPath *Directory,
) *Container {
	// daprPlacement := dag.Container().From("daprio/placement:1.12.2").
	// WithEntrypoint([]string{"./placement"}).
	// WithExposedPort(50005).
	// AsService()

	args := []string{"./daprd", "-components-path", "/components", "-app-id", appId, "-log-level", "debug"}

	if appPort != nil {
		args = append(args, "-app-port", strconv.Itoa(*appPort))
	}

	if appChannelAddress != nil {
		args = append(args, "-app-channel-address", *appChannelAddress)
	}

	dapr := dag.Container().From("docker.io/daprio/daprd:1.13.0-rc.2").
		With(func(c *Container) *Container {
			if componentsPath != nil {
				c = c.WithDirectory("/components", componentsPath)
			}
			return c
		}).
		// WithServiceBinding("placement", daprPlacement).
		// TODO set app id to right value
		WithEntrypoint(args).
		WithExposedPort(50001).
		WithExposedPort(3500)
	return dapr
}
