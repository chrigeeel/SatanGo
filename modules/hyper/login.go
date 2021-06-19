package hyper

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

func login(hyperAccountId string, site string, profiles []loader.ProfileStruct) []loader.ProfileStruct {
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
		r, _ := regexp.Compile(`\?code=(\w*)`)
		client := &http.Client{}
		dUrl := "https://discord.com/api/v9/oauth2/authorize?client_id=648234176805470248&response_type=code&redirect_uri=https%3A%2F%2Fapi.hyper.co%2Fportal%2Fauth%2Fdiscord%2Fcallback&scope=identify%20email%20guilds%20guilds.join&state=%7B%22account%22%3A%22" + hyperAccountId + "%22%7D"

		payload, _ := json.Marshal(map[string]string{
			"permissions": "0",
			"authorize":   "true",
		})

		req, err := http.NewRequest("POST", dUrl, bytes.NewBuffer(payload))
		if err != nil {
			profiles[id].DiscordSession = "error"
			profiles[id].HyperUserId = "error"
			fmt.Println(colors.Prefix() + colors.Red("Failed to login to profile ") + colors.White(profiles[id].Name) + colors.Red(" - Not running tasks on this profile"))
			return
		}
		req.Header.Set("authorization", token)
		req.Header.Set("content-type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			profiles[id].DiscordSession = "error"
			profiles[id].HyperUserId = "error"
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

		callBackUrl := "https://api.hyper.co/portal/auth/discord/callback?code=" + code + "&state=%7B%22account%22%3A%22" + hyperAccountId + "%22%7D"
		resp, err = client.Get(callBackUrl)
		if err != nil {
			profiles[id].DiscordSession = "error"
			profiles[id].HyperUserId = "error"
			fmt.Println(colors.Prefix() + colors.Red("Failed to login to profile ") + colors.White(profiles[id].Name) + colors.Red(" - Not running tasks on this profile"))
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

		req, _ = http.NewRequest("GET", site+"ajax/user", nil)

		req.Header.Set("cookie", "authorization="+hyperToken)
		req.Header.Set("hyper-account", hyperAccountId)

		resp, err = client.Do(req)
		if err != nil {
			profiles[id].DiscordSession = "error"
			profiles[id].HyperUserId = "error"
			fmt.Println(colors.Prefix() + colors.Red("Failed to login to profile ") + colors.White(profiles[id].Name) + colors.Red(" - Not running tasks on this profile"))
			return
		}
		defer resp.Body.Close()

		user := hyperUserStruct{}

		body, _ = ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &user)

		if user.ID == "" {
			profiles[id].DiscordSession = "error"
			profiles[id].HyperUserId = "error"
			fmt.Println(colors.Prefix() + colors.Red("Failed to login to profile ") + colors.White(profiles[id].Name) + colors.Red(" - Not running tasks on this profile"))
			return
		}

		profiles[id].DiscordSession = hyperToken
		profiles[id].HyperUserId = user.ID
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