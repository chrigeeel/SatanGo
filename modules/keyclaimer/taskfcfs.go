package keyclaimer

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

func taskfcfs(wg *sync.WaitGroup, userData loader.UserDataStruct, id int, key string, profile loader.ProfileStruct, site siteStruct) {
	defer wg.Done()

	type checkoutDataStruct struct {
		License_key string `json:"license_key"`
	}

	type claimResponseStruct struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
	
	beginTime := time.Now()

	checkoutData := checkoutDataStruct{}
	checkoutData.License_key = key

	payload, _ := json.Marshal(checkoutData)
	client := http.DefaultClient

	fmt.Println(colors.TaskPrefix(id) + colors.Yellow("Submitting claim..."))

	req, err := http.NewRequest("POST", site.URL + "api/activate", bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to initiate claim!"))
		return
	}
	req.Header.Set("usera-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.101 Safari/537.36")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("cookie", "dashboard_session="+profile.DiscordSession)
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
	go utility.SendWebhook(userData.Webhook, utility.WebhookContentStruct{
		Speed:   diff.String(),
		Module:  "Key Claimer",
		Site:    site.DisplayName,
		Profile: profile.Name,
	})
	payload, _ = json.Marshal(map[string]string{
		"site":     site.DisplayName,
		"module":   "Key Claimer",
		"speed":    diff.String(),
		"mode":     "Normal",
		"password": "Unknown",
		"user":     userData.Username,
	})
	req, _ = http.NewRequest("POST", "http://ec2-13-52-240-112.us-west-1.compute.amazonaws.com:3000/checkouts", bytes.NewBuffer(payload))
	req.Header.Set("content-type", "application/json")
	client.Do(req)
}