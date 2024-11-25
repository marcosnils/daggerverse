// container utils module
//
// This file contains the container utils module, which is a set of functions that
// can be used to modify a container and enhance its functionality in different ways.
package main

import (
	"bytes"
	"context"
	"main/internal/dagger"
	"strings"
	"text/template"
)

type Utils struct{}

// it'd be awesome if this function could be define as follows:
//func (m *Utils) WithCommands(cmds [][]string) WithContainerFunc {
//	return func(c *Container) *Container {
//		return c
//	}
//}

const cmdsTemplate = `#!/bin/bash

# Start processes
{{range $index, $element := .}}
{{$element}} > /tmp/jobs/{{$index}}.out 2> /tmp/jobs/{{$index}}.err &
{{end}}


wait -n
# Exit with status of process that exited first
exit $?
`

var logsCache = dag.CacheVolume("logs")

// WithCommands takes a container and a slice of command arguments and returns a new container
// with the specified commands. It generates a shell script based on the provided commands
// and sets it as the entrypoint for the container. The commands are executed when the container
// is started.
//
// If the bash executable is not present, an error will be returned.
//
// Parameters:
//   - c: The original container.
//   - cmds: A slice of command arguments. Each command is represented as a slice of strings.
//
// Returns:
//   - A new container with the specified commands.
//   - An error if the bash executable is not present.
func (m *Utils) WithCommands(c *dagger.Container, cmds [][]string) (*dagger.Container, error) {
	// TODO check if bash is present and return more meaningful error

	// Create a new template and parse the letter into it.
	t := template.Must(template.New("cmds").Parse(cmdsTemplate))

	b := &bytes.Buffer{}

	strCmds := []string{}

	for _, cmd := range cmds {
		cmd := cmd
		strCmds = append(strCmds, strings.Join(cmd, " "))
	}

	t.Execute(b, strCmds)

	return c.WithNewFile("/entrypoint.sh", b.String(), dagger.ContainerWithNewFileOpts{
		Permissions: 0755,
	}).
		WithMountedCache("/tmp/jobs", logsCache).
		WithEntrypoint([]string{"/entrypoint.sh"}).
		WithExec(nil), nil
}

func (m *Utils) WithEnvVariables(ctx context.Context, c *dagger.Container, envs *dagger.Secret) *dagger.Container {

	plainEnvs, _ := envs.Plaintext(ctx)

	for _, e := range strings.Split(plainEnvs, "\n") {
		e := strings.Split(e, "=")
		if len(e) != 2 {
			continue
		}
		c = c.WithEnvVariable(e[0], e[1])
	}

	return c
}
