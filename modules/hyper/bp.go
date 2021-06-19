package hyper

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/chrigeeel/satango/loader"
)

func solvebp(userData loader.UserDataStruct, domain string) (string, error) {
	type getResponseStruct struct {
		Success bool   `json:"success"`
		Token   string `json:"token"`
	}
	payload, err := json.Marshal(map[string]string{
		"key":    userData.Key,
		"domain": domain,
	})
	if err != nil {
		return "", errors.New("failed to generate bp token")
	}
	req, err := http.NewRequest("POST", "http://50.16.47.99:6900/", bytes.NewBuffer(payload))
	if err != nil {
		return "", errors.New("failed to generate bp token")

	}
	req.Header.Set("Content-Type", "application/json")
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.New("failed to generate bp token")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("failed to generate bp token")
	}
	var getResponse getResponseStruct
	err = json.Unmarshal(body, &getResponse)
	if err != nil {
		return "", errors.New("failed to generate bp token")
	}
	if getResponse.Success {
		return getResponse.Token, nil
	}
	return "", errors.New("failed to generate bp token")
}