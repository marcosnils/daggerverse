// Runs a k3s server than can be accessed both locally and in your pipelines

package main

import (
	"context"
	"time"
)

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

type K3S struct {
	// +private
	Name string

	// +private
	ConfigCache *CacheVolume

	Container *Container
}

func New(name string) *K3S {
	ccache := dag.CacheVolume("k3s_config_" + name)
	ctr := dag.Container().
		From("rancher/k3s").
		WithNewFile("/usr/bin/entrypoint.sh", ContainerWithNewFileOpts{
			Contents:    entrypoint,
			Permissions: 0o755,
		}).
		WithEntrypoint([]string{"entrypoint.sh"}).
		WithMountedCache("/etc/rancher/k3s", ccache).
		WithMountedTemp("/etc/lib/cni").
		WithMountedTemp("/var/lib/kubelet").
		WithMountedTemp("/var/lib/rancher/k3s").
		WithMountedTemp("/var/log").
		WithExec([]string{"sh", "-c", "k3s server --bind-address $(ip route | grep src | awk '{print $NF}') --disable traefik --disable metrics-server"}, ContainerWithExecOpts{InsecureRootCapabilities: true}).
		WithExposedPort(6443)
	return &K3S{
		Name:        name,
		ConfigCache: ccache,
		Container:   ctr,
	}
}

// Returns a newly initialized kind cluster
func (m *K3S) Server() *Service {
	return m.Container.AsService()
}

// returns the config file for the k3s cluster
func (m *K3S) Config(ctx context.Context,
	// default=false
	local bool,
) *File {
	return dag.Container().
		From("alpine").
		// we need to bust the cache so we don't fetch the same file each time.
		WithEnvVariable("CACHE", time.Now().String()).
		WithMountedCache("/cache/k3s", m.ConfigCache).
		WithExec([]string{"cp", "/cache/k3s/k3s.yaml", "k3s.yaml"}).
		With(func(c *Container) *Container {
			if local {
				c = c.WithExec([]string{"sed", "-i", `s/https:.*:6443/https:\/\/localhost:6443/g`, "k3s.yaml"})
			}
			return c
		}).
		File("k3s.yaml")
}

// runs kubectl on the target k3s cluster
func (m *K3S) Kubectl(ctx context.Context, args string) (string, error) {
	return dag.Container().
		From("bitnami/kubectl").
		WithoutEntrypoint().
		WithMountedCache("/cache/k3s", m.ConfigCache).
		WithEnvVariable("CACHE", time.Now().String()).
		WithFile("/.kube/config", m.Config(ctx, false), ContainerWithFileOpts{Permissions: 1001}).
		WithUser("1001").
		WithExec([]string{"sh", "-c", "kubectl " + args}).Stdout(ctx)
}