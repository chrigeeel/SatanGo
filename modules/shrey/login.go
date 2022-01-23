package shrey

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"sync"

	"github.com/chrigeeel/satango/colors"
	"github.com/chrigeeel/satango/loader"
	"github.com/chrigeeel/satango/utility"
)

func login(profiles []loader.ProfileStruct, site string) []loader.ProfileStruct {
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

		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
		url := site + "login"
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			profiles[id].DiscordSession = "error"
			fmt.Println(colors.Prefix() + colors.Red("Invalid Discord Token on profile ") + colors.White(profiles[id].Name) + colors.Red(" - Not running tasks on this profile"))
			return
		}
		req.Header.Add("Turbolinks-Referrer", site)
		resp, err := client.Do(req)
		if err != nil {
			profiles[id].DiscordSession = "error"
			fmt.Println(colors.Prefix() + colors.Red("Invalid Discord Token on profile ") + colors.White(profiles[id].Name) + colors.Red(" - Not running tasks on this profile"))
			return
		}
		var shreyAuthSession string
		for _, cookie := range resp.Cookies() {
			if cookie.Name == "_shreyauth_session" {
				shreyAuthSession = cookie.Value
			}
		}
		
		req, err = http.NewRequest("GET", url, nil)
		if err != nil {
			profiles[id].DiscordSession = "error"
			fmt.Println(colors.Prefix() + colors.Red("Invalid Discord Token on profile ") + colors.White(profiles[id].Name) + colors.Red(" - Not running tasks on this profile"))
			return
		}
		req.Header.Add("sec-ch-ua", "\" Not;A Brand\";v=\"99\", \"Google Chrome\";v=\"91\", \"Chromium\";v=\"91\"")
		req.Header.Add("sec-ch-ua-mobile", "?0")
		req.Header.Add("Upgrade-Insecure-Requests", "1")
		req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
		req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
		req.Header.Add("Sec-Fetch-Site", "cross-site")
		req.Header.Add("Sec-Fetch-Mode", "navigate")
		req.Header.Add("Sec-Fetch-User", "?1")
		req.Header.Add("Sec-Fetch-Dest", "document")
		req.Header.Add("Cookie", "_shreyauth_session=" + shreyAuthSession)
		resp, err = client.Do(req)
		if err != nil {
			profiles[id].DiscordSession = "error"
			fmt.Println(colors.Prefix() + colors.Red("Invalid Discord Token on profile ") + colors.White(profiles[id].Name) + colors.Red(" - Not running tasks on this profile"))
			return
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			profiles[id].DiscordSession = "error"
			fmt.Println(colors.Prefix() + colors.Red("Invalid Discord Token on profile ") + colors.White(profiles[id].Name) + colors.Red(" - Not running tasks on this profile"))
			return
		}
		r := regexp.MustCompile(`href="([^\"]*)`)
		match := r.FindStringSubmatch(string(body))
		if len(match) != 2 {
			profiles[id].DiscordSession = "error"
			fmt.Println(colors.Prefix() + colors.Red("Invalid Discord Token on profile ") + colors.White(profiles[id].Name) + colors.Red(" - Not running tasks on this profile"))
			return
		}
		url = match[1]
		for _, cookie := range resp.Cookies() {
			if cookie.Name == "_shreyauth_session" {
				shreyAuthSession = cookie.Value
			}
		}
		req, err = http.NewRequest("GET", url, nil)
		if err != nil {
			profiles[id].DiscordSession = "error"
			fmt.Println(colors.Prefix() + colors.Red("Invalid Discord Token on profile ") + colors.White(profiles[id].Name) + colors.Red(" - Not running tasks on this profile"))
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
		_, err = client.Do(req)
		if err != nil {
			profiles[id].DiscordSession = "error"
			fmt.Println(colors.Prefix() + colors.Red("Invalid Discord Token on profile ") + colors.White(profiles[id].Name) + colors.Red(" - Not running tasks on this profile"))
			return
		}
		var shreyAuthSession2 string
		for _, cookie := range resp.Cookies() {
			if cookie.Name == "_shreyauth_session" {
				shreyAuthSession2 = cookie.Value
			}
		}

		r = regexp.MustCompile(`\?code=(\w*)`)
		url = "https://discord.com/api/v9/oauth2/authorize?client_id=601262335713345542&response_type=code&redirect_uri=https%3A%2F%2Fshreyauth.com%2Fdiscord%2Fconnect&scope=identify%20email%20guilds%20guilds.join"

		payload, _ := json.Marshal(map[string]string{
			"permissions": "0",
			"authorize":   "true",
		})

		req, err = http.NewRequest("POST", url, bytes.NewBuffer(payload))
		if err != nil {
			profiles[id].DiscordSession = "error"
			fmt.Println(colors.Prefix() + colors.Red("Failed to login to profile ") + colors.White(profiles[id].Name) + colors.Red(" - Not running tasks on this profile"))
			return
		}
		req.Header.Set("authorization", token)
		req.Header.Set("content-type", "application/json")

		resp, err = client.Do(req)
		if err != nil {
			profiles[id].DiscordSession = "error"
			fmt.Println(colors.Prefix() + colors.Red("Failed to login to profile ") + colors.White(profiles[id].Name) + colors.Red(" - Not running tasks on this profile"))
			return
		}
		defer resp.Body.Close()
		body, _ = ioutil.ReadAll(resp.Body)
		resp2 := login{}

		json.Unmarshal([]byte(body), &resp2)
		if resp2.Location == "" {
			profiles[id].DiscordSession = "error"
			fmt.Println(colors.Prefix() + colors.Red("Invalid Discord Token on profile ") + colors.White(profiles[id].Name) + colors.Red(" - Not running tasks on this profile"))
			return
		}

		code := resp2.Location
		code = r.FindStringSubmatch(code)[1]
		url = "https://shreyauth.com/discord/connect?code=" + code
		req, err = http.NewRequest("GET", url, nil)
		if err != nil {
			profiles[id].DiscordSession = "error"
			fmt.Println(colors.Prefix() + colors.Red("Invalid Discord Token on profile ") + colors.White(profiles[id].Name) + colors.Red(" - Not running tasks on this profile"))
			return
		}
		req.Header.Add("sec-ch-ua", "\" Not;A Brand\";v=\"99\", \"Google Chrome\";v=\"91\", \"Chromium\";v=\"91\"")
		req.Header.Add("sec-ch-ua-mobile", "?0")
		req.Header.Add("Upgrade-Insecure-Requests", "1")
		req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
		req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
		req.Header.Add("Sec-Fetch-Site", "cross-site")
		req.Header.Add("Sec-Fetch-Mode", "navigate")
		req.Header.Add("Sec-Fetch-User", "?1")
		req.Header.Add("Sec-Fetch-Dest", "document")
		req.Header.Add("Cookie", "_shreyauth_session=" + shreyAuthSession2)
		_, err = client.Do(req)
		if err != nil {
			profiles[id].DiscordSession = "error"
			fmt.Println(colors.Prefix() + colors.Red("Invalid Discord Token on profile ") + colors.White(profiles[id].Name) + colors.Red(" - Not running tasks on this profile"))
			return
		}

		url = site + "discord/complete?code=" + code

		req, err = http.NewRequest("GET", url, nil)
		if err != nil {
			profiles[id].DiscordSession = "error"
			fmt.Println(colors.Prefix() + colors.Red("Invalid Discord Token on profile ") + colors.White(profiles[id].Name) + colors.Red(" - Not running tasks on this profile"))
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
		req.Header.Add("Cookie", "_shreyauth_session=" + shreyAuthSession)
		resp, err = client.Do(req)
		if err != nil {
			profiles[id].DiscordSession = "error"
			fmt.Println(colors.Prefix() + colors.Red("Invalid Discord Token on profile ") + colors.White(profiles[id].Name) + colors.Red(" - Not running tasks on this profile"))
			return
		}
		for _, cookie := range resp.Cookies() {
			if cookie.Name == "_shreyauth_session" {
				shreyAuthSession = cookie.Value
			}
		}
		profiles[id].DiscordSession = shreyAuthSession
		fmt.Println(colors.Prefix() + colors.Green("Successfully logged in on profile ") + colors.White(profiles[id].Name))
	}

	for id := range profiles {
		wg.Add(1)
		go loginLocal(&wg, id)
	}
	wg.Wait()

	for i := len(profiles) - 1; i >= 0; {
		if profiles[i].DiscordSession == "error" {
			profiles = utility.RemoveIndex(profiles, i)
		}
		i--
	}
	return profiles
}