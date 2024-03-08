// container utils module
//
// This file contains the container utils module, which is a set of functions that
// can be used to modify a container and enhance its functionality in different ways.
package main

import (
	"bytes"
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

func (m *Utils) WithCommands(c *Container, cmds [][]string) (*Container, error) {
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

	return c.WithNewFile("/entrypoint.sh", ContainerWithNewFileOpts{
		Contents:    b.String(),
		Permissions: 0755,
	}).
		WithMountedCache("/tmp/jobs", logsCache).
		WithEntrypoint([]string{"/entrypoint.sh"}).
		WithExec(nil), nil
}
