package torpedo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"sync"
	"time"

	"github.com/chrigeeel/satango/colors"
	"github.com/chrigeeel/satango/loader"
	"github.com/chrigeeel/satango/utility"
)

func login(profiles []loader.ProfileStruct) []loader.ProfileStruct {

	type login struct {
		Location string `json:"location"`
	}

	var wg sync.WaitGroup

	fmt.Println(colors.Prefix() + colors.Yellow("Logging in on all profiles..."))

	loginLocal := func(wg *sync.WaitGroup, id int) {
		defer wg.Done()
	
		token := profiles[id].DiscordToken
		if token == "" {
			profiles[id].DiscordSession = "error"
			fmt.Println(colors.Prefix() + colors.Red("No Discord Token on profile ") + colors.White(profiles[id].Name) + colors.Red(" - Not running tasks on this profile"))
			return
		}
		r, _ := regexp.Compile(`\?code=(\w*)`)
		client := &http.Client{
			Timeout: 7 * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
		url := "https://discord.com/api/v9/oauth2/authorize?client_id=771970580257570816&response_type=code&redirect_uri=https%3A%2F%2Fdashboard.torpedoindustries.com%2Flogin%2Fcallback&scope=identify%20email%20guilds.join%20guilds"

		payload, _ := json.Marshal(map[string]string{
			"permissions": "0",
			"authorize":   "true",
		})

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
		if err != nil {
			profiles[id].DiscordSession = "error"
			fmt.Println(colors.Prefix() + colors.Red("Failed to login to profile ") + colors.White(profiles[id].Name) + colors.Red(" - Not running tasks on this profile"))
			return
		}
		req.Header.Set("authorization", token)
		req.Header.Set("content-type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			profiles[id].DiscordSession = "error"
			fmt.Println(colors.Prefix() + colors.Red("Failed to login to profile ") + colors.White(profiles[id].Name) + colors.Red(" - Not running tasks on this profile"))
			return
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

		callBackUrl := "https://dashboard.torpedoindustries.com/login/callback?code=" + code
		resp, err = client.Get(callBackUrl)
		if err != nil {
			profiles[id].DiscordSession = "error"
			fmt.Println(colors.Prefix() + colors.Red("Failed to login to profile ") + colors.White(profiles[id].Name) + colors.Red(" - Not running tasks on this profile"))
			return
		}
		defer resp.Body.Close()

		var sessionCookie string

		for _, cookie := range resp.Cookies() {
			if cookie.Name == "dashboard_session" {
				sessionCookie = cookie.Value
			}
		}
		if sessionCookie == "" {
			profiles[id].DiscordSession = "error"
			fmt.Println(colors.Prefix() + colors.Red("Failed to login to profile ") + colors.White(profiles[id].Name) + colors.Red(" - Not running tasks on this profile"))
			return
		}
		profiles[id].DiscordSession = sessionCookie
		fmt.Println(colors.Prefix() + colors.Green("Successfully logged in on profile ") + colors.White(profiles[id].Name))
	}

	for id := range profiles {
		wg.Add(1)
		go loginLocal(&wg, id)
	}
	wg.Wait()

	for i := len(profiles) - 1; i >= 0; {
		if profiles[i].DiscordId == "error" {
			profiles = utility.RemoveIndex(profiles, i)
		}
		i--
	}
	return profiles
}