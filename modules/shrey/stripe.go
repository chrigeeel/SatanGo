package shrey

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"github.com/chrigeeel/satango/colors"
	"github.com/chrigeeel/satango/loader"
	"github.com/chrigeeel/satango/utility"
)

func stripe(profiles []loader.ProfileStruct, site string) []loader.ProfileStruct {

	var wg sync.WaitGroup

	tokenLocal := func(wg *sync.WaitGroup, id int) {
		defer wg.Done()
		type tokenStruct struct {
			ID string `json:"id"`
		}

		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
		url := site + "dashboard"
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			profiles[id].StripeToken = "error"
			fmt.Println(colors.Prefix() + colors.Red("Stripe rejected your profile ") + colors.White(profiles[id].Name) + colors.Red(" - Not running tasks on this profile"))
			return
		}
		req.Header.Add("Upgrade-Insecure-Requests", "1")
		req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
		req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
		req.Header.Add("Sec-Fetch-Site", "cross-site")
		req.Header.Add("Sec-Fetch-Mode", "navigate")
		req.Header.Add("Sec-Fetch-User", "?1")
		req.Header.Add("Sec-Fetch-Dest", "document")
		req.Header.Add("sec-ch-ua", "\" Not;A Brand\";v=\"99\", \"Google Chrome\";v=\"91\", \"Chromium\";v=\"91\"")
		req.Header.Add("sec-ch-ua-mobile", "?0")
		req.Header.Add("Cookie", "_shreyauth_session="+profiles[id].DiscordSession)
		resp, err := client.Do(req)
		if err != nil {
			profiles[id].StripeToken = "error"
			fmt.Println(colors.Prefix() + colors.Red("Stripe rejected your profile ") + colors.White(profiles[id].Name) + colors.Red(" - Not running tasks on this profile"))
			return
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			profiles[id].StripeToken = "error"
			fmt.Println(colors.Prefix() + colors.Red("Stripe rejected your profile ") + colors.White(profiles[id].Name) + colors.Red(" - Not running tasks on this profile"))
			return
		}
		for _, cookie := range resp.Cookies() {
			if cookie.Name == "_shreyauth_session" {
				profiles[id].DiscordSession = cookie.Value
			}
		}
		r := regexp.MustCompile(`csrf-token" content="([^"]*)`)
		match := r.FindStringSubmatch(string(body))
		if len(match) != 2 {
			profiles[id].StripeToken = "error"
			fmt.Println(colors.Prefix() + colors.Red("Stripe rejected your profile ") + colors.White(profiles[id].Name) + colors.Red(" - Not running tasks on this profile"))
			return
		}
		csrfToken := match[1]
		r = regexp.MustCompile(`Stripe\('([^']*)`)
		match = r.FindStringSubmatch(string(body))
		if len(match) != 2 {
			profiles[id].StripeToken = "error"
			fmt.Println(colors.Prefix() + colors.Red("Stripe rejected your profile ") + colors.White(profiles[id].Name) + colors.Red(" - Not running tasks on this profile"))
			return
		}
		stripeKey := match[1]

		url = site + "payment_methods"
		req, err = http.NewRequest("POST", url, nil)
		if err != nil {
			profiles[id].StripeToken = "error"
			fmt.Println(colors.Prefix() + colors.Red("Stripe rejected your profile ") + colors.White(profiles[id].Name) + colors.Red(" - Not running tasks on this profile"))
			return
		}
		req.Header.Add("sec-ch-ua", `\" Not;A Brand\";v=\"99\", \"Google Chrome\";v=\"91\", \"Chromium\";v=\"91\"`)
		req.Header.Add("Accept", "text/javascript, application/javascript, application/ecmascript, application/x-ecmascript, */*; q=0.01")
		req.Header.Add("X-CSRF-Token", csrfToken)
		req.Header.Add("X-Requested-With", "XMLHttpRequest")
		req.Header.Add("sec-ch-ua-mobile", "?0")
		req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
		req.Header.Add("Sec-Fetch-Site", "cross-site")
		req.Header.Add("Sec-Fetch-Mode", "cors")
		req.Header.Add("Sec-Fetch-Dest", "empty")
		req.Header.Add("Cookie", "_shreyauth_session=" + profiles[id].DiscordSession)
		resp, err = client.Do(req)
		if err != nil {
			profiles[id].StripeToken = "error"
			fmt.Println(colors.Prefix() + colors.Red("Stripe rejected your profile ") + colors.White(profiles[id].Name) + colors.Red(" - Not running tasks on this profile"))
			return
		}
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			profiles[id].StripeToken = "error"
			fmt.Println(colors.Prefix() + colors.Red("Stripe rejected your profile ") + colors.White(profiles[id].Name) + colors.Red(" - Not running tasks on this profile"))
			return
		}
		r = regexp.MustCompile(`sessionId: '([^']*)`)
		match = r.FindStringSubmatch(string(body))
		if len(match) != 2 {
			profiles[id].StripeToken = "error"
			fmt.Println(colors.Prefix() + colors.Red("Stripe rejected your profile ") + colors.White(profiles[id].Name) + colors.Red(" - Not running tasks on this profile"))
			return
		}
		stripeSessionId := match[1]

		url = "https://api.stripe.com/v1/payment_methods"

		if len(profiles[id].PaymentDetails.CardExpYear) < 2 {
			profiles[id].StripeToken = "error"
			fmt.Println(colors.Prefix() + colors.Red("Stripe rejected your profile ") + colors.White(profiles[id].Name) + colors.Red(" - Not running tasks on this profile"))
			return
		}
		payload := strings.NewReader(
			`type=card` +
			`&card[number]=` + profiles[id].PaymentDetails.CardNumber +
			`&card[cvc]=` + profiles[id].PaymentDetails.CardCvv +
			`&card[exp_month]=` + profiles[id].PaymentDetails.CardExpMonth +
			`&card[exp_year]=` + profiles[id].PaymentDetails.CardExpYear[len(profiles[id].PaymentDetails.CardExpYear)-2:] +
			`&billing_details[name]=` + profiles[id].BillingAddress.Name +
			`&billing_details[email]=` + profiles[id].BillingAddress.Email +
			`&billing_details[address][country]=` + "US" +
			`&billing_details[address][line1]=` + profiles[id].BillingAddress.Line1 +
			`&billing_details[address][city]=` + profiles[id].BillingAddress.City +
			`&billing_details[address][postal_code]=` + profiles[id].BillingAddress.PostCode +
			`&key=` + stripeKey)

		req, err = http.NewRequest("POST", url, payload)
		if err != nil {
			profiles[id].StripeToken = "error"
			fmt.Println(colors.Prefix() + colors.Red("Stripe rejected your profile ") + colors.White(profiles[id].Name) + colors.Red(" - Not running tasks on this profile"))
			return
		}

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err = client.Do(req)
		if err != nil {
			profiles[id].StripeToken = "error"
			fmt.Println(colors.Prefix() + colors.Red("Stripe rejected your profile ") + colors.White(profiles[id].Name) + colors.Red(" - Not running tasks on this profile"))
			return
		}
		defer resp.Body.Close()

		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			profiles[id].StripeToken = "error"
			fmt.Println(colors.Prefix() + colors.Red("Stripe rejected your profile ") + colors.White(profiles[id].Name) + colors.Red(" - Not running tasks on this profile"))
			return
		}

		var token tokenStruct
		json.Unmarshal([]byte(body), &token)
		if token.ID == "" {
			profiles[id].StripeToken = "error"
			fmt.Println(colors.Prefix() + colors.Red("Stripe rejected your profile ") + colors.White(profiles[id].Name) + colors.Red(" - Not running tasks on this profile"))
			return
		}
		profiles[id].StripeToken = token.ID
		fmt.Println(colors.Prefix() + colors.Green("Successfully fetched Payment Method for profile ") + colors.White(profiles[id].Name))
		fmt.Println(colors.Prefix() + colors.Yellow("Attaching Payment Method to user..."))

		confirmUrl := "https://api.stripe.com/v1/payment_pages/" + stripeSessionId + "/confirm"
		payload = strings.NewReader(
			`eid=NA` +
			`&payment_method=` + profiles[id].StripeToken +
			`&key=` + stripeKey)
	
		req, err = http.NewRequest("POST", confirmUrl, payload)
		if err != nil {
			profiles[id].StripeToken = "error"
			fmt.Println(colors.Prefix() + colors.Red("Stripe rejected your profile ") + colors.White(profiles[id].Name) + colors.Red(" - Not running tasks on this profile"))
			return
		}
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
		var token2 tokenStruct
		json.Unmarshal([]byte(body), &token2)
		if token2.ID == "" {
			profiles[id].StripeToken = "error"
			fmt.Println(colors.Prefix() + colors.Red("Stripe rejected your profile ") + colors.White(profiles[id].Name) + colors.Red(" - Not running tasks on this profile"))
			return
		}

		url = site + "payment_methods?session_id=" + stripeSessionId
		req, err = http.NewRequest("GET", url, nil)
		if err != nil {
			profiles[id].StripeToken = "error"
			fmt.Println(colors.Prefix() + colors.Red("Stripe rejected your profile ") + colors.White(profiles[id].Name) + colors.Red(" - Not running tasks on this profile"))
			return
		}
		req.Header.Add("Upgrade-Insecure-Requests", "1")
		req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
		req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
		req.Header.Add("Sec-Fetch-Site", "cross-site")
		req.Header.Add("Sec-Fetch-Mode", "navigate")
		req.Header.Add("Sec-Fetch-User", "?1")
		req.Header.Add("Sec-Fetch-Dest", "document")
		req.Header.Add("sec-ch-ua", "\" Not;A Brand\";v=\"99\", \"Google Chrome\";v=\"91\", \"Chromium\";v=\"91\"")
		req.Header.Add("sec-ch-ua-mobile", "?0")
		req.Header.Add("Cookie", "_shreyauth_session="+profiles[id].DiscordSession)
		resp, err = client.Do(req)
		if err != nil {
			profiles[id].StripeToken = "error"
			fmt.Println(colors.Prefix() + colors.Red("Stripe rejected your profile ") + colors.White(profiles[id].Name) + colors.Red(" - Not running tasks on this profile"))
			return
		}
		for _, cookie := range resp.Cookies() {
			if cookie.Name == "_shreyauth_session" {
				profiles[id].DiscordSession = cookie.Value
			}
		}

		url = site + "subscribe"
		req, err = http.NewRequest("GET", url, nil)
		if err != nil {
			profiles[id].StripeToken = "error"
			fmt.Println(colors.Prefix() + colors.Red("Stripe rejected your profile ") + colors.White(profiles[id].Name) + colors.Red(" - Not running tasks on this profile"))
			return
		}
		req.Header.Add("Upgrade-Insecure-Requests", "1")
		req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
		req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
		req.Header.Add("Sec-Fetch-Site", "cross-site")
		req.Header.Add("Sec-Fetch-Mode", "navigate")
		req.Header.Add("Sec-Fetch-User", "?1")
		req.Header.Add("Sec-Fetch-Dest", "document")
		req.Header.Add("sec-ch-ua", "\" Not;A Brand\";v=\"99\", \"Google Chrome\";v=\"91\", \"Chromium\";v=\"91\"")
		req.Header.Add("sec-ch-ua-mobile", "?0")
		req.Header.Add("Cookie", "_shreyauth_session="+profiles[id].DiscordSession)
		resp, err = client.Do(req)
		if err != nil {
			profiles[id].StripeToken = "error"
			fmt.Println(colors.Prefix() + colors.Red("Stripe rejected your profile ") + colors.White(profiles[id].Name) + colors.Red(" - Not running tasks on this profile"))
			return
		}
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			profiles[id].StripeToken = "error"
			fmt.Println(colors.Prefix() + colors.Red("Stripe rejected your profile ") + colors.White(profiles[id].Name) + colors.Red(" - Not running tasks on this profile"))
			return
		}
		for _, cookie := range resp.Cookies() {
			if cookie.Name == "_shreyauth_session" {
				profiles[id].DiscordSession = cookie.Value
			}
		}
		r = regexp.MustCompile(`csrf-token" content="([^"]*)`)
		match = r.FindStringSubmatch(string(body))
		if len(match) != 2 {
			profiles[id].StripeToken = "error"
			fmt.Println(colors.Prefix() + colors.Red("Stripe rejected your profile ") + colors.White(profiles[id].Name) + colors.Red(" - Not running tasks on this profile"))
			return
		}
		csrfToken = match[1]
		profiles[id].ShreyCSRF = csrfToken

		fmt.Println(colors.Prefix() + colors.Green("Successfully attached Payment Method to profile ") + colors.White(profiles[id].Name) + colors.Green("!"))
	}


	for id := range profiles {
		wg.Add(1)
		go tokenLocal(&wg, id)
	}
	wg.Wait()


	for i := len(profiles) - 1; i >= 0; {
		if profiles[i].StripeToken == "error" {
			profiles = utility.RemoveIndex(profiles, i)
		} else if profiles[i].StripeToken2 == "error" {
			profiles = utility.RemoveIndex(profiles, i)
		}
		i--
	}

	return profiles
}