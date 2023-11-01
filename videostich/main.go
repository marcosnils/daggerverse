package main

import (
	"context"
	"fmt"
	"path/filepath"
)

type Videostich struct{}

// example usage: "dagger call container-echo --string-arg yo"
func (m *Videostich) Stitch(ctx context.Context, srcDir *Directory) (*File, error) {
	files, err := srcDir.Entries(ctx)
	if err != nil {
		return nil, err
	}

	mp4s := []string{}
	for _, f := range files {
		if filepath.Ext(f) == ".mp4" {
			mp4s = append(mp4s, f)
		}
	}

	ffc := dag.Container().From("jrottenberg/ffmpeg")

	fmt.Printf("Files to process %v", mp4s)

	concat := ""
	for i, mp4 := range mp4s {

		intermediateName := fmt.Sprintf("intermediate%d.ts", i)
		ffc = ffc.WithFile(mp4, srcDir.File(mp4))
		ffc = ffc.WithExec([]string{"-i", mp4, "-c", "copy", intermediateName})
		if i == 0 {
			concat += intermediateName
		} else {
			concat += "|" + intermediateName
		}
	}

	ffc = ffc.WithExec([]string{"-i", "concat:" + concat, "-c", "copy", "output.mp4"})

	return ffc.File("output.mp4"), nil
}
