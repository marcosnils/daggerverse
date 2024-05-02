// A generated module for Asktube functions
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
	"dagger/asktube/internal/dagger"
)

type GPTools struct {
	BaseCtr *Container
	YTCtr   *Container
	YTURL   string
}

func New() *GPTools {
	base := dag.Container().
		From("python:3-slim-bookworm")

	return &GPTools{
		BaseCtr: base,
		YTCtr: base.WithExec([]string{"apt", "update"}).
			WithExec([]string{"apt", "install", "-y", "ffmpeg"}),
		YTURL: "https://github.com/ytdl-org/ytdl-nightly/releases/latest/download/youtube-dl",
	}
}

// Returns a container that echoes whatever string argument is provided
func (m *GPTools) RAG(openaiApiKey *Secret, source *Directory, question string) (string, error) {
	return m.BaseCtr.
		WithExec([]string{"apt", "update"}).WithExec([]string{"apt", "install", "-y", "build-essential"}).
		WithExec([]string{"pip", "install", "llama-index", "llama-index-vector-stores-chroma"}).
		WithSecretVariable("OPENAI_API_KEY", openaiApiKey).
		WithMountedDirectory("/files", source).
		// WithMountedCache("/tmp/llama_index/rag_cli", dag.CacheVolume("llama_index_rag")).
		WithExec([]string{"llamaindex-cli", "rag", "-v", "--files", "/files", "-q", question}).Stdout(context.Background())
}

// Returns a container that echoes whatever string argument is provided
func (m *GPTools) YtChat(openaiApiKey *Secret, url, question string) (string, error) {
	t := m.Transcript(url)
	d := dag.Directory().WithFile("yt-transcript.txt", t)
	return m.RAG(openaiApiKey, d, question)
}

// Returns the video audio as an mp3 encoded file
func (m *GPTools) Transcript(url string) *File {
	audio := m.Audio(url)
	return m.YTCtr.WithExec([]string{"pip", "install", "openai-whisper"}).
		WithMountedFile("audio.mp3", audio).
		WithMountedCache("/root/.cache/whisper", dag.CacheVolume("whisper")).
		WithExec([]string{"whisper", "audio.mp3", "--model", "base", "--fp16", "False", "-f", "txt"}).
		File("audio.txt")
}

// Returns the video audio as an mp3 encoded file
func (m *GPTools) Audio(url string) *File {
	ytdlCLI := dag.HTTP(m.YTURL)

	return m.YTCtr.WithFile("/usr/local/bin/youtube-dl", ytdlCLI, dagger.ContainerWithFileOpts{Permissions: 0777}).
		WithExec([]string{"youtube-dl", "-x", url, "-o", "audio.%(ext)s", "--audio-format", "mp3"}).
		File("audio.mp3")
}
