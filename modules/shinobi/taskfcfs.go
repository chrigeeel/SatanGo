package shinobi

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
	"github.com/chrigeeel/satango/modules/utility"
)

func taskfcfs(wg *sync.WaitGroup, userData loader.UserDataStruct, id int, password string, profile loader.ProfileStruct) {
	type postResponseStruct struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
	defer wg.Done()
	beginTime := time.Now()
	client, err := utility.CoolClient("localhost")
	if err != nil {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Invalid Proxy!"))
		return
	}
	fmt.Println(colors.TaskPrefix(id) + colors.Yellow("Submitting checkout..."))
	payload, err := json.Marshal(map[string]string{
		"token":    profile.StripeToken,
		"email":    profile.BillingAddress.Email,
		"password": password,
	})
	if err != nil {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to checkout!"))
	}
	shinobiUrl := "https://dashboard.shinobi-scripts.com/api/payment/purchase"
	req, err := http.NewRequest("POST", shinobiUrl, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to checkout!"))
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.93 Safari/537.36")
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to checkout!"))
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to load release!"))
		return
	}

	var postResponse postResponseStruct
	json.Unmarshal([]byte(body), &postResponse)

	if postResponse.Success {
		stopTime := time.Now()
		diff := stopTime.Sub(beginTime)
		fmt.Println(colors.TaskPrefix(id) + colors.Green("Successfully checked out, Check your email!"))
		go utility.SendWebhook(userData.Webhook, utility.WebhookContentStruct{
			Speed:   diff.String(),
			Module:  "Shinobi",
			Site:    "Shinobi",
			Profile: profile.Name,
		})
		payload, _ := json.Marshal(map[string]string{
			"site":     "Shinobi",
			"module":   "Shinobi",
			"speed":    diff.String(),
			"mode":     "Brr mode",
			"password": password,
			"user":     userData.Username,
		})
		req, _ := http.NewRequest("POST", "http://ec2-13-52-240-112.us-west-1.compute.amazonaws.com:3000/checkouts", bytes.NewBuffer(payload))
		req.Header.Set("content-type", "application/json")
		client.Do(req)
		return
	}
	if postResponse.Message == "sold_out" {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Release OOS!"))
		return
	}
	if postResponse.Message == "Cannot read property 'stock' of null" {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Password incorrect!"))
		return
	}
	fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to checkout with reason: "+string(body)))
}
