package tldash

import (
	"fmt"
	"net/http"
	"time"

	"github.com/chrigeeel/satango/colors"
)

func findApi() string {
	client := http.Client{Timeout: 2 * time.Second}
	req, err := http.NewRequest("GET", "http://35.80.125.25:5069", nil)
	if err != nil {
		fmt.Println(colors.Prefix() + colors.Green("Using 24/7 API"))
		return "http://54.159.151.181:5069/v1/solve"
	}
	_, err = client.Do(req)
	if err != nil {
		fmt.Println(colors.Prefix() + colors.Green("Using 24/7 API"))
		return "http://54.159.151.181:5069/v1/solve"
	}
	return "http://35.80.125.25:5069/v1/solve"
}