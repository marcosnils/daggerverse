// Build Ghostty (ghostty.org) from source
//
// Convenient module to build Ghostty from source without the need to install
// Nix or Zig in your local environment

package main

import (
	"dagger/ghostty/internal/dagger"
	"fmt"
	"runtime"
)

type Ghostty struct {
	// +private
	Src *dagger.Directory

	// +private
	Base *dagger.Container
}

func New(
	// +default="tip"
	ghosttyTag string,
	// +default="0.13.0"
	zigVersion string,

	// +default="latest"
	ubuntuVersion string,

	// +optional
	pBase *dagger.Container,
) *Ghostty {

	arch := runtime.GOARCH

	switch arch {
	case "arm64":
		arch = "aarch64"
	case "amd64":
		arch = "x86_64"
	}

	src := dag.Git("github.com/ghostty-org/ghostty").Tag(ghosttyTag).Tree()

	var base *dagger.Container
	if pBase != nil {
		base = pBase

	} else {

		zigBundle := dag.HTTP(fmt.Sprintf("https://ziglang.org/download/%s/zig-linux-%s-%s.tar.xz", zigVersion, arch, zigVersion))

		base = dag.Container().From("ubuntu:"+ubuntuVersion).
			WithExec([]string{"apt", "update"}).
			WithExec([]string{"apt", "install", "--no-install-recommends", "-y", "ncurses-base", "xz-utils", "libgtk-4-dev", "libadwaita-1-dev"}).
			WithMountedCache("/root/.cache/zig", dag.CacheVolume("zig-cache")).
			WithFile("zig.tar.xz", zigBundle).
			WithExec([]string{"mkdir", "/zig"}).
			WithExec([]string{"tar", "-C", "/zig", "-xf", "zig.tar.xz"}).
			WithEnvVariable("PATH", fmt.Sprintf("/zig/zig-linux-%s-%s:$PATH", arch, zigVersion), dagger.ContainerWithEnvVariableOpts{Expand: true}).
			WithWorkdir("/ghostty").
			WithMountedDirectory(".", src)

	}

	return &Ghostty{
		Base: base,
	}
}

// Returns a container that echoes whatever string argument is provided
func (m *Ghostty) Binary() *dagger.File {
	return m.Base.
		WithExec([]string{"zig", "build", "-Doptimize=ReleaseFast"}).
		File("zig-out/bin/ghostty")

}
