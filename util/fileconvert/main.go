package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	wd, _ := os.Getwd()

	path := filepath.Join(filepath.Dir(filepath.Dir(wd)), "audio")
	filepath.WalkDir(path, process)
}

func process(path string, d fs.DirEntry, err error) error {
	if strings.HasSuffix(d.Name(), ".mp3") {
		formatAudio(strings.TrimSuffix(path, ".mp3"))
	}
	return nil
}

// this would be simpler to do like so:
// https://stackoverflow.com/questions/10781516/how-to-pipe-several-commands-in-go
//
// but I wanted to figure it out, and once it works then there's no point changing it.
func formatAudio(path string) error {
	_, err := exec.LookPath("ffmpeg")
	if err != nil {
		panic("could not find ffmpeg")
	}
	_, err = exec.LookPath("dca")
	if err != nil {
		panic("could not find dca")
	}

	ffmpegFileArg := fmt.Sprintf("%s.mp3", path)
	ffmpegCmd := exec.Command("ffmpeg", "-i", ffmpegFileArg, "-f", "s16le", "-ar", "48000", "-ac", "2", "pipe:1")

	ffmpegErrOut := bytes.Buffer{}
	ffmpegCmd.Stderr = &ffmpegErrOut

	ffmpegCmdOut, err := ffmpegCmd.StdoutPipe()
	if err != nil {
		panic(err)
	}

	dcaCmd := exec.Command("dca")
	dcaErrOut := bytes.Buffer{}

	dcaCmd.Stderr = &dcaErrOut
	dcaCmd.Stdin = ffmpegCmdOut

	outfile, _ := os.Create(path + ".dca")
	dcaCmd.Stdout = outfile

	dcaCmd.Start()
	ffmpegCmd.Start()

	ffmpegErr := ffmpegCmd.Wait()
	dcaErr := dcaCmd.Wait()

	if ffmpegErr != nil {
		panic(ffmpegErrOut.String())
	}

	if dcaErr != nil {
		panic(dcaErrOut.String())
	}

	err = outfile.Close()
	return err
}
