package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
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
	cmd := exec.Command("ffprobe", "-select_streams", "s", "-show_entries", "stream=index,codec_long_name:stream_tags=language,duration", "-of", "json", filePath, "-v", "0")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	out, err := cmd.Output()
	if err != nil {
		log.Printf("Error %s: %s", err, stderr.String())
		panic(err)
	}

	var ffmpegOut FfprobeOutput
	jsonErr := json.Unmarshal(out, &ffmpegOut)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return ffmpegOut.Streams
}

func ExtractSubtitle(path string, subtitleStream FfmpegStream, progressWriter func(float64)) string {
	//ffmpeg -i video.mkv -map 0:s:0 subs.srt
	streamArg := fmt.Sprintf("0:%d", subtitleStream.Index)
	fileName := fileNameWithoutExtension(path)
	parentPath := filepath.Dir(path)
	//TODO: language in name
	srtName := fmt.Sprintf("%s%c%s.srt", parentPath, filepath.Separator, fileName)

	duration := subtitleStream.Tags["DURATION"]
	layout := "15:04:05.000000000"
	durationTime, parseErr := time.Parse(layout, duration)
	if parseErr != nil {
		fmt.Println("Error parsing time:", parseErr)
		panic(parseErr)
	}

	cmd := exec.Command("ffmpeg", "-i", path, "-map", streamArg, "-progress", "unix://"+TempSock(durationInMs(durationTime), progressWriter), srtName)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Printf("Error %s: %s", err, stderr.String())
		panic(err)
	}

	return srtName
}

func durationInMs(duration time.Time) int64 {
	hours := duration.Hour()
	minutes := duration.Minute()
	seconds := duration.Second()
	nanoseconds := duration.Nanosecond()

	return int64(hours)*3600000 +
		int64(minutes)*60000 +
		int64(seconds)*1000 +
		int64(nanoseconds)/1_000_000
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

func TempSock(totalDuration int64, progressWriter func(float64)) string {
	// serve

	sockFileName := path.Join(os.TempDir(), fmt.Sprintf("%d_sock", rand.Int()))
	l, err := net.Listen("unix", sockFileName)
	if err != nil {
		panic(err)
	}

	go func() {
		re := regexp.MustCompile(`out_time_ms=(\d+)`)
		fd, err := l.Accept()
		if err != nil {
			log.Fatal("accept error:", err)
		}
		buf := make([]byte, 16)
		data := ""
		progress := 0.0
		for {
			_, err := fd.Read(buf)
			if err != nil {
				return
			}
			data += string(buf)
			a := re.FindAllStringSubmatch(data, -1)
			cp := 0.0
			if len(a) > 0 && len(a[len(a)-1]) > 0 {
				c, _ := strconv.Atoi(a[len(a)-1][len(a[len(a)-1])-1])
				cp = float64(c) / float64(totalDuration)
				// fmt.Print("Values: ")
				// fmt.Print(a)
				// fmt.Print()
				// fmt.Printf("Read progress: %d / %d\r\n", c, totalDuration)
			}
			time.Sleep(2 * time.Second)
			if strings.Contains(data, "progress=end") {
				cp = 100
			}
			if cp != progress {
				progress = cp
				progressWriter(progress)
				// fmt.Println("progress: ", fmt.Sprintf("%.2f", progress))
			}
		}
	}()

	return sockFileName
}
