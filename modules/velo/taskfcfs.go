package velo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"github.com/chrigeeel/satango/colors"
	"github.com/chrigeeel/satango/loader"
	"github.com/chrigeeel/satango/modules/getpw"
	"github.com/chrigeeel/satango/modules/utility"
)

func VeloTask(wg *sync.WaitGroup, userData loader.UserDataStruct, id int, site siteStruct, password string, profile loader.ProfileStruct) {
	defer wg.Done()

	type getResponseStruct struct {
		Success  bool   `json:"success"`
		Status   string `json:"status"`
		Checkout string `json:"checkout"`
	}

	type checkoutResponseStruct struct {
		Success bool   `json:"success"`
		Status  string `json:"status"`
	}

	client, err := utility.CoolClient("localhost")
	if err != nil {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Invalid Proxy!"))
		return
	}

	//beginTime := time.Now()

	fmt.Println(colors.TaskPrefix(id) + colors.Yellow("Loading Release..."))

	VUrl := "https://api.velo.gg/api/purchase/create?password=" + password + "&referral=null"
	req, err := http.NewRequest("GET", VUrl, nil)
	if err != nil {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to load release!"))
		return
	}

	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.93 Safari/537.36")
	req.Header.Set("x-velo-host", site.BackendName)
	req.Header.Set("x-velo-authorization", profile.DiscordSession)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to load release!"))
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to load release!"))
		return
	}

	var getResponse getResponseStruct
	json.Unmarshal([]byte(body), &getResponse)

	if getResponse.Success == false {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Wrong password or release OOS!"))
		return
	}

	fmt.Println(getResponse.Checkout)

	go getpw.PWSharingSend(userData, password, site.BackendName)

	fmt.Println(colors.TaskPrefix(id) + colors.Yellow("Submitting Stripe..."))

	confirmUrl := "https://api.stripe.com/v1/payment_pages/" + getResponse.Checkout + "/confirm"
	payload := strings.NewReader(
		`eid=NA` +
			`&payment_method=` + profile.StripeToken +
			`&key=` + site.Stripe_public_key)

	req, err = http.NewRequest("POST", confirmUrl, payload)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err = client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(colors.TaskPrefix(id) + colors.Yellow("Confirming Order..."))

	processUrl := "https://api.velo.gg/api/process?type=checkout&checkoutSession=" + getResponse.Checkout + "&password=" + password + "&appliedCode=null"
	req, err = http.NewRequest("GET", processUrl, nil)
	if err != nil {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to load release!"))
		return
	}

	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.93 Safari/537.36")
	req.Header.Set("x-velo-host", site.BackendName)
	req.Header.Set("x-velo-authorization", profile.DiscordSession)

	resp, err = client.Do(req)
	if err != nil {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to load release!"))
		return
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to load release!"))
		return
	}

	var checkoutResponse checkoutResponseStruct
	json.Unmarshal([]byte(body), &checkoutResponse)

	if checkoutResponse.Success == true {
		fmt.Println(colors.TaskPrefix(id) + colors.Green("Successfully checked out! Please give a kiss to Chrigeeel and check email"))
		go utility.SendWebhook(userData.Webhook, utility.WebhookContentStruct{
			Speed:   "bruh idk bro",
			Module:  "Velo",
			Site:    site.DisplayName,
			Profile: profile.Name,
		})
		payload, _ := json.Marshal(map[string]string{
			"site":     site.DisplayName,
			"module":   "Velo",
			"speed":    "idk bro",
			"mode":     "Brr mode",
			"password": "Unknown",
			"user":     userData.Username,
		})
		req, _ := http.NewRequest("POST", "http://ec2-13-52-240-112.us-west-1.compute.amazonaws.com:3000/checkouts", bytes.NewBuffer(payload))
		req.Header.Set("content-type", "application/json")
		client.Do(req)
		return
	}

	if checkoutResponse.Success == false {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to checkout :/ Maybe still success though idk how to log success correctly, check email!!!"))
		return
	}
}