package main

import (
	"context"
	"strings"
)

type Docker struct{}

func (m *Docker) Import(ctx context.Context, c *Container, ref string) error {
	dockerSock := dag.Host().UnixSocket("/var/run/docker.sock")

	cliCtr := dag.Container().
		From("docker:cli").
		WithUnixSocket("/var/run/docker.sock", dockerSock)

	out, err := cliCtr.WithMountedFile("/ctr.tar", c.AsTarball()).
		WithExec([]string{"load", "-q", "-i", "/ctr.tar"}).Stdout(ctx)
	if err != nil {
		return err
	}

	imgSHA := strings.Split(out, "sha256:")[1][:31]

	_, err = cliCtr.WithExec([]string{"tag", imgSHA, ref}).Sync(ctx)
	if err != nil {
		return err
	}

	return nil
}
