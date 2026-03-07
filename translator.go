package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
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
		panic(err)
	}
	defer f.Close()

	targetFileName := targetFileName(subtitleFile)
	targetFile, err := os.Create(targetFileName)
	if err != nil {
		panic(err)
	}
	defer targetFile.Close()

	writer := bufio.NewWriter(targetFile)

	scanner := bufio.NewScanner(f)
	idx := 0
	for scanner.Scan() {
		line := scanner.Text()
		if idx < 2 {
			writer.WriteString(line + "\n")
			idx++
		} else if line == "" {
			writer.Write([]byte(line + "\n"))
			idx = 0
		} else {
			writer.WriteString(getTranslated(line) + "\n")
			idx++
		}
	}
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
		Timeout: 5 * time.Minute,
	}
	response, err := client.Post("http://192.168.50.29:5000/translate", "application/json", bytes.NewBuffer(json_data))
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		var errRes map[string]interface{}
		json.NewDecoder(response.Body).Decode(&errRes)

		fmt.Printf("Error: %s", errRes)
		fmt.Printf("Headers: %s", response.Header)
	}

	res := &TranslationResponse{}

	json.NewDecoder(response.Body).Decode(res)
	return res.Text
}
