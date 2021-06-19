package tldash

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/chrigeeel/satango/colors"
	"github.com/chrigeeel/satango/loader"
	"github.com/chrigeeel/satango/modules/getpw"
	"github.com/chrigeeel/satango/modules/utility"
)

func taskfcfs(wg *sync.WaitGroup, userData loader.UserDataStruct, id int, password string, solveIp string, task taskStruct) {
	type getResponseStruct struct {
		Stripe_public_key string `json:"stripe_public_key"`
		Price_with_symbol string `json:"price_with_symbol"`
		Captcha           string `json:"captcha"`
	}
	type tlError struct {
		Message string `json:"message"`
	}
	type postResponseStruct struct {
		Success bool    `json:"success"`
		Message string  `json:"message"`
		Error   tlError `json:"error"`
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

	var getResponse getResponseStruct
	json.Unmarshal([]byte(body), &getResponse)

	if string(getResponse.Stripe_public_key) == "" {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Wrong password or release OOS!"))
		return
	}

	var cfCookie string

	for _, cookie := range resp.Cookies() {
		if cookie.Name == "__cf_bm" {
			cfCookie = cookie.Value
		}
	}
	if cfCookie != "" {
		go getpw.PWSharingSend(userData, password, site)
	}

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

	var captchaSolution string

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

	var payload []byte

	if captchaSolution == "" {
		payload, err = json.Marshal(map[string]string{
			"email": profile.BillingAddress.Email,
			"token": stripeToken,
		})
		if err != nil {
			fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to checkout!"))
		}
	} else {
		payload, err = json.Marshal(map[string]string{
			"captcha": captchaSolution,
			"email":   profile.BillingAddress.Email,
			"token":   stripeToken,
		})
		if err != nil {
			fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to checkout!"))
		}
	}

	TLUrl = "https://button-backend.tldash.ai/api/register/" + site + "/" + password
	req, err = http.NewRequest("POST", TLUrl, bytes.NewBuffer(payload))
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
		fmt.Println(colors.TaskPrefix(id) + colors.Green("Successfully checked out, Check your email!"))
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
	fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to checkout for some reason, DM the below to owners:"))
	fmt.Println(colors.TaskPrefix(id) + colors.White(string(body)))
}