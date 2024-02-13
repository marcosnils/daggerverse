# Utils

## WithCommands

Allows to run multiple commands in the same container. The first command
that exits will stop the execution of the container and will return its status code.

This solution requires **bash** to be present in the container passed to the
`WithCommands` function.

example:

```go
	c, _ := dag.Utils().WithCommands(dag.Container().From("alpine").
		WithExec([]string{"apk", "add", "bash"}), [][]string{
		{"sh", "-c", `"echo hello; sleep 2"`},
		{"sh", "-c", `"echo bye; sleep 5"`},
	})
}
```

> **Note**
> The executed commands stdout and stderr won't be shown in the Dagger log as
> for this to work it's necessary to redirect them so the Dagger shim
> allows the entrypoint to exit successfully. They will be present under the
> /tmp/jobs folder cache volume in case they're needed.
