package tldash

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/chrigeeel/satango/colors"
	"github.com/chrigeeel/satango/loader"
	"github.com/chrigeeel/satango/modules/getpw"
	"github.com/chrigeeel/satango/utility"
)

func taskfcfs(wg *sync.WaitGroup, userData loader.UserDataStruct, id int, password string, solveIp string, task taskStruct) {
	type getResponseStruct struct {
		Stripe_public_key string `json:"stripe_public_key"`
		Price_with_symbol string `json:"price_with_symbol"`
		Captcha           string `json:"captcha"`
		AddressRequired string `json:"address_required"`
		NameRequired string `json:"name_required"`
		Country string `json:"country"`
	}
	type tlError struct {
		Message string `json:"message"`
	}
	type postResponseStruct struct {
		Success bool    `json:"success"`
		Message string  `json:"message"`
		Error   tlError `json:"error"`
	}
	type payloadStruct struct {
		Address struct {
			City string `json:"city,omitempty"`
			Country string `json:"country,omitempty"`
			Line1 string `json:"line1,omitempty"`
			PostalCode string `json:"postal_code,omitempty"`
		} `json:"address,omitempty"`
		Captcha string `json:"captcha,omitempty"`
		Email string `json:"email"`
		Name string `json:"name,omitempty"`
		Token string `json:"token"`
	}
	defer wg.Done()

	proxy := task.Proxy
	profile := task.Profile
	site := task.Site
	stripeToken := task.Profile.StripeToken
	discordSession := task.Profile.DiscordSession
	client, err := utility.CoolClient(proxy)
	if err != nil {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Invalid Proxy!"))
		return
	}
	beginTime := time.Now()

	var captchaSolution string
	var getResponse getResponseStruct
	var cfCookie string

	fmt.Println(colors.TaskPrefix(id) + colors.Yellow("Loading Release..."))

	TLUrl := "https://button-backend.tldash.ai/api/purchase/" + site + "/" + password
	req, err := http.NewRequest("GET", TLUrl, nil)
	if err != nil {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to load release!"))
		return
	}

	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.93 Safari/537.36")
	if discordSession != "" {
		req.Header.Set("authorization", "Bearer "+discordSession)
	}

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

	json.Unmarshal([]byte(body), &getResponse)

	if string(getResponse.Stripe_public_key) == "" {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Wrong password or release OOS!"))
		return
	}

	for _, cookie := range resp.Cookies() {
		if cookie.Name == "__cf_bm" {
			cfCookie = cookie.Value
		}
	}
	if cfCookie != "" {
		pwData := getpw.PWStruct{
			Username: userData.Username,
			Password: password,
			Site: site,
			SiteType: "tldash",
		}
		go getpw.PWSharingSend2(pwData)
		go newstripe(task.Site, getResponse.Stripe_public_key)
	}

	fmt.Println(colors.TaskPrefix(id) + colors.Green("Successfully loaded Release!"))

	if stripeToken == "" {
		type tokenStruct struct {
			ID string `json:"id"`
		}

		strpClient := &http.Client{}
		url := "https://api.stripe.com/v1/tokens"
		payload := strings.NewReader(
			`card[number]=` + profile.PaymentDetails.CardNumber +
				`&card[cvc]=` + profile.PaymentDetails.CardCvv +
				`&card[exp_month]=` + profile.PaymentDetails.CardExpMonth +
				`&card[exp_year]=` + profile.PaymentDetails.CardExpYear[len(profile.PaymentDetails.CardExpYear)-2:] +
				`&key=` + getResponse.Stripe_public_key)

		req, err := http.NewRequest("POST", url, payload)
		if err != nil {
			fmt.Println(colors.TaskPrefix(id) + colors.Red("Stripe rejected your profile ") + colors.White(profile.Name) + colors.Red("!"))
			return
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := strpClient.Do(req)
		if err != nil {
			fmt.Println(colors.TaskPrefix(id) + colors.Red("Stripe rejected your profile ") + colors.White(profile.Name) + colors.Red("!"))
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(colors.TaskPrefix(id) + colors.Red("Stripe rejected your profile ") + colors.White(profile.Name) + colors.Red("!"))
			return
		}
		var token tokenStruct
		json.Unmarshal([]byte(body), &token)
		if token.ID == "" {
			fmt.Println(colors.TaskPrefix(id) + colors.Red("Stripe rejected your profile ") + colors.White(profile.Name) + colors.Red("!"))
			return
		}
		fmt.Println(colors.TaskPrefix(id) + colors.Green("Successfully fetched Stripe token for profile ") + colors.White(profile.Name))
		stripeToken = token.ID
	}

	if getResponse.Captcha != "" {
		type captchaResponseStruct struct {
			Solution        string `json:"solution"`
			Processing_time string `json:"processing_time"`
		}

		fmt.Println(colors.TaskPrefix(id) + colors.White("Capcha enabled!"))

		captchaClient := &http.Client{}
		webhookUrl := "https://discord.com/api/webhooks/820084465497669663/0VZgCoLaBWAuIJ_osAzhaGEGOjsgQp7v_N6gL_GTxIQoUX6rh_AQZJGn74O4f_1Q9AmM"

		payload, _ := json.Marshal(map[string]string{
			"b64":        getResponse.Captcha,
			"key":        userData.Key,
			"username":   userData.Username,
			"webhookUrl": webhookUrl,
		})

		req, err := http.NewRequest("POST", solveIp, bytes.NewBuffer(payload))
		if err != nil {
			fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed solving Captcha!"))
			return
		}

		req.Header.Set("content-type", "application/json")
		req.Header.Set("x-api-key", userData.Key+"-TL")

		resp, err := captchaClient.Do(req)
		if err != nil {
			fmt.Println(err)
			fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed solving Captcha!"))
			return
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to load release!"))
			return
		}

		var captchaResponse captchaResponseStruct
		json.Unmarshal([]byte(body), &captchaResponse)
		captchaSolution = captchaResponse.Solution
		if captchaSolution == "AAAAAA" {
			fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed solving Captcha!"))
			return
		}
		fmt.Println(colors.TaskPrefix(id) + colors.Green("Successfully solved Captcha: ") + colors.White(captchaSolution))
	}


	payload := payloadStruct{
		Email: profile.BillingAddress.Email,
		Token: stripeToken,
	}

	if captchaSolution != "" {
		payload.Captcha = captchaSolution
	}

	if getResponse.NameRequired == "true" {
		payload.Name = profile.BillingAddress.Name
	}

	if getResponse.AddressRequired == "true" {
		payload.Address.City = profile.BillingAddress.Name
		payload.Address.Country = getResponse.Country
		payload.Address.Line1 = profile.BillingAddress.Line1
		payload.Address.PostalCode = profile.BillingAddress.PostCode 
	}
	
	p, err := json.Marshal(payload)
	if err != nil {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to checkout!"))
	}

	fmt.Println(colors.TaskPrefix(id) + colors.Yellow("Submitting Checkout..."))

	TLUrl = "https://button-backend.tldash.ai/api/register/" + site + "/" + password
	req, err = http.NewRequest("POST", TLUrl, bytes.NewBuffer(p))
	if err != nil {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to checkout!"))
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.93 Safari/537.36")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", "__cf_bm="+cfCookie)
	if discordSession != "" {
		req.Header.Set("authorization", "Bearer "+discordSession)
	}

	resp, err = client.Do(req)
	if err != nil {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to checkout!"))
		return
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to load release!"))
		return
	}

	var postResponse postResponseStruct
	json.Unmarshal([]byte(body), &postResponse)

	if postResponse.Success {
		stopTime := time.Now()
		diff := stopTime.Sub(beginTime)
		fmt.Println(colors.TaskPrefix(id) + colors.Green("Successfully checked out on profile ") + colors.White("\"") + colors.Green(profile.Name) + colors.White("\""))
		go utility.SendWebhook(userData.Webhook, utility.WebhookContentStruct{
			Speed:   diff.String(),
			Module:  "TL Dashboards",
			Site:    site,
			Profile: profile.Name,
		})
		payload, _ := json.Marshal(map[string]string{
			"site":     task.Site,
			"module":   "TL Dash",
			"speed":    diff.String(),
			"mode":     "Brr mode",
			"password": "Unknown",
			"user":     userData.Username,
		})
		req, _ := http.NewRequest("POST", "http://ec2-13-52-240-112.us-west-1.compute.amazonaws.com:3000/checkouts", bytes.NewBuffer(payload))
		req.Header.Set("content-type", "application/json")
		client.Do(req)
		return
	}
	if postResponse.Error.Message != "" {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to checkout with reason: ") + colors.White("\"") + colors.Red(postResponse.Error.Message) + colors.White("\"") + colors.Red("!"))
		return
	}
	if postResponse.Message !=  "" {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to checkout with reason: ") + colors.White("\"") + colors.Red(postResponse.Message) + colors.White("\"") + colors.Red("!"))
		return
	}
	fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to checkout for some reason, DM the below to owners:"))
	fmt.Println(colors.TaskPrefix(id) + colors.White(string(body)))
}

func taskfcfsbypass(wg *sync.WaitGroup, userData loader.UserDataStruct, id int, password string, solveIp string, task taskStruct, bypassLevel int) {
	type tlError struct {
		Message string `json:"message"`
	}
	type postResponseStruct struct {
		Success bool    `json:"success"`
		Message string  `json:"message"`
		Error   tlError `json:"error"`
	}
	type payloadStruct struct {
		Address struct {
			City string `json:"city,omitempty"`
			Country string `json:"country,omitempty"`
			Line1 string `json:"line1,omitempty"`
			PostalCode string `json:"postal_code,omitempty"`
		} `json:"address,omitempty"`
		Captcha string `json:"captcha,omitempty"`
		Email string `json:"email"`
		Name string `json:"name,omitempty"`
		Token string `json:"token"`
	}
	defer wg.Done()

	proxy := task.Proxy
	profile := task.Profile
	site := task.Site
	stripeToken := task.Profile.StripeToken
	discordSession := task.Profile.DiscordSession
	client, err := utility.CoolClient(proxy)
	if err != nil {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Invalid Proxy!"))
		return
	}

	beginTime := time.Now()

	payload := payloadStruct{
		Email: profile.BillingAddress.Email,
		Token: stripeToken,
	}

	if bypassLevel >= 1 {
		payload.Name = profile.BillingAddress.Name
	}

	if bypassLevel == 2 {
		payload.Address.City = profile.BillingAddress.Name
		payload.Address.Country = "US"
		payload.Address.Line1 = profile.BillingAddress.Line1
		payload.Address.PostalCode = profile.BillingAddress.PostCode 
	}

	p, err := json.Marshal(payload)
	if err != nil {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to checkout!"))
		return
	}

	fmt.Println(colors.TaskPrefix(id) + colors.Yellow("Attemping Bypass with level " + strconv.Itoa(bypassLevel) + "..."))
	fmt.Println(colors.TaskPrefix(id) + colors.Yellow("Submitting Checkout..."))

	TLUrl := "https://button-backend.tldash.ai/api/register/" + site + "/" + password
	req, err := http.NewRequest("POST", TLUrl, bytes.NewBuffer(p))
	if err != nil {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to checkout!"))
		return
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.93 Safari/537.36")
	req.Header.Set("Content-Type", "application/json")
	if discordSession != "" {
		req.Header.Set("authorization", "Bearer "+discordSession)
	}

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
		fmt.Println(colors.TaskPrefix(id) + colors.Green("THE BYPASS WORKED!! PLEASE GIVE Chrigeeel#9456 A KISS AND POST #SUCCESS"))
		fmt.Println(colors.TaskPrefix(id) + colors.Green("Successfully checked out on profile ") + colors.White("\"") + colors.Green(profile.Name) + colors.White("\""))
		go utility.SendWebhook(userData.Webhook, utility.WebhookContentStruct{
			Speed:   diff.String(),
			Module:  "TL Dashboards - Bypass",
			Site:    site,
			Profile: profile.Name,
		})
		payload, _ := json.Marshal(map[string]string{
			"site":     task.Site,
			"module":   "TL Dash",
			"speed":    diff.String(),
			"mode":     "Bypass",
			"password": "Unknown",
			"user":     userData.Username,
		})
		req, _ := http.NewRequest("POST", "http://ec2-13-52-240-112.us-west-1.compute.amazonaws.com:3000/checkouts", bytes.NewBuffer(payload))
		req.Header.Set("content-type", "application/json")
		client.Do(req)
		return
	}
	if postResponse.Error.Message != "" {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to checkout with reason: ") + colors.White("\"") + colors.Red(postResponse.Error.Message) + colors.White("\"") + colors.Red("!"))
		return
	}
	if postResponse.Message !=  "" {
		if postResponse.Message == "Sorry! Sold out!" {
			fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed executing Bypass!"))
			return
		}
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to checkout with reason: ") + colors.White("\"") + colors.Red(postResponse.Message) + colors.White("\"") + colors.Red("!"))
		return
	}
	fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to checkout for some reason, DM the below to owners:"))
	fmt.Println(colors.TaskPrefix(id) + colors.White(string(body)))

}

func newstripe(site string, stripe string) {
	req, err := http.NewRequest("POST", "https://hardcore.astolfoporn.com/newstripe?site=" + site + "&stripe=" + stripe, nil)
	if err != nil {
		return
	}
	req.Header.Set("x-auth", "BruhPlsStop")
	client := http.DefaultClient
	client.Do(req)
}