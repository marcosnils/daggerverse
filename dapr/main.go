package main

import (
	"context"
	"strconv"
)

type Component struct {
	Name string
	Type string
}

type Dapr struct {
	components []*Component
}

func (m *Dapr) Dapr(
	ctx context.Context,
	appId string,
	appPort Optional[int],
	componentsPath Optional[*Directory],
) *Service {
	// daprPlacement := dag.Container().From("daprio/placement:1.12.2").
	// WithEntrypoint([]string{"./placement"}).
	// WithExposedPort(50005).
	// AsService()

	args := []string{"./daprd", "-components-path", "/components", "-app-id", appId, "-log-level", "debug"}

	if p, ok := appPort.Get(); ok {
		args = append(args, "-app-port", strconv.Itoa(p))
	}

	dapr := dag.Container().From("daprio/daprd:1.12.2").
		With(func(c *Container) *Container {
			if d, ok := componentsPath.Get(); ok {
				c = c.WithDirectory("/components", d)
			}
			return c
		}).
		// WithServiceBinding("placement", daprPlacement).
		// TODO set app id to right value
		WithEntrypoint(args).
		WithExposedPort(50001).
		WithExposedPort(3500).AsService()
	return dapr
}

func (m *Dapr) WithComponents(c ...*Component) *Dapr {
	m.components = append(m.components, c...)
	return m
}

func (m *Dapr) Run(dir *Directory, componentsPath Optional[*Directory]) *Container {
	return dag.Go().WithSource(dir).Container().
		WithServiceBinding("dapr", m.Dapr(context.Background(), "app", OptEmpty[int](), componentsPath)).
		WithEnvVariable("DAPR_GRPC_PORT", "50001").
		WithEnvVariable("DAPR_GRPC_ENDPOINT", "dapr").
		WithExec([]string{"go", "run", "app.go"})
}

//func (m *Dapr) GoWithDapr(dir *Directory) *Container {
//return dag.Go().Container().
//With(m.WithComponents(
//&Component{"kvstore", "state.in-memory"},
//&Component{"orderpubsub", "pubsub.in-memory"},
//).Start).
//WithDirectory("/app", dir).
//WithWorkdir("/app").
//WithExec([]string{"go", "run", "app.go"})
//}
