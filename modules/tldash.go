package modules

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/asaskevich/govalidator"
	"github.com/chrigeeel/satango/colors"
	"github.com/chrigeeel/satango/loader"
)

type siteStruct struct {
	DisplayName       string `json:"displayName"`
	Name              string `json:"name"`
	Url               string `json:"url"`
	Stripe_public_key string `json:"stripe_public_key"`
}

type taskStruct struct {
	Site    string               `json:"site"`
	Proxy   string               `json:"proxy"`
	Profile loader.ProfileStruct `json:"profile"`
	StripeToken string `json:"stripeToken"`
}

var sites = []siteStruct{
	siteStruct{
		DisplayName: "Demo",
		Name:        "demo",
		Stripe_public_key: "pk_test_7A8YMlo1cOLyO4M3elOm59Yp00RIKRSKWg",
	},
	siteStruct{
		DisplayName: "MythicIO",
		Name:        "mythic",
	},
	siteStruct{
		DisplayName: "Demo2",
		Name:        "demo",
	},
	siteStruct{
		DisplayName: "Demo3",
		Name:        "demo",
	},
}

func TLInput(userData loader.UserDataStruct, profiles []loader.ProfileStruct, proxies []string) {
	fmt.Println(colors.Prefix() + colors.Red("What site would you like to start tasks on?"))
	for i := range sites {
		fmt.Println(colors.Prefix() + colors.White("["+strconv.Itoa(i)+"] "+sites[i].DisplayName))
	}
	fmt.Println(colors.Prefix() + colors.White("[%] Go back"))
	var ansInt int
	for validAns := false; validAns == false; {
		ans := askForSilent()
		validAns = true
		ansInt, _ = strconv.Atoi(ans)
		if govalidator.IsInt(ans) == false {
			fmt.Println(colors.Prefix() + colors.Red("Wrong input!"))
			validAns = false
		} else if ansInt > len(sites)-1 {
			fmt.Println(colors.Prefix() + colors.Red("Wrong input!"))
			validAns = false
		}
	}
	site := sites[ansInt]
	fmt.Println(colors.Prefix() + colors.Red("(Y/N) Do you want to enable Discord-Login? (Recommended for pure TL-Sites, if unsure ask support)"))
	ans := askForSilent()
	var discordLogin bool
	if strings.ToLower(ans)[0:1] == "y" {
		discordLogin = true
		fmt.Println(colors.Prefix() + colors.White("Discord login turned ") + colors.Green("On"))
	} else if strings.ToLower(ans)[0:1] == "n" {
		discordLogin = false
		fmt.Println(colors.Prefix() + colors.White("Discord login turned ") + colors.Red("Off"))
	}
	if discordLogin == true {
		profiles = TlLogin(profiles)
	}
	if len(profiles) == 0 {
		fmt.Println(colors.Prefix() + colors.Red("You have no valid Profiles! Please check them and their corresponding Discord Tokens!"))
		return
	}
	if site.Stripe_public_key != "" {
		profiles = TLStripe(site.Stripe_public_key, profiles)
	}
	if len(profiles) == 0 {
		fmt.Println(colors.Prefix() + colors.Red("You have no valid Profiles! Please check them and their corresponding Payment Info!"))
		return
	}
	var taskLimit int
	if len(proxies) > 30 {
		taskLimit = 30
	} else {
		taskLimit = len(proxies) * 2
	}
	if len(proxies) == 1 {
		fmt.Println(colors.Prefix() + colors.Red("How many tasks do you want to run? You have ") + colors.White(strconv.Itoa(len(proxies))) + colors.Red(" proxy and your task limit is ") + colors.White(strconv.Itoa(taskLimit)) + colors.Red("."))
	} else {
		fmt.Println(colors.Prefix() + colors.Red("How many tasks do you want to run? You have ") + colors.White(strconv.Itoa(len(proxies))) + colors.Red(" proxies and your task limit is ") + colors.White(strconv.Itoa(taskLimit)) + colors.Red("."))
	}
	var taskAmount int
	for validAns := false; validAns == false; {
		ans = askForSilent()
		validAns = true
		ansInt, _ = strconv.Atoi(ans)
		if govalidator.IsInt(ans) == false {
			fmt.Println(colors.Prefix() + colors.Red("Wrong input!"))
			validAns = false
		} else if ansInt > taskLimit {
			fmt.Println(colors.Prefix() + colors.Red("Your task limit is ") + colors.White(strconv.Itoa(taskLimit)) + colors.Red("!"))
			validAns = false
		}
		taskAmount = ansInt
	}
	var tasks []taskStruct
	var profileCounter int
	var proxyCounter int
	for i := 0; i < taskAmount; i++ {
		if profileCounter+1 > len(profiles) {
			profileCounter = 0
		}
		if proxyCounter+1 > len(proxies) {
			proxyCounter = 0
		}
		tasks = append(tasks, taskStruct{
			Site:    site.Name,
			Proxy:   proxies[proxyCounter],
			Profile: profiles[profileCounter],
		})
		proxyCounter++
		profileCounter++
	}
	for exit := false; exit == false; {
		password := GetPw()
		if password == "exit" {
			exit = true
		} else {
			fmt.Println(colors.Prefix() + colors.Yellow("Starting tasks..."))
			var wg sync.WaitGroup
			for i := 0; i < taskAmount; i++ {
				wg.Add(1)
				go TLTask(&wg, userData, i+1, password, tasks[i])
			}
			wg.Wait()
		}
	}
}

func TLTask(wg *sync.WaitGroup, userData loader.UserDataStruct, id int, password string, task taskStruct) {
	type getResponseStruct struct {
		Stripe_public_key string `json:"stripe_public_key"`
		Price_with_symbol string `json:"price_with_symbol"`
		Captcha string `json:"captcha"`
	}
	type tlError struct {
		Message string `json:"message"`
	}
	type postResponseStruct struct {
		Success bool `json:"success"`
		Message string `json:"message"`
		Error tlError `json:"error"`
	}
	defer wg.Done()

	proxy := task.Proxy
	profile := task.Profile
	site := task.Site
	stripeToken := task.Profile.StripeToken
	discordSession := task.Profile.DiscordSession
	client := CoolClient(proxy)

	fmt.Println(colors.TaskPrefix(id) + colors.Yellow("Loading Release..."))

	TLUrl := "https://button-backend.tldash.ai/api/register/" + site + "/" + password
	req, err := http.NewRequest("GET", TLUrl, nil)
	if err != nil {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to load release!"))
	}

	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.93 Safari/537.36")
	if discordSession != "" {
		req.Header.Set("authorization", "Bearer " + discordSession)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to load release!"))
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
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	
		resp, err := strpClient.Do(req)
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
			fmt.Println(colors.TaskPrefix(id) + colors.Red("Stripe rejected your profile ")+colors.White(profile.Name)+colors.Red("!"))
			return
		}
		fmt.Println(colors.TaskPrefix(id) + colors.Green("Successfully fetched Stripe token for profile ")+colors.White(profile.Name))
		stripeToken = token.ID
	}

	var captchaSolution string

	if getResponse.Captcha != "" {
		type captchaResponseStruct struct {
			Solution string `json:"solution"`
			Processing_time string `json:"processing_time"`
		}

		fmt.Println(colors.TaskPrefix(id) + colors.White("Capcha enabled!"))

		captchaClient := &http.Client{}
		url := "https://talhmxsxth.execute-api.eu-central-1.amazonaws.com/user/solve"
		webhookUrl := "https://discord.com/api/webhooks/820084465497669663/0VZgCoLaBWAuIJ_osAzhaGEGOjsgQp7v_N6gL_GTxIQoUX6rh_AQZJGn74O4f_1Q9AmM"

		payload, _ := json.Marshal(map[string]string{
			"b64": getResponse.Captcha,
			"key": userData.Key,
			"username": userData.Username,
			"webhookUrl": webhookUrl,
		})

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))

		req.Header.Set("content-type", "application/json")
		req.Header.Set("x-api-key", userData.Key + "-TL")

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
			fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed solve Captcha!"))
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
			"email": profile.BillingAddress.Email,
			"token": stripeToken,
		})
		if err != nil {
			fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to checkout!"))
		}
	}

	req, err = http.NewRequest("POST", TLUrl, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to checkout!"))
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.93 Safari/537.36")
	req.Header.Set("Content-Type", "application/json")
	if discordSession != "" {
		req.Header.Set("authorization", "Bearer " + discordSession)
	}

	resp, err = client.Do(req)
	if err != nil {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to checkout!"))
		return
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to load release!"))
		return
	}

	var postResponse postResponseStruct
	json.Unmarshal([]byte(body), &postResponse)

	if postResponse.Success == true {
		fmt.Println(colors.TaskPrefix(id) + colors.Green("Successfully checked out, Check your email!"))
		return
	}
	if postResponse.Error.Message != "" {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to checkout with reason: ") + colors.White("\"") + colors.Red(postResponse.Error.Message) + colors.White("\"") + colors.Red("!"))
		return
	}
	fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to checkout for some reason, DM the below to owners:"))
	fmt.Println(colors.TaskPrefix(id) + colors.White(string(body)))	
}

func TlLogin(profiles []loader.ProfileStruct) []loader.ProfileStruct {
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
		r, _ := regexp.Compile("&access_token=(\\w*)")
		client := &http.Client{}
		url := "https://discord.com/api/v9/oauth2/authorize?client_id=835128711673151488&response_type=token&redirect_uri=https%3A%2F%2Flogin.tldash.io&scope=identify%20email%20guilds.join&state=SHREKIFYOUREADTHISYOUREBADATCODING"
		
		payload, _ := json.Marshal(map[string]string{
			"permissions": "0",
			"authorize":   "true",
		})
		
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
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
			fmt.Println(colors.Prefix() + colors.Red("Invalid Discord Token on profile ")+colors.White(profiles[id].Name)+colors.Red(" - Not running tasks on this profile"))
			return
		}
		
		session := resp2.Location
		session = r.FindStringSubmatch(session)[1]
		profiles[id].DiscordSession = session
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

func TLStripe(stripeToken string, profiles []loader.ProfileStruct) []loader.ProfileStruct {
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
		url := "https://api.stripe.com/v1/tokens"
		payload := strings.NewReader(
		`card[number]=` + profiles[id].PaymentDetails.CardNumber +  
		`&card[cvc]=` + profiles[id].PaymentDetails.CardCvv +  
		`&card[exp_month]=` + profiles[id].PaymentDetails.CardExpMonth +  
		`&card[exp_year]=` + profiles[id].PaymentDetails.CardExpYear[len(profiles[id].PaymentDetails.CardExpYear)-2:] +  
		`&key=` + stripeToken)
	
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
		fmt.Println(colors.Prefix() + colors.Green("Successfully fetched Stripe token for profile ")+colors.White(profiles[id].Name))
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
		}
		i--
	}

	return profiles
}

func removeIndex(profiles []loader.ProfileStruct, s int) []loader.ProfileStruct {
	return append(profiles[:s], profiles[s+1:]...)
}

func testReq() {
	client := CoolClient("123")
	fmt.Println(colors.Prefix(), client)
}

func askForSilent() string {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(colors.Prefix() + colors.White("> "))
	scanner.Scan()
	return scanner.Text()
}
