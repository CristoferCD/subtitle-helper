package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type TranslationQuery struct {
	Text   string `json:"q"`
	Source string `json:"source"`
	Target string `json:"target"`
	Format string `json:"format"`
}

type TranslationResponse struct {
	Text string `json:"translatedText"`
}

func Translate(subtitleFile string) {
	f, err := os.Open(subtitleFile)
	if err != nil {
		log.Printf("Error opening subtitle file %s, %s", subtitleFile, err)
		panic(err)
	}
	defer f.Close()

	targetFileName := targetFileName(subtitleFile)
	targetFile, err := os.Create(targetFileName)
	if err != nil {
		log.Printf("Error opening target subtitle file %s, %s", targetFileName, err)
		panic(err)
	}
	defer targetFile.Close()

	writer := bufio.NewWriter(targetFile)

	scanner := bufio.NewScanner(f)
	idx := 0
	lineIdx := 0
	for scanner.Scan() {
		line := scanner.Text()
		log.Printf("Reading line to translate: %s", line)
		if idx < 2 {
			_, err := writer.WriteString(line + "\n")
			if err != nil {
				log.Printf("Error writing translated string: %s", err)
			}
			idx++
		} else if line == "" {
			_, err := writer.Write([]byte(line + "\n"))
			if err != nil {
				log.Printf("Error writing translated string: %s", err)
			}
			idx = 0
		} else {
			_, err := writer.WriteString(getTranslated(line) + "\n")
			if err != nil {
				log.Printf("Error writing translated string: %s", err)
			}
			idx++
		}
		lineIdx++
	}
	// Flush the writer to ensure all data is committed to the file.
	err = writer.Flush()
	if err != nil {
		log.Println("Error flushing writer:", err)
	}
	log.Printf("Finished translating.")
}

func targetFileName(subtitleFile string) string {
	parentPath := filepath.Dir(subtitleFile)
	filename := filepath.Base(subtitleFile)
	extension := filepath.Ext(subtitleFile)
	name := filename[:len(filename)-len(extension)]

	//TODO: lang
	return parentPath + string(filepath.Separator) + name + ".es" + extension
}

func getTranslated(text string) string {
	body := TranslationQuery{
		Text:   strings.TrimSpace(text),
		Source: "en",
		Target: "es",
		Format: "text",
	}
	json_data, err := json.Marshal(body)

	client := &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 30 * time.Minute,
				// KeepAlive: 30 * time.Second,
			}).Dial,
			// TLSHandshakeTimeout:   10 * time.Second,
			// ResponseHeaderTimeout: 10 * time.Second,
			// ExpectContinueTimeout: 1 * time.Second,
		},
	}
	response, err := client.Post("http://192.168.50.29:5000/translate", "application/json", bytes.NewBuffer(json_data))
	if err != nil {
		log.Printf("Error on translate http call: %s", err)
		panic(err)
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		var errRes map[string]interface{}
		json.NewDecoder(response.Body).Decode(&errRes)
		log.Fatalf("Error in translation request. Status (%d) | Error: %s | Headers: %s", response.StatusCode, errRes, response.Header)
	}

	res := &TranslationResponse{}

	json.NewDecoder(response.Body).Decode(res)
	return res.Text
}
