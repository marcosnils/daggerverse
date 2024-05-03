// A toolbox of GPT functions for fun
//
// This module serves as toolbox with different GPT related functions that can
// be used to interact with GPT models. The functions are designed to be used
// either in a stand-alone manner or in combination with other functions for
// more elaborated pipelines.

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

// Runs a RAG model on the provided source directory and question
func (m *GPTools) RAG(openaiApiKey *Secret, source *Directory, question string) (string, error) {
	return m.BaseCtr.
		WithExec([]string{"apt", "update"}).WithExec([]string{"apt", "install", "-y", "build-essential"}).
		WithExec([]string{"pip", "install", "llama-index", "llama-index-vector-stores-chroma"}).
		WithSecretVariable("OPENAI_API_KEY", openaiApiKey).
		WithMountedDirectory("/files", source).
		// WithMountedCache("/tmp/llama_index/rag_cli", dag.CacheVolume("llama_index_rag")).
		WithExec([]string{"llamaindex-cli", "rag", "-v", "--files", "/files", "-q", question}).Stdout(context.Background())
}

// Asks a question to a youtube video
func (m *GPTools) YtChat(openaiApiKey *Secret, url, question string) (string, error) {
	a := m.Audio(url)
	t := m.Transcript(a)
	d := dag.Directory().WithFile("yt-transcript.txt", t)
	return m.RAG(openaiApiKey, d, question)
}

// Returns the video transcript as a txt file
func (m *GPTools) Transcript(src *File) *File {
	return m.YTCtr.WithExec([]string{"pip", "install", "openai-whisper"}).
		WithMountedFile("audio.mp3", src).
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
