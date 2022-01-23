package wrath

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/chrigeeel/satango/colors"
	"github.com/chrigeeel/satango/loader"
	"github.com/chrigeeel/satango/utility"
)

func taskfcfs(wg *sync.WaitGroup, userData loader.UserDataStruct, id int, key string, profile loader.ProfileStruct) {
	defer wg.Done()

	beginTime := time.Now()

	type checkoutDataStruct struct {
		DiscordId string `json:"discordId"`
		Key       string `json:"key"`
	}

	type claimResponseStruct struct {
		Success bool `json:"success"`
	}

	if len(key) != 29 {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Key must be 29 characters long!"))
		return
	}

	checkoutData := checkoutDataStruct{}
	checkoutData.DiscordId = profile.DiscordId
	checkoutData.Key = key

	payload, _ := json.Marshal(checkoutData)
	client := http.DefaultClient

	fmt.Println(colors.TaskPrefix(id) + colors.Yellow("Submitting claim..."))

	req, err := http.NewRequest("POST", "https://server.wrathbots.co/keybind", bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to initiate claim!"))
		return
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.101 Safari/537.36")
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to claim the key!"))
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	rdata := new(claimResponseStruct)
	json.Unmarshal([]byte(body), &rdata)
	if !rdata.Success {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to claim key! Either already claimed or wrong key."))
		return
	}
	fmt.Println(colors.TaskPrefix(id) + colors.Green("Successfully claimed key on profile ") + colors.White("\"") + colors.Green(profile.Name) + colors.White("\""))
	stopTime := time.Now()
	diff := stopTime.Sub(beginTime)
	go utility.NewSuccess(userData.Webhook, utility.SuccessStruct{
		Site: "Wrath",
		Module: "Wrath",
		Mode: "Normal",
		Time: diff.String(),
		Profile: profile.Name,
	})
}