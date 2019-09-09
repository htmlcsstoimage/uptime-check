package hcti

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type Image struct {
	Url string
}

type ImageError struct {
	Error   string
	Message string
}

func GenerateImage(html string, css string) (url string, elaspedMS int64, err error) {
	userID := os.Getenv("API_ID")
	apiKey := os.Getenv("API_KEY")

	data := map[string]string{
		"html": html,
		"css":  css,
	}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(data)

	if err != nil {
		return "", 0, err
	}

	req, err := http.NewRequest("POST", "https://hcti.io/v1/image", bytes.NewReader(reqBody))

	if err != nil {
		return "", 0, err
	}

	req.SetBasicAuth(userID, apiKey)

	start := time.Now()

	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)

	if err != nil {
		return "", 0, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", 0, err
	}

	if resp.StatusCode == 200 {
		elapsedMS := int64(time.Since(start) / time.Millisecond)

		var image Image
err := json.NewDecoder(resp.Body).Decode(&image)
		url = image.Url

		return url, elapsedMS, err
	} else {
		var imageError ImageError
    err := json.NewDecoder(resp.Body).Decode(&imageError)

		return "", 0, errors.New(imageError.Message)
	}
}
