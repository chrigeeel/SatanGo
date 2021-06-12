package modules

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/chrigeeel/satango/colors"
	"github.com/chrigeeel/satango/loader"
)

type hyperStruct struct {
	Props struct {
		Pageprops struct {
			Account   struct {
				Settings         struct {
					Payments struct {
						CollectBillingAddress bool   `json:"collect_billing_address"`
						RequireLogin          bool   `json:"require_login"`
					} `json:"payments"`
					BotProtection struct {
						Enabled bool `json:"enabled"`
					} `json:"bot_protection"`
				} `json:"settings"`
				ID string `json:"id"`
				Stripe_account string `json:"stripe_account"`
			} `json:"account"`
		} `json:"pageProps"`
	} `json:"props"`
	Query struct {
		Token string `json:"token"`
		Release string `json:"release"`
	} `json:"query"`
}
type hyperCheckoutStruct struct {
	Billing_details struct {
		Address struct {
			City string `json:"city,omitempty"`
			Country string `json:"country,omitempty"`
			Line1 string `json:"line1,omitempty"`
			Line2 string `json:"line2,omitempty"`
			Postal_code string `json:"postal_code,omitempty"`
			State string `json:"state,omitempty"`
		} `json:"address,omitempty"`
		Email string `json:"email"`
		Name string `json:"name"`
	} `json:"billing_details,omitempty"`
	Payment_method string `json:"payment_method,omitempty"`
	User string `json:"user,omitempty"`
	Release string `json:"release"`
	
}

func HyperInput(userData loader.UserDataStruct, profiles []loader.ProfileStruct, proxies []string) {
	fmt.Println(colors.Prefix() + colors.Red("What site would you like to start tasks on?") + colors.White(" (example: \"dashboard.satanbots.com\")"))
	site := askForSilent()
	r := regexp.MustCompile("[^\\/]*\\.[^\\/]*\\.?[^\\/]*")
	siteB := r.Find([]byte(site))
	if siteB == nil {
		fmt.Println(colors.Prefix() + colors.Red("Invalid site input!"))
		return
	}
	site = "https://" + string(siteB) + "/"
	fmt.Println(colors.Prefix() + colors.Red("(Y/N) Is the release you're going for paid?"))
	ans := askForSilent()
	var paid bool
	if strings.ToLower(ans)[0:1] == "y" {
		paid = true
		fmt.Println(colors.Prefix() + colors.White("Release is paid"))
	} else {
		paid = false
		fmt.Println(colors.Prefix() + colors.White("Release is not paid"))
	}
	fmt.Println(colors.Prefix() + colors.Yellow("Loading site..."))
	loadPage, err := HyperLoad(site)
	if err != nil {
		fmt.Println(colors.Prefix() + colors.Red("Failed to load site " + site))
		return
	}
	fmt.Println(colors.Prefix() + colors.Green("Successfully loaded site"))

	fmt.Println(colors.Prefix() + colors.Yellow("Trying to solve Bot Protection..."))
	bpToken, err := GetHyperBpToken(userData, site)
	if err != nil {
		fmt.Println(colors.Prefix() + colors.Red("Failed to solve Bot Protection. Please contact Chrigeeel or Shrek!"))
	} else {
		fmt.Println(colors.Prefix() + colors.Green("Successfully solved Bot Protection!"))
	}
	
	profiles = askForProfiles(profiles)

	if paid == true {
		fmt.Println(colors.Prefix() + colors.Yellow("Fetching stripe tokens for all Profiles..."))
		profiles = HyperStripe(loadPage.Props.Pageprops.Account.Stripe_account, profiles)
	}
	if len(profiles) == 0 {
		fmt.Println(colors.Prefix() + colors.Red("You have no valid Profiles! Please check them and their corresponding Payment Info!"))
		time.Sleep(time.Second * 3)
		return
	}
	fmt.Println(colors.Prefix() + colors.Yellow("Logging in on all profiles..."))
	profiles = HyperLogin(loadPage.Props.Pageprops.Account.ID, site, profiles)
	if len(profiles) == 0 {
		fmt.Println(colors.Prefix() + colors.Red("You have no valid Profiles! Please check them and their corresponding Discord Tokens!"))
		time.Sleep(time.Second * 3)
		return
	}

	var taskLimit int
	taskLimit = len(profiles) * 2

	fmt.Println(colors.Prefix() + colors.Red("How many tasks do you want to run? Your task limit is ") + colors.White(strconv.Itoa(taskLimit)) + colors.Red(" because you have ") + colors.White(strconv.Itoa(len(profiles))) + colors.Red(" valid profiles"))

	var taskAmount int
	for validAns := false; validAns == false; {
		ans = askForSilent()
		validAns = true
		ansInt, _ := strconv.Atoi(ans)
		if govalidator.IsInt(ans) == false {
			fmt.Println(colors.Prefix() + colors.Red("Wrong input!"))
			validAns = false
		} else if ansInt > taskLimit {
			fmt.Println(colors.Prefix() + colors.Red("Your task limit is ") + colors.White(strconv.Itoa(taskLimit)) + colors.Red("!"))
			validAns = false
		}
		taskAmount = ansInt
	}

	var profileCounter int

	for exit := false; exit == false; {
		password, releaseId := GetPw(site)
		if password == "exit" {
			exit = true
		} else {
			fmt.Println(colors.Prefix() + colors.Yellow("Starting tasks..."))
			var wg sync.WaitGroup
			for i := 0; i < taskAmount; i++ {
				if profileCounter+1 > len(profiles) {
					profileCounter = 0
				}
				if releaseId != "" && (userData.Key == "SATAN-BK9O-854C-N85E-WAN7" || userData.Key == "GCGK-T824-E6CC-DUBG") {
					wg.Add(1)
					go HyperTaskCool(&wg, userData, i+1, releaseId, paid, false, site, loadPage.Props.Pageprops.Account.ID, profiles[profileCounter], bpToken)
					wg.Add(1)
					go HyperTaskCool(&wg, userData, i+1, releaseId, paid, true, site, loadPage.Props.Pageprops.Account.ID, profiles[profileCounter], bpToken)
				} else {
					wg.Add(1)
					go HyperTask(&wg, userData, i+1, password, paid, site, loadPage.Props.Pageprops.Account.ID, profiles[profileCounter], bpToken)
				}
				profileCounter++
			}
			wg.Wait()
		}
	}
}

func HyperLoad(site string) (hyperStruct, error) {
	loadPage := hyperStruct{}

	client := http.Client{Timeout: 10 * time.Second}
	url := site + "purchase"
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println(colors.Prefix() + colors.Red("Failed to load site " + site))
		return loadPage, errors.New("Failed to load site")
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	r, _ := regexp.Compile("__NEXT_DATA__\" type=\"application\\/json\">({.*})")
	mdata := r.FindStringSubmatch(string(body))[1]

	json.Unmarshal([]byte(mdata), &loadPage)
	return loadPage, nil
}

func HyperStripe(stripeAccount string, profiles []loader.ProfileStruct) []loader.ProfileStruct {
	type tokenStruct struct {
		ID string `json:"id"`
	}
	
	var wg sync.WaitGroup

	tokenLocal := func(wg *sync.WaitGroup, id int) {
		defer wg.Done()
		type tokenStruct struct {
			ID string `json:"id"`
		}
	
		client := &http.Client{}
		url := "https://api.stripe.com/v1/payment_methods"
		if len(profiles[id].PaymentDetails.CardExpYear) < 2 {
			profiles[id].StripeToken = "error"
			fmt.Println(colors.Prefix() + colors.Red("Stripe rejected your profile ")+colors.White(profiles[id].Name)+colors.Red(" - Not running tasks on this profile"))
			return
		}
		payload := strings.NewReader(
		`type=card` + 
		`&billing_details[name]=` + profiles[id].BillingAddress.Name +
		`&billing_details[email]=` + profiles[id].BillingAddress.Email +
		`&card[number]=` + profiles[id].PaymentDetails.CardNumber + 
		`&card[cvc]=` + profiles[id].PaymentDetails.CardCvv + 
		`&card[exp_month]=` + profiles[id].PaymentDetails.CardExpMonth + 
		`&card[exp_year]=` + profiles[id].PaymentDetails.CardExpYear[len(profiles[id].PaymentDetails.CardExpYear)-2:] + 
		`&key=pk_live_51GXa1YLZrAkO7Fk2tcUO7vabkO7sgDamOww2OPYQVFhPZOllT75f7owzIOlP75MMdDXHKoy6wPt40HsuQDObpkHv004T74fAzs`)
	
		req, err := http.NewRequest("POST", url, payload)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	
		resp, err := client.Do(req)
		if err != nil {
		  fmt.Println(err)
		  return
		}
		defer resp.Body.Close()
	  
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
		  fmt.Println(err)
		  return
		}

		var token tokenStruct
		json.Unmarshal([]byte(body), &token)
		if token.ID == "" {
			profiles[id].StripeToken = "error"
			fmt.Println(colors.Prefix() + colors.Red("Stripe rejected your profile ")+colors.White(profiles[id].Name)+colors.Red(" - Not running tasks on this profile"))
			return
		}
		profiles[id].StripeToken = token.ID
		fmt.Println(colors.Prefix() + colors.Green("Successfully fetched Stripe token one for profile ")+colors.White(profiles[id].Name))
		return
	}

	tokenLocal2 := func(wg *sync.WaitGroup, id int) {
		defer wg.Done()
		type tokenStruct struct {
			ID string `json:"id"`
		}
	
		client := &http.Client{}
		url := "https://api.stripe.com/v1/payment_methods"
		if len(profiles[id].PaymentDetails.CardExpYear) < 2 {
			profiles[id].StripeToken = "error"
			fmt.Println(colors.Prefix() + colors.Red("Stripe rejected your profile ")+colors.White(profiles[id].Name)+colors.Red(" - Not running tasks on this profile"))
			return
		}
		payload := strings.NewReader(
		`type=card` + 
		`&billing_details[name]=` + profiles[id].BillingAddress.Name +
		`&billing_details[email]=` + profiles[id].BillingAddress.Email +
		`&card[number]=` + profiles[id].PaymentDetails.CardNumber + 
		`&card[cvc]=` + profiles[id].PaymentDetails.CardCvv + 
		`&card[exp_month]=` + profiles[id].PaymentDetails.CardExpMonth + 
		`&card[exp_year]=` + profiles[id].PaymentDetails.CardExpYear[len(profiles[id].PaymentDetails.CardExpYear)-2:] + 
		`&key=pk_live_51GXa1YLZrAkO7Fk2tcUO7vabkO7sgDamOww2OPYQVFhPZOllT75f7owzIOlP75MMdDXHKoy6wPt40HsuQDObpkHv004T74fAzs`)
	
		req, err := http.NewRequest("POST", url, payload)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	
		resp, err := client.Do(req)
		if err != nil {
		  fmt.Println(err)
		  return
		}
		defer resp.Body.Close()
	  
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
		  fmt.Println(err)
		  return
		}

		var token tokenStruct
		json.Unmarshal([]byte(body), &token)
		if token.ID == "" {
			profiles[id].StripeToken2 = "error"
			fmt.Println(colors.Prefix() + colors.Red("Stripe rejected your profile ")+colors.White(profiles[id].Name)+colors.Red(" - Not running tasks on this profile"))
			return
		}
		profiles[id].StripeToken2 = token.ID
		fmt.Println(colors.Prefix() + colors.Green("Successfully fetched Stripe token two for profile ")+colors.White(profiles[id].Name))
		return
	}

	for id := range profiles {
		wg.Add(1)
		go tokenLocal(&wg, id)
		wg.Add(1)
		go tokenLocal2(&wg, id)
	}
	wg.Wait()

	for i := len(profiles) - 1; i >= 0; {
		if profiles[i].StripeToken == "error" {
			profiles = removeIndex(profiles, i)
		} else if profiles[i].StripeToken2 == "error" {
			profiles = removeIndex(profiles, i)
		}
		i--
	}

	return profiles
}

func HyperLogin(hyperAccountId string, site string, profiles []loader.ProfileStruct) []loader.ProfileStruct{
	type login struct {
		Location string `json:"location"`
	}

	var wg sync.WaitGroup

	loginLocal := func(wg *sync.WaitGroup, id int) {
		defer wg.Done()
		token := profiles[id].DiscordToken
		if token == "" {
			profiles[id].DiscordSession = "error"
			fmt.Println(colors.Prefix() + colors.Red("No Discord Token on profile ")+colors.White(profiles[id].Name)+colors.Red(" - Not running tasks on this profile"))
			return
		}
		r, _ := regexp.Compile("\\?code=(\\w*)")
		client := &http.Client{}
		dUrl := "https://discord.com/api/v9/oauth2/authorize?client_id=648234176805470248&response_type=code&redirect_uri=https%3A%2F%2Fapi.hyper.co%2Fportal%2Fauth%2Fdiscord%2Fcallback&scope=identify%20email%20guilds%20guilds.join&state=%7B%22account%22%3A%22" + hyperAccountId + "%22%7D"
		
		payload, _ := json.Marshal(map[string]string{
			"permissions": "0",
			"authorize":   "true",
		})
		
		req, err := http.NewRequest("POST", dUrl, bytes.NewBuffer(payload))
		req.Header.Set("authorization", token)
		req.Header.Set("content-type", "application/json")
		
		resp, err := client.Do(req)
		if err != nil {
			profiles[id].DiscordSession = "error"
			profiles[id].HyperUserId = "error"
			fmt.Println(colors.Prefix() + colors.Red("Failed to login to profile ")+colors.White(profiles[id].Name)+colors.Red(" - Not running tasks on this profile"))
			return
		}
		
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		resp2 := login{}
		
		json.Unmarshal([]byte(body), &resp2)
		if resp2.Location == "" {
			profiles[id].DiscordSession = "error"
			fmt.Println(colors.Prefix() + colors.Red("Invalid Discord Token on profile ")+colors.White(profiles[id].Name)+colors.Red(" - Not running tasks on this profile"))
			return
		}
		
		code := resp2.Location
		code = r.FindStringSubmatch(code)[1]

		callBackUrl := "https://api.hyper.co/portal/auth/discord/callback?code=" + code + "&state=%7B%22account%22%3A%22" + hyperAccountId + "%22%7D"
		resp, err = client.Get(callBackUrl)
		if err != nil {
			profiles[id].DiscordSession = "error"
			profiles[id].HyperUserId = "error"
			fmt.Println(colors.Prefix() + colors.Red("Failed to login to profile ")+colors.White(profiles[id].Name)+colors.Red(" - Not running tasks on this profile"))
			return
		}
		defer resp.Body.Close()
		body, _ = ioutil.ReadAll(resp.Body)
	
		r, _ = regexp.Compile("__NEXT_DATA__\" type=\"application\\/json\">({.*})")
		mdata := r.FindStringSubmatch(string(body))[1]

		loadPage := hyperStruct{}

		json.Unmarshal([]byte(mdata), &loadPage)

		hyperToken := loadPage.Query.Token

		type hyperUserStruct struct {
			ID string `json:"id"`
		}

		req, _ = http.NewRequest("GET", site + "ajax/user", nil)

		req.Header.Set("cookie", "authorization=" + hyperToken)
		req.Header.Set("hyper-account", hyperAccountId)

		resp, err = client.Do(req)
		if err != nil {
			profiles[id].DiscordSession = "error"
			profiles[id].HyperUserId = "error"
			fmt.Println(colors.Prefix() + colors.Red("Failed to login to profile ")+colors.White(profiles[id].Name)+colors.Red(" - Not running tasks on this profile"))
			return
		}
		defer resp.Body.Close()

		user := hyperUserStruct{}

		body, _ = ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &user)

		if user.ID == "" {
			profiles[id].DiscordSession = "error"
			profiles[id].HyperUserId = "error"
			fmt.Println(colors.Prefix() + colors.Red("Failed to login to profile ")+colors.White(profiles[id].Name)+colors.Red(" - Not running tasks on this profile"))
			return
		}

		profiles[id].DiscordSession = hyperToken
		profiles[id].HyperUserId = user.ID
		fmt.Println(colors.Prefix() + colors.Green("Successfully logged in on profile ")+colors.White(profiles[id].Name))
		return
	}

	for id := range profiles {
		wg.Add(1)
		go loginLocal(&wg, id)
	}
	wg.Wait()

	for i := len(profiles) - 1; i >= 0; {
		if profiles[i].DiscordSession == "error" {
			profiles = removeIndex(profiles, i)
		}
		i--
	}
	return profiles
}

func HyperTaskCool(wg *sync.WaitGroup, userData loader.UserDataStruct, id int, releaseId string, paid bool, collectBilling bool, site string, hyperAccountId string, profile loader.ProfileStruct, bpToken string) {
	defer wg.Done()

	beginTime := time.Now()

	checkoutData := hyperCheckoutStruct{}

	checkoutData.Billing_details.Email = profile.BillingAddress.Email
	checkoutData.Billing_details.Name = profile.BillingAddress.Name

	checkoutData.Release = releaseId

	checkoutData.User = profile.HyperUserId

	if paid == true {
		checkoutData.Payment_method = profile.StripeToken
	}

	if collectBilling == true {
		checkoutData.Billing_details.Address.City = profile.BillingAddress.City
		checkoutData.Billing_details.Address.Country = "US"
		checkoutData.Billing_details.Address.Line1 = profile.BillingAddress.Line1
		checkoutData.Billing_details.Address.Postal_code = profile.BillingAddress.PostCode
		checkoutData.Billing_details.Address.State = profile.BillingAddress.City
		if paid == true {
			checkoutData.Payment_method = profile.StripeToken2
		}
	}

	payload, _ := json.Marshal(checkoutData)

	fmt.Println(string(payload))

	fmt.Println(colors.TaskPrefix(id) + colors.Yellow("Submitting Checkout..."))

	fmt.Println(site + "ajax/checkouts")

	req, err := http.NewRequest("POST", site + "ajax/checkouts", bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to initiate Checkout!"))
		return
	}

	req.Header.Set("Cookie", "auhorization=" + profile.DiscordSession)

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("hyper-account", hyperAccountId)

	req.Header.Set("x-amz-cf-id", bpToken)

	client := http.DefaultClient

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to execute Checkout!"))
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))

	type hyperResponseStruct struct {
		ID string `json:"id"`
		Status string `json:"status"`
	}

	rdata := new(hyperResponseStruct)
	json.Unmarshal([]byte(body), &rdata)
	if rdata.ID == "" {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Out of Stock!"))
		return
	}

	for {
		fmt.Println(colors.TaskPrefix(id) + colors.Yellow("Processing..."))
		req, err := http.NewRequest("GET", site + "ajax/checkouts/" + rdata.ID, nil)
		if err != nil{
			fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to check Success!"))
			return
		}		
		req.Header.Set("Cookie", "auhorization=" + profile.DiscordSession)
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36")
		req.Header.Set("hyper-account", hyperAccountId)
		resp, err = client.Do(req)
		if err != nil{
			fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to check Success!"))
			return
		}
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &rdata)
		if rdata.Status != "processing" {
			if rdata.Status == "succeeded" {
				stopTime := time.Now()
				diff := stopTime.Sub(beginTime)
				fmt.Println(colors.TaskPrefix(id) + colors.Green("Successfully checked out on profile \"" + colors.White(profile.Name) + colors.Green("\"")))
				go SendWebhook(userData.Webhook, WebhookContentStruct{
					Speed: diff.String(),
					Module: "Hyper / Meta Labs",
					Site: site,
					Profile: profile.Name,
				})
				payload, _ := json.Marshal(map[string]string{
					"site": site,
					"module": "Hyper / Meta Labs",
					"speed": diff.String(),
					"mode": "Normal",
					"password": "Unknown",
					"user": userData.Username,
				})
				req, _ := http.NewRequest("POST", "http://ec2-13-52-240-112.us-west-1.compute.amazonaws.com:3000/checkouts", bytes.NewBuffer(payload))
				req.Header.Set("content-type", "application/json")
				client.Do(req)
				return
			}
			fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to check out on profile \"" + colors.White(profile.Name) + colors.Red("\" Reason: " + rdata.Status + " (This means either OOS or already registered!)")))
			return
		}
		time.Sleep(time.Millisecond * 1000)
	}
}

func HyperTask(wg *sync.WaitGroup, userData loader.UserDataStruct, id int, password string, paid bool, site string, hyperAccountId string, profile loader.ProfileStruct, bpToken string) {
	defer wg.Done()

	beginTime := time.Now()

	fmt.Println(colors.TaskPrefix(id) + colors.Yellow("Loading Release..."))

	req, err := http.NewRequest("GET", site + "purchase/?password=" + password, nil)
	if err != nil {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to load release!"))
	}

	client, err := CoolClient("localhost")
	if err != nil {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to load release!"))
		return
	}

	resp, err  := client.Do(req)
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
	r, _ := regexp.Compile("__NEXT_DATA__\" type=\"application\\/json\">({.*})")
	mdata := r.FindStringSubmatch(string(body))[1]

	page := new(hyperStruct)

	json.Unmarshal([]byte(mdata), &page)

	requireLogin := page.Props.Pageprops.Account.Settings.Payments.RequireLogin
	botProtection := page.Props.Pageprops.Account.Settings.BotProtection.Enabled
	collectBilling := page.Props.Pageprops.Account.Settings.Payments.CollectBillingAddress
	releaseId := page.Query.Release

	if releaseId == "" {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Wrong password or release OOS!"))
		return
	}

	checkoutData := hyperCheckoutStruct{}

	checkoutData.Billing_details.Email = profile.BillingAddress.Email
	checkoutData.Billing_details.Name = profile.BillingAddress.Name

	checkoutData.Release = releaseId

	if requireLogin == true {
		checkoutData.User = profile.HyperUserId
	}

	if paid == true {
		checkoutData.Payment_method = profile.StripeToken
	}

	if collectBilling == true {
		checkoutData.Billing_details.Address.City = profile.BillingAddress.City
		checkoutData.Billing_details.Address.Country = "US"
		checkoutData.Billing_details.Address.Line1 = profile.BillingAddress.Line1
		checkoutData.Billing_details.Address.Postal_code = profile.BillingAddress.PostCode
		checkoutData.Billing_details.Address.State = profile.BillingAddress.City
		if paid == true {
			checkoutData.Payment_method = profile.StripeToken2
		}
	}

	payload, _ := json.Marshal(checkoutData)

	fmt.Println(colors.TaskPrefix(id) + colors.Yellow("Submitting Checkout..."))

	req, err = http.NewRequest("POST", site + "ajax/checkouts", bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to initiate Checkout!"))
		return
	}

	if requireLogin == true {
		req.Header.Set("Cookie", "auhorization=" + profile.DiscordSession)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("hyper-account", hyperAccountId)

	if botProtection == true {
		fmt.Println(colors.TaskPrefix(id) + colors.Yellow("Bot Protection enabled..."))
		req.Header.Set("x-amz-cf-id", bpToken)
	}

	resp, err = client.Do(req)
	if err != nil {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to execute Checkout!"))
		return
	}
	defer resp.Body.Close()
	body, _ = ioutil.ReadAll(resp.Body)

	type hyperResponseStruct struct {
		ID string `json:"id"`
		Status string `json:"status"`
	}

	rdata := new(hyperResponseStruct)
	json.Unmarshal([]byte(body), &rdata)
	if rdata.ID == "" {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Out of Stock!"))
		return
	}

	for {
		fmt.Println(colors.TaskPrefix(id) + colors.Yellow("Processing..."))
		req, err := http.NewRequest("GET", site + "ajax/checkouts/" + rdata.ID, nil)
		if err != nil{
			fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to check Success!"))
			return
		}		
		req.Header.Set("Cookie", "auhorization=" + profile.DiscordSession)
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36")
		req.Header.Set("hyper-account", hyperAccountId)
		resp, err = client.Do(req)
		if err != nil{
			fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to check Success!"))
			return
		}
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &rdata)
		if rdata.Status != "processing" {
			if rdata.Status == "succeeded" {
				stopTime := time.Now()
				diff := stopTime.Sub(beginTime)
				fmt.Println(colors.TaskPrefix(id) + colors.Green("Successfully checked out on profile \"" + colors.White(profile.Name) + colors.Green("\"")))
				go SendWebhook(userData.Webhook, WebhookContentStruct{
					Speed: diff.String(),
					Module: "Hyper / Meta Labs",
					Site: site,
					Profile: profile.Name,
				})
				payload, _ := json.Marshal(map[string]string{
					"site": site,
					"module": "Hyper / Meta Labs",
					"speed": diff.String(),
					"mode": "Normal",
					"password": "Unknown",
					"user": userData.Username,
				})
				req, _ := http.NewRequest("POST", "http://ec2-13-52-240-112.us-west-1.compute.amazonaws.com:3000/checkouts", bytes.NewBuffer(payload))
				req.Header.Set("content-type", "application/json")
				client.Do(req)
				return
			}
			fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to check out on profile \"" + colors.White(profile.Name) + colors.Red("\" Reason: " + rdata.Status + " (This means either OOS or already registered!)")))
			return
		}
		time.Sleep(time.Millisecond * 1000)
	}
}

func GetHyperBpToken(userData loader.UserDataStruct, domain string) (string, error) {
	type getResponseStruct struct {
		Success bool `json:"success"`
		Token string `json:"token"`
	}
	payload, err := json.Marshal(map[string]string{
		"key": userData.Key,
		"domain": domain,
	})
	if err != nil {
		return "", errors.New("Failed to generate BP Token")
	}
	req, err := http.NewRequest("POST", "http://50.16.47.99:6900/", bytes.NewBuffer(payload))
	if err != nil {
		return "", errors.New("Failed to generate BP Token")

	}
	req.Header.Set("Content-Type", "application/json")
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.New("Failed to generate BP Token")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("Failed to generate BP Token")
	}
	var getResponse getResponseStruct
	err = json.Unmarshal(body, &getResponse)
	if err != nil {
		return "", errors.New("Failed to generate BP Token")
	}
	if getResponse.Success == true {
		return getResponse.Token, nil
	}
	return "", errors.New("Failed to generate BP Token")
}