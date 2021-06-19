package tldash

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
	"github.com/chrigeeel/satango/modules/utility"
)

func login(profiles []loader.ProfileStruct) []loader.ProfileStruct {
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
		r, _ := regexp.Compile(`&access_token=(\w*)`)
		client := &http.Client{}
		url := "https://discord.com/api/v9/oauth2/authorize?client_id=835128711673151488&response_type=token&redirect_uri=https%3A%2F%2Flogin.tldash.io&scope=identify%20email%20guilds.join&state=SHREKIFYOUREADTHISYOUREBADATCODING"

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

		session := resp2.Location
		session = r.FindStringSubmatch(session)[1]
		profiles[id].DiscordSession = session
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