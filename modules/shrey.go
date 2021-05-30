package modules

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
)


func ShreyInput(userData loader.UserDataStruct, profiles []loader.ProfileStruct, proxies []string) {
	fmt.Println("bruh")
	return
}

func ShreyLogin(profiles []loader.ProfileStruct) []loader.ProfileStruct{
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
		dUrl := "https://discord.com/api/v9/oauth2/authorize?client_id=601262335713345542&response_type=code&redirect_uri=https%3A%2F%2Fshreyauth.com%2Fdiscord%2Fconnect&scope=identify%20email%20guilds%20guilds.join"
		
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

		callBackUrl := "https://shreyauth.com/discord/connect?code=" + code
		resp, err = client.Get(callBackUrl)
		if err != nil {
			profiles[id].ShreyCookie = "error"
			fmt.Println(colors.Prefix() + colors.Red("Failed to login to profile ")+colors.White(profiles[id].Name)+colors.Red(" - Not running tasks on this profile"))
			return
		}
		defer resp.Body.Close()

		callBackUrl2 := "https://dash.guap.io/discord/complete?code=" + code + code + code + code + code
		req, err = http.NewRequest("GET", callBackUrl2, nil)
		if err != nil {
			profiles[id].ShreyCookie = "error"
			fmt.Println(colors.Prefix() + colors.Red("Failed to login to profile ")+colors.White(profiles[id].Name)+colors.Red(" - Not running tasks on this profile"))
			return
		}

		resp, err = client.Do(req)
		if err != nil {
			profiles[id].ShreyCookie = "error"
			fmt.Println(colors.Prefix() + colors.Red("Failed to login to profile ")+colors.White(profiles[id].Name)+colors.Red(" - Not running tasks on this profile"))
			return
		}
		defer resp.Body.Close()
		var session2 string
		for _, cookie := range resp.Cookies() {
			if cookie.Name == "_shreyauth_session" {
				session2 = cookie.Value
			}
		}
		if session2 == "" {
			profiles[id].ShreyCookie = "error"
		}
		profiles[id].ShreyCookie = session2
		fmt.Println(colors.Prefix() + colors.Green("Successfully logged in on profile ")+colors.White(profiles[id].Name))


		req, err = http.NewRequest("GET", "https://dash.guap.io/subscribe", nil)
		if err != nil {
			fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to load release!"))
			return
		}
		session2 = "9M8swceuEXNqoVIFVQtWHav9eF8zdaPsOYjAnrurQ38%2BKq0TG5aG951YPF3P22l95NW1HlBDtlcVsBKmURmAOAeV4tFfxMTtRmouRJf6WX%2BpXZK52YOXYI0fXkuOdAbeBRQIoLNUXjYtE0hZuNgKIhB8CRz%2BgSERrJX41vNFP%2FhBpxqBKCeTB0GG70%2FZSHVcHh9T--t%2FaFQoXlZXKwjMsK--6VVWXQmqstQXUTJBAFxlrA%3D%3D"
		req.Header.Set("Cookie", "_shreyauth_session" + session2)
		resp, err = client.Do(req)
		if err != nil {
			fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to load release!"))
			return
		}
		defer resp.Body.Close()
		body, _ = ioutil.ReadAll(resp.Body)
		fmt.Println(string(body))
		r, _ = regexp.Compile("<meta name=\"csrf-token\" content=\"([^\"]*)")
		mdata := r.FindStringSubmatch(string(body))[1]
		fmt.Println(mdata)
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