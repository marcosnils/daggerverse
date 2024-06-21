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

// entrypoint to setup cgroup nesting since k3s only does it
// when running as PID 1. This doesn't happen in Dagger given that we're using
// our custom shim
const entrypoint = `#!/bin/sh

set -o errexit
set -o nounset

#########################################################################################################################################
# DISCLAIMER																																																														#
# Copied from https://github.com/moby/moby/blob/ed89041433a031cafc0a0f19cfe573c31688d377/hack/dind#L28-L37															#
# Permission granted by Akihiro Suda <akihiro.suda.cz@hco.ntt.co.jp> (https://github.com/k3d-io/k3d/issues/493#issuecomment-827405962)	#
# Moby License Apache 2.0: https://github.com/moby/moby/blob/ed89041433a031cafc0a0f19cfe573c31688d377/LICENSE														#
#########################################################################################################################################
if [ -f /sys/fs/cgroup/cgroup.controllers ]; then
  echo "[$(date -Iseconds)] [CgroupV2 Fix] Evacuating Root Cgroup ..."
	# move the processes from the root group to the /init group,
  # otherwise writing subtree_control fails with EBUSY.
  mkdir -p /sys/fs/cgroup/init
  busybox xargs -rn1 < /sys/fs/cgroup/cgroup.procs > /sys/fs/cgroup/init/cgroup.procs || :
  # enable controllers
  sed -e 's/ / +/g' -e 's/^/+/' <"/sys/fs/cgroup/cgroup.controllers" >"/sys/fs/cgroup/cgroup.subtree_control"
  echo "[$(date -Iseconds)] [CgroupV2 Fix] Done"
fi

exec "$@"
`

func (m *Examples) Server(ctx context.Context) (*Service, error) {
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

func (m *Examples) Kubectl(ctx context.Context, args string) (string, error) {
	return dag.K3S("test").Kubectl(ctx, args)
}
