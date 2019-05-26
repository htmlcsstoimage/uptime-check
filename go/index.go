package handler

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
	hcti "uptime-check/hcti"
)

type pingdom_http_custom_check struct {
	Status       string `xml:"status"`
	ResponseTime int64  `xml:"response_time"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	password := r.FormValue("password")
	if password != os.Getenv("UPTIME_PASSWORD") {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Unauthorized")
		return
	}

	w.Header().Set("Content-Type", "application/xml")

	now := time.Now()
	url, createTime, err := hcti.GenerateImage(now.Format(time.RFC3339Nano), "css")

	if err != nil {
		response := pingdom_http_custom_check{"DOWN", 0}
		x, _ := xml.MarshalIndent(response, "", "  ")

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(x)
		return
	}

	downloadTime, err := downloadImage(url)

	totalTime := downloadTime + createTime

	if err != nil {
		response := pingdom_http_custom_check{"DOWN", 0}
		x, _ := xml.MarshalIndent(response, "", "  ")

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(x)
	} else {
		response := pingdom_http_custom_check{"UP", totalTime}
		x, _ := xml.MarshalIndent(response, "", "  ")

		w.WriteHeader(http.StatusOK)
		w.Write(x)
	}
}

func downloadImage(url string) (timeElapsed int64, err error) {
	start := time.Now()

	resp, err := http.Get(url)

	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	return int64(time.Since(start) / time.Millisecond), err
}
