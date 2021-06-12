package modules

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/chrigeeel/satango/colors"
	"github.com/chrigeeel/satango/loader"
)

func ShinobiInput(userData loader.UserDataStruct, profiles []loader.ProfileStruct, proxies []string) {
	fmt.Println(colors.Prefix() + colors.Yellow("Fetching Stripe Tokens for all profiles..."))
	profiles = ShinobiStripe(profiles)
	if len(profiles) == 0 {
		fmt.Println(colors.Prefix() + colors.Red("You have no valid Profiles! Please check them and their corresponding Payment Info!"))
		time.Sleep(time.Second * 3)
		return
	}

	profiles = askForProfiles(profiles)

	var taskLimit int
	taskLimit = len(profiles) * 5

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
		password, _ := GetPw("shinobi")
		if password == "exit" {
			exit = true
		} else {
			fmt.Println(colors.Prefix() + colors.Yellow("Starting tasks..."))
			var wg sync.WaitGroup
			for i :=0; i < taskAmount; i++ {
				if profileCounter+1 > len(profiles) {
					profileCounter = 0
				}
				wg.Add(1)
				go ShinobiTask(&wg, userData, i+1, password, profiles[profileCounter])
				profileCounter++
			}
			wg.Wait()
		}
	}

}

func ShinobiTask(wg *sync.WaitGroup, userData loader.UserDataStruct, id int, password string, profile loader.ProfileStruct) {
	type postResponseStruct struct {
		Success bool `json:"success"`
		Message string `json:"message"`
	}
	defer wg.Done()
	beginTime := time.Now()
	client, err := CoolClient("localhost")
	if err != nil {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Invalid Proxy!"))
		return
	}
	fmt.Println(colors.TaskPrefix(id) + colors.Yellow("Submitting checkout..."))
	payload, err := json.Marshal(map[string]string{
		"token": profile.StripeToken,
		"email": profile.BillingAddress.Email,
		"password": password,
	})
	if err != nil {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to checkout!"))
	}
	shinobiUrl := "https://dashboard.shinobi-scripts.com/api/payment/purchase"
	req, err := http.NewRequest("POST", shinobiUrl, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to checkout!"))
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.93 Safari/537.36")
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to checkout!"))
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to load release!"))
		return
	}

	var postResponse postResponseStruct
	json.Unmarshal([]byte(body), &postResponse)

	if postResponse.Success == true {
		stopTime := time.Now()
		diff := stopTime.Sub(beginTime)
		fmt.Println(colors.TaskPrefix(id) + colors.Green("Successfully checked out, Check your email!"))
		go SendWebhook(userData.Webhook, WebhookContentStruct{
			Speed: string(diff),
			Module: "Shinobi",
			Site: "Shinobi",
			Profile: profile.Name,
		})
		payload, _ := json.Marshal(map[string]string{
			"site": "Shinobi",
			"module": "Shinobi",
			"speed": diff.String(),
			"mode": "Brr mode",
			"password": password,
			"user": userData.Username,
		})
		req, _ := http.NewRequest("POST", "http://ec2-13-52-240-112.us-west-1.compute.amazonaws.com:3000/checkouts", bytes.NewBuffer(payload))
		req.Header.Set("content-type", "application/json")
		client.Do(req)
		return
	}
	if postResponse.Message == "sold_out" {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Release OOS!"))
		return
	}
	if postResponse.Message == "Cannot read property 'stock' of null" {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Password incorrect!"))
		return
	}
	fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to checkout with reason: " + string(body)))

}

func ShinobiStripe(profiles []loader.ProfileStruct) []loader.ProfileStruct {
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
		if len(profiles[id].PaymentDetails.CardExpYear) < 2 {
			profiles[id].StripeToken = "error"
			fmt.Println(colors.Prefix() + colors.Red("Stripe rejected your profile ") + colors.White(profiles[id].Name) + colors.Red(" - Not running tasks on this profile"))
			return
		}
		payload := strings.NewReader(
			`card[number]=` + profiles[id].PaymentDetails.CardNumber +
				`&card[cvc]=` + profiles[id].PaymentDetails.CardCvv +
				`&card[exp_month]=` + profiles[id].PaymentDetails.CardExpMonth +
				`&card[exp_year]=` + profiles[id].PaymentDetails.CardExpYear[len(profiles[id].PaymentDetails.CardExpYear)-2:] +
				`&key=` + `pk_live_CoRu9eKCeKLA85AZFqQ8B5lk002ljSbgQl`)

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
			fmt.Println(colors.Prefix() + colors.Red("Stripe rejected your profile ") + colors.White(profiles[id].Name) + colors.Red(" - Not running tasks on this profile"))
			return
		}
		profiles[id].StripeToken = token.ID
		fmt.Println(colors.Prefix() + colors.Green("Successfully fetched Stripe token for profile ") + colors.White(profiles[id].Name))
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