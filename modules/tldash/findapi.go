package tldash

import (
	"fmt"
	"net/http"
	"time"

	"github.com/chrigeeel/satango/colors"
)

func findApi() string {
	client := http.Client{Timeout: 2 * time.Second}
	req, err := http.NewRequest("GET", "http://3.140.14.227:5069", nil)
	if err != nil {
		fmt.Println(colors.Prefix() + colors.Green("Using 24/7 API"))
		return "http:/44.195.187.171:5069/v1/solve"
	}
	_, err = client.Do(req)
	if err != nil {
		fmt.Println(colors.Prefix() + colors.Green("Using 24/7 API"))
		return "http://44.195.187.171:5069/v1/solve"
	}
	fmt.Println(colors.Prefix() + colors.Green("Using BIG API"))
	return "http://3.140.14.227:5069/v1/solve"
}