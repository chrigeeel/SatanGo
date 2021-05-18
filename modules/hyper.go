package modules

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

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


func HyperInput(userData loader.UserDataStruct, profiles []loader.ProfileStruct, proxies []string) {
	fmt.Println(colors.Prefix() + colors.Red("What site would you like to start tasks on?") + colors.White(" (example: \"dashboard.satanbots.com\")"))
	site := askForSilent()
	r, _ := regexp.Compile("https:\\/\\/\\w")
	formatted := r.MatchString(site)
	if formatted == false {
		site = "https://" + site
	}
	if site[len(site)-1:] != "/" {
		site = site + "/"
	}
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

	loadPage, err := HyperLoad(site)
	if err != nil {
		fmt.Println(colors.Red("Failed to load site " + site))
		return
	}
	fmt.Println(loadPage)
	if paid == true {
		fmt.Println(colors.Yellow("Fetching stripe tokens for all Profiles..."))
		profiles = HyperStripe(loadPage.Props.Pageprops.Account.Stripe_account, profiles)
	}
	if len(profiles) == 0 {
		fmt.Println(colors.Prefix() + colors.Red("You have no valid Profiles! Please check them and their corresponding Payment Info!"))
		time.Sleep(time.Second * 3)
		return
	}
	fmt.Println(colors.Yellow("Logging in on all profiles..."))
	profiles = HyperLogin(loadPage.Props.Pageprops.Account.ID, site, profiles)
	if len(profiles) == 0 {
		fmt.Println(colors.Prefix() + colors.Red("You have no valid Profiles! Please check them and their corresponding Discord Tokens!"))
		time.Sleep(time.Second * 3)
		return
	}
}

func HyperLoad(site string) (hyperStruct, error) {
	loadPage := hyperStruct{}

	client := http.Client{Timeout: 5 * time.Second}
	url := site + "purchase"
	fmt.Println(url)
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println(colors.Prefix() + colors.Red("Failed to load site " + site))
		return loadPage, errors.New("Failed to load site")
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	r, _ := regexp.Compile("__NEXT_DATA__\" type=\"application\\/json\">({.*})")
	mdata := r.FindStringSubmatch(string(body))[1]
	
	fmt.Println(mdata)

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
		payload := strings.NewReader(
		`type=card` + 
		`&billing_details[name]=` + profiles[id].BillingAddress.Name + 
		`&card[number]=` + profiles[id].PaymentDetails.CardNumber +  
		`&card[cvc]=` + profiles[id].PaymentDetails.CardCvv +  
		`&card[exp_month]=` + profiles[id].PaymentDetails.CardExpMonth +  
		`&card[exp_year]=` + profiles[id].PaymentDetails.CardExpYear +  
		`&key=pk_live_51GXa1YLZrAkO7Fk2tcUO7vabkO7sgDamOww2OPYQVFhPZOllT75f7owzIOlP75MMdDXHKoy6wPt40HsuQDObpkHv004T74fAzs` +
		`&_stripe_account=` + stripeAccount)
	
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

func HyperTask(wg *sync.WaitGroup, userData loader.UserDataStruct, id int, password string, profile loader.ProfileStruct) {
	defer wg.Done()
}