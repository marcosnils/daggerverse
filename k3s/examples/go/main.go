// A generated module for Examples functions
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
	"time"
)

type Examples struct{}

// starts a k3s server with a local registry and a pre-loaded alpine image
func (m *Examples) K3SServer(ctx context.Context) (*Service, error) {
	regSvc := dag.Container().From("registry:2.8").
		WithExposedPort(5000).AsService()

	_, err := dag.Container().From("quay.io/skopeo/stable").
		WithServiceBinding("registry", regSvc).
		WithEnvVariable("BUST", time.Now().String()).
		WithExec([]string{"copy", "--dest-tls-verify=false", "docker://docker.io/alpine:latest", "docker://registry:5000/alpine:latest"}).Sync(ctx)
	if err != nil {
		return nil, err
	}

	return dag.K3S("test").With(func(k *K3S) *K3S {
		return k.WithContainer(
			k.Container().
				WithEnvVariable("BUST", time.Now().String()).
				WithExec([]string{"sh", "-c", `
cat <<EOF > /etc/rancher/k3s/registries.yaml
mirrors:
  "registry:5000":
    endpoint:
      - "http://registry:5000"
EOF`}, ContainerWithExecOpts{SkipEntrypoint: true}).
				WithServiceBinding("registry", regSvc),
		)
	}).Server(), nil
}

func (m *Examples) K3SKubectl(ctx context.Context, args string) (string, error) {
	return dag.K3S("test").Kubectl(ctx, args)
}
