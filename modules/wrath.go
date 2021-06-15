package modules

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"

	"github.com/asaskevich/govalidator"
	"github.com/chrigeeel/satango/colors"
	"github.com/chrigeeel/satango/loader"
)

func WrathKeyInput(userData loader.UserDataStruct, profiles []loader.ProfileStruct) {
	
	profiles = askForProfiles(profiles)
	
	profiles = WrathKeyLogin(profiles)

	var taskLimit int
	taskLimit = len(profiles) * 2
	if taskLimit > 15 {
		taskLimit = 15
	}
	fmt.Println(colors.Prefix() + colors.Red("How many tasks do you want to run? You have ") + colors.White(strconv.Itoa(len(profiles))) + colors.Red(" profiles and your task limit is ") + colors.White(strconv.Itoa(taskLimit)) + colors.Red("."))
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
		key, _ := GetPw("wrath")
		if key == "exit" {
			exit = true
		} else {
			fmt.Println(colors.Prefix() + colors.Yellow("Starting tasks..."))
			var wg sync.WaitGroup
			for i := 0; i < taskAmount; i++ {
				if profileCounter+1 > len(profiles) {
					profileCounter = 0
				}
				go WrathKeyTask(&wg, userData, i+1, key, profiles[profileCounter])
				profileCounter++
			}
			wg.Wait()
		}
	}
}

func WrathKeyTask(wg *sync.WaitGroup, userData loader.UserDataStruct, id int, key string, profile loader.ProfileStruct) {
	type checkoutDataStruct struct {
		DiscordId string `json:"discordId"`
		Key string `json:"key"`
	}

	type claimResponseStruct struct {
		Success bool `json:"success"`
	}
	
	if len(key) != 29 {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Key must be 29 characters long!"))
		return
	}

	checkoutData := checkoutDataStruct{}
	checkoutData.DiscordId = profile.DiscordId
	checkoutData.Key = key

	payload, _ := json.Marshal(checkoutData)
	client := http.DefaultClient

	fmt.Println(colors.TaskPrefix(id) + colors.Yellow("Submitting claim..."))

	req, err := http.NewRequest("POST", "https://server.wrathbots.co/keybind", bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to initiate claim!"))
		return
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.101 Safari/537.36")
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to claim the key!"))
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	rdata := new(claimResponseStruct)
	json.Unmarshal([]byte(body), &rdata)
	if rdata.Success == false {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to claim key! Either already claimed or wrong key."))
		return
	}
	fmt.Println(colors.TaskPrefix(id) + colors.Green("Successfully claimed key on profile ") + colors.White("\"") + colors.Green(profile.Name) + colors.White("\""))
	go SendWebhook(userData.Webhook, WebhookContentStruct{
		Speed: "idk bro",
		Module: "Wrath",
		Site: "Wrath",
		Profile: profile.Name,
	})
	payload, _ = json.Marshal(map[string]string{
		"site": "Wrath",
		"module": "Wrath",
		"speed": "idk",
		"mode": "Normal",
		"password": "Unknown",
		"user": userData.Username,
	})
	req, _ = http.NewRequest("POST", "http://ec2-13-52-240-112.us-west-1.compute.amazonaws.com:3000/checkouts", bytes.NewBuffer(payload))
	req.Header.Set("content-type", "application/json")
	client.Do(req)
	return
}

func WrathKeyLogin(profiles []loader.ProfileStruct) []loader.ProfileStruct {
	type login struct {
		Id string `json:"id"`
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
		client := &http.Client{}
		url := "https://discord.com/api/v9/users/@me"

		req, err := http.NewRequest("GET", url, nil)
		req.Header.Set("authorization", token)
		req.Header.Set("content-type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			profiles[id].DiscordId = "error"
			fmt.Println(colors.Prefix() + colors.Red("Failed to login to profile ") + colors.White(profiles[id].Name) + colors.Red(" - Not running tasks on this profile"))
			return
		}
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		resp2 := login{}

		json.Unmarshal([]byte(body), &resp2)
		if resp2.Id == "" {
			profiles[id].DiscordId = "error"
			fmt.Println(colors.Prefix() + colors.Red("Invalid Discord Token on profile ") + colors.White(profiles[id].Name) + colors.Red(" - Not running tasks on this profile"))
			return
		}

		session := resp2.Id
		profiles[id].DiscordId = session
		fmt.Println(colors.Prefix() + colors.Green("Successfully logged in on profile ") + colors.White(profiles[id].Name))
		return
	}

	for id := range profiles {
		wg.Add(1)
		go loginLocal(&wg, id)
	}
	wg.Wait()

	for i := len(profiles) - 1; i >= 0; {
		if profiles[i].DiscordId == "error" {
			profiles = removeIndex(profiles, i)
		}
		i--
	}
	return profiles
}