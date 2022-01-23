package velo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"github.com/chrigeeel/satango/colors"
	"github.com/chrigeeel/satango/loader"
	"github.com/chrigeeel/satango/modules/getpw"
	"github.com/chrigeeel/satango/utility"
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

	if !getResponse.Success {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Wrong password or release OOS!"))
		return
	}

	pwData := getpw.PWStruct{
		Username: userData.Username,
		Password: password,
		Site: site.BackendName,
		SiteType: "velo",
	}
	go getpw.PWSharingSend2(pwData)

	fmt.Println(colors.TaskPrefix(id) + colors.Yellow("Submitting Stripe..."))

	confirmUrl := "https://api.stripe.com/v1/payment_pages/" + getResponse.Checkout + "/confirm"
	payload := strings.NewReader(
		`eid=NA` +
			`&payment_method=` + profile.StripeToken +
			`&key=` + site.Stripe_public_key)

	req, err = http.NewRequest("POST", confirmUrl, payload)
	if err != nil {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to submit Stripe!"))
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err = client.Do(req)
	if err != nil {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to submit Stripe!"))
		return
	}
	defer resp.Body.Close()

	fmt.Println(colors.TaskPrefix(id) + colors.Yellow("Confirming Order..."))

	processUrl := "https://api.velo.gg/api/process?type=checkout&checkoutSession=" + getResponse.Checkout + "&password=" + password + "&appliedCode=null"
	req, err = http.NewRequest("GET", processUrl, nil)
	if err != nil {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to confirm order!"))
		return
	}

	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.93 Safari/537.36")
	req.Header.Set("x-velo-host", site.BackendName)
	req.Header.Set("x-velo-authorization", profile.DiscordSession)

	resp, err = client.Do(req)
	if err != nil {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to confirm order!"))
		return
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to confirm order!"))
		return
	}

	var checkoutResponse checkoutResponseStruct
	json.Unmarshal([]byte(body), &checkoutResponse)

	if checkoutResponse.Success {
		fmt.Println(colors.TaskPrefix(id) + colors.Green("Maybe checked out! Please check your email to confirm and post #success"))
		return
	}

	if !checkoutResponse.Success {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to checkout :/ Maybe still success though idk how to log success correctly, check email!!!"))
		return
	}
}