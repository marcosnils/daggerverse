// Compiles rust code into ARM architectures

package main

import (
	"context"
	"dagger/rustarmv-6/internal/dagger"
	"fmt"
)

const patchRawURL = "https://gist.githubusercontent.com/marcosnils/d8ce7c128c344b4d812f5c290a44ef28/raw/f3e302762da830d86235b8ec477eca779fd41d31/ring-armv6.patch"

type Rustarm struct{}

// builds the supplied source code for the armv6hf architecture
// example usage: dagger -m github.com/marcosnils/daggerverse/rustarm call build-armv-6 --path ./src export --path target
func (m *Rustarm) BuildArmv6(ctx context.Context, src *Directory) *Directory {
	armCompiler := dag.HTTP("https://github.com/raspberrypi/tools/archive/5caa7046982f0539cf5380f94da04b31129ed521.tar.gz")

	return m.base().
		WithMountedFile("/mnt/armcc.tar.gz", armCompiler).
		WithExec([]string{"tar", "xzf", "/mnt/armcc.tar.gz", "--strip-components=1", "-C", "/usr/local"}).
		WithEnvVariable("PATH", "/usr/local/arm-bcm2708/arm-linux-gnueabihf/bin:$PATH", dagger.ContainerWithEnvVariableOpts{Expand: true}).
		WithEnvVariable("PATH", "/usr/local/arm-bcm2708/arm-linux-gnueabihf/libexec/gcc/arm-linux-gnueabihf/4.9.3:$PATH", dagger.ContainerWithEnvVariableOpts{Expand: true}).
		WithMountedCache("/root/.cargo", dag.CacheVolume("cargo-cache")).
		WithMountedCache("/root/.rustup", dag.CacheVolume("rustup-cache")).
		WithExec([]string{"sh", "-c", "curl https://sh.rustup.rs -sSf | sh -s -- -y --verbose"}).
		WithEnvVariable("PATH", "/root/.cargo/bin:$PATH", dagger.ContainerWithEnvVariableOpts{Expand: true}).
		WithMountedDirectory("/src", src).
		WithMountedDirectory("/src/ring", m.PatchedRingArmv6()).
		WithWorkdir("/src").
		WithExec([]string{"rustup", "target", "add", "arm-unknown-linux-gnueabihf"}).
		WithExec([]string{"sh", "-c", `echo '[target.arm-unknown-linux-gnueabihf]\nlinker = "arm-linux-gnueabihf-gcc"' > /root/.cargo/config.toml`}).
		WithExec([]string{"sh", "-c", `echo '[patch.crates-io]\nring = { path = "/src/ring", version = "0.17.8"}' >> Cargo.toml`}).
		WithExec([]string{"cargo", "build", "--target=arm-unknown-linux-gnueabihf", "--release"}).
		Directory("target")
}

// patches the rust ring library to enable armv6 compilation
// example: dagger -m github.com/marcosnils/daggerverse/rustarm call patched-ring-armv-6 export --path ./ring
func (m *Rustarm) PatchedRingArmv6() *Directory {
	// ring v0.17.8
	ring := dag.Git("github.com/briansmith/ring", dagger.GitOpts{KeepGitDir: true}).
		Commit("fa98b490bcbc99a01ff150896ec74c1813242d7f").
		Tree()
	return m.base().
		WithMountedDirectory("/src", ring).
		WithWorkdir("/src").
		WithExec([]string{"sh", "-c", fmt.Sprintf("curl -sS %s | git apply", patchRawURL)}).
		Directory("/src")
}

func (m *Rustarm) base() *Container {
	return dag.Container().From("ubuntu:22.04").
		WithExec([]string{"apt", "update"}).
		WithExec([]string{"apt", "install", "-y", "curl", "git"})
}
