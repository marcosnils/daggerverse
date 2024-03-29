// A dapr module to setup a dagger sidecar in your applications
//
// This module allows injecting a dapr sidecar into your application so
// you can use dapr features like service invocation, state management,
// pub/sub, etc.

package main

import (
	"context"
	"strconv"
)

type Dapr struct {
	Image string
}

func New(
	// dapr image to use
	// +optional
	// +default="docker.io/daprio/daprd:1.13.0-rc.7"
	image string,
) *Dapr {
	return &Dapr{Image: image}
}

// Dapr creates a new Dapr container with the specified configuration.
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

	dapr := dag.Container().From(m.Image).
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
