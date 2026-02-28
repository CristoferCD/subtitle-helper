package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
)

type FfprobeOutput struct {
	Streams []FfmpegStream `json:"streams"`
}

type FfmpegStream struct {
	Index         int               `json:"index"`
	CodecLongName string            `json:"codec_long_name"`
	Tags          map[string]string `json:"tags"`
}

func ListSubtitleStreams(filePath string) []FfmpegStream {
	cmd := exec.Command("ffprobe", "-select_streams", "s", "-show_entries", "stream=index,codec_long_name:stream_tags=language", "-of", "json", filePath, "-v", "0")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	out, err := cmd.Output()
	if err != nil {
		log.Fatalf("Error %s: %s", err, stderr.String())
		panic(err)
	}

	var ffmpegOut FfprobeOutput
	jsonErr := json.Unmarshal(out, &ffmpegOut)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return ffmpegOut.Streams
}

func ExtractSubtitle(path string, subtitleStream int) string {
	//ffmpeg -i video.mkv -map 0:s:0 subs.srt
	streamArg := fmt.Sprintf("0:%d", subtitleStream)
	fileName := fileNameWithoutExtension(path)
	parentPath := filepath.Dir(path)
	//TODO: language in name
	srtName := fmt.Sprintf("%s%c%s.srt", parentPath, filepath.Separator, fileName)

	cmd := exec.Command("ffmpeg", "-i", path, "-map", streamArg, srtName)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Fatalf("Error %s: %s", err, stderr.String())
		panic(err)
	}

	return srtName
}

func fileNameWithoutExtension(path string) string {
	// Get full filename with extension
	filename := filepath.Base(path)
	// Get only the extension
	extension := filepath.Ext(path)
	// Get filename without extension
	name := filename[:len(filename)-len(extension)]
	return name
}
