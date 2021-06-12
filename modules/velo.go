package modules

import (
	"bytes"
	"encoding/json"
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

var veloSites = []siteStruct{
	siteStruct{
		DisplayName: "NEST",
		BackendName: "dashboard.nestbycardinal.com",
		Stripe_public_key: "pk_live_51IA68dJ8wbUpNsnLYiHjewIyVRpWxRwHICvlxrB53GHnorqNvs5IL8H8BYgeW5U0Dl7CXzIs2l3el741uDkU2Gh000TH4f7mfE",
	},
	siteStruct{
		DisplayName: "CryptoClub",
		BackendName: "dash.cryptoclub.group",
		Stripe_public_key: "pk_live_51IibakGMt5G1CmqPqzOg6k2RKavZa2PF2BoQ0c5BM53GKANiJCGq7CZwL1uHAzjjcAYD4jaCYsYtsa6M3QV7ivMX00n8thUefS",
	},
}

func VeloInput(userData loader.UserDataStruct, profiles []loader.ProfileStruct, proxies []string) {
	fmt.Println(colors.Prefix() + colors.Red("What site would you like to start tasks on?"))
	for i := range veloSites {
		fmt.Println(colors.Prefix() + colors.White("["+strconv.Itoa(i+1)+"] "+veloSites[i].DisplayName))
	}
	fmt.Println(colors.Prefix() + colors.White("[%] Go back"))
	var ansInt int
	for validAns := false; validAns == false; {
		ans := askForSilent()
		if ans == "%" {
			return
		}
		validAns = true
		ansInt, _ = strconv.Atoi(ans)
		if govalidator.IsInt(ans) == false {
			fmt.Println(colors.Prefix() + colors.Red("Wrong input!"))
			validAns = false
		} else if ansInt > len(veloSites) {
			fmt.Println(colors.Prefix() + colors.Red("Wrong input!"))
			validAns = false
		}
	}
	site := veloSites[ansInt-1]
	fmt.Println(colors.Prefix() + colors.Yellow("Loggin in on all profiles..."))
	profiles = VeloLogin(profiles, site.BackendName)
	if len(profiles) == 0 {
		fmt.Println(colors.Prefix() + colors.Red("You have no valid Profiles! Please check them and their corresponding Discord Tokens!"))
		time.Sleep(time.Second * 3)
		return
	}
	fmt.Println(colors.Prefix() + colors.Yellow("Fetching Stripe on all profiles..."))
	profiles = VeloStripe(profiles, site.Stripe_public_key)
	if len(profiles) == 0 {
		fmt.Println(colors.Prefix() + colors.Red("You have no valid Profiles! Please check them and their corresponding Discord Tokens!"))
		time.Sleep(time.Second * 3)
		return
	}
	profiles = askForProfiles(profiles)

	var taskLimit int
	taskLimit = len(profiles) * 4

	fmt.Println(colors.Prefix() + colors.Red("How many tasks do you want to run? Your task limit is ") + colors.White(strconv.Itoa(taskLimit)) + colors.Red(" because you have ") + colors.White(strconv.Itoa(len(profiles))) + colors.Red(" valid profiles"))

	var taskAmount int
	for validAns := false; validAns == false; {
		ans := askForSilent()
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
		password, _ := GetPw(site.BackendName)
		if password == "exit" {
			exit = true
		} else {
			fmt.Println(colors.Prefix() + colors.Yellow("Starting tasks..."))
			var wg sync.WaitGroup
			for i := 0; i < taskAmount; i++ {
				if profileCounter+1 > len(profiles) {
					profileCounter = 0
				}
				wg.Add(1)
				go VeloTask(&wg, userData, i+1, site, password, profiles[profileCounter])
				profileCounter++
			}
			wg.Wait()
		}
	}
}

func VeloLogin(profiles []loader.ProfileStruct, site string) []loader.ProfileStruct {
	type login struct {
		Location string `json:"location"`
	}

	var wg sync.WaitGroup

	loginLocal := func(wg *sync.WaitGroup, id int) {
		defer wg.Done()
		token := profiles[id].DiscordToken
		if token == "" {
			profiles[id].DiscordSession = "error"
			fmt.Println(colors.Prefix() + colors.Red("No Discord Token on profile ") + colors.White(profiles[id].Name) + colors.Red(" - Not running tasks on this profile"))
			return
		}
		r, _ := regexp.Compile("\\?code=(\\w*)")
		client := &http.Client{
			Timeout: 7 * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
		dUrl := "https://discord.com/api/v9/oauth2/authorize?client_id=646888831596888085&response_type=code&redirect_uri=https%3A%2F%2Fvlo.to%2Flogin%2Fcomplete&scope=identify%20email%20guilds.join%20guilds"

		payload, _ := json.Marshal(map[string]string{
			"permissions": "0",
			"authorize":   "true",
		})

		req, err := http.NewRequest("POST", dUrl, bytes.NewBuffer(payload))
		req.Header.Set("authorization", token)
		req.Header.Set("content-type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}

		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		resp2 := login{}

		json.Unmarshal([]byte(body), &resp2)
		if resp2.Location == "" {
			profiles[id].DiscordSession = "error"
			fmt.Println(colors.Prefix() + colors.Red("Invalid Discord Token on profile ") + colors.White(profiles[id].Name) + colors.Red(" - Not running tasks on this profile"))
			return
		}

		code := resp2.Location
		code = r.FindStringSubmatch(code)[1]

		callBackUrl := "https://vlo.to/login/complete?code=" + code
		req, err = http.NewRequest("GET", callBackUrl, nil)
		if err != nil {
			profiles[id].DiscordSession = "error"
			fmt.Println(colors.Prefix() + colors.Red("Failed to login to profile ") + colors.White(profiles[id].Name) + colors.Red(" - Not running tasks on this profile"))
			return
		}

		req.Header.Set("cookie", "host=" + site)
		resp, err = client.Do(req)
		if err != nil {
			profiles[id].DiscordSession = "error"
			fmt.Println(colors.Prefix() + colors.Red("Failed to login to profile ") + colors.White(profiles[id].Name) + colors.Red(" - Not running tasks on this profile"))
			return
		}
		defer resp.Body.Close()
		var jwtToken string
		for name, values := range resp.Header {
			if name == "Location" {
				location := values[0]
				r2 := regexp.MustCompile("&token=(.*)")
				jwtToken = r2.FindStringSubmatch(location)[1]
			}
		}

		_, err = client.Get("https://vlo.to/dashboard/redirect?host=" + site + "&token=" + jwtToken)
		if err != nil {
			profiles[id].DiscordSession = "error"
			fmt.Println(colors.Prefix() + colors.Red("Failed to login to profile ") + colors.White(profiles[id].Name) + colors.Red(" - Not running tasks on this profile"))
			return
		}
		_, err = client.Get("https://" + site + "/token?token=" + token)
		if err != nil {
			profiles[id].DiscordSession = "error"
			fmt.Println(colors.Prefix() + colors.Red("Failed to login to profile ") + colors.White(profiles[id].Name) + colors.Red(" - Not running tasks on this profile"))
			return
		}
		profiles[id].DiscordSession = jwtToken
		fmt.Println(colors.Prefix() + colors.Green("Successfully logged in on profile ")+colors.White(profiles[id].Name))
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

func VeloStripe(profiles []loader.ProfileStruct, stripeKey string) []loader.ProfileStruct {
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

	for id := range profiles {
		wg.Add(1)
		go tokenLocal(&wg, id)
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

func VeloTask(wg *sync.WaitGroup, userData loader.UserDataStruct, id int, site siteStruct, password string, profile loader.ProfileStruct) {
	defer wg.Done()

	type getResponseStruct struct {
		Success bool `json:"success"`
		Status string `json:"status"`
		Checkout string `json:"checkout"`
	}

	type checkoutResponseStruct struct {
		Success bool `json:"success"`
		Status string `json:"status"`
	}

	client, err := CoolClient("localhost")
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

	go PwSharingSend(password, userData.Username, site.BackendName)

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
		go SendWebhook(userData.Webhook, WebhookContentStruct{
			Speed: "bruh idk bro",
			Module: "Velo",
			Site: site.DisplayName,
			Profile: profile.Name,
		})
		payload, _ := json.Marshal(map[string]string{
			"site": site.DisplayName,
			"module": "Velo",
			"speed": "idk bro",
			"mode": "Brr mode",
			"password": "Unknown",
			"user": userData.Username,
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