package velo

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

		req.Header.Set("cookie", "host="+site)
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