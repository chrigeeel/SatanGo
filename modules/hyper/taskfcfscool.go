package hyper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/chrigeeel/satango/colors"
	"github.com/chrigeeel/satango/loader"
	"github.com/chrigeeel/satango/modules/utility"
)

func taskfcfsCool(wg *sync.WaitGroup, userData loader.UserDataStruct, id int, releaseId string, paid bool, collectBilling bool, site string, hyperAccountId string, profile loader.ProfileStruct, bpToken string) {
	defer wg.Done()

	beginTime := time.Now()

	checkoutData := hyperCheckoutStruct{}

	checkoutData.Billing_details.Email = profile.BillingAddress.Email
	checkoutData.Billing_details.Name = profile.BillingAddress.Name

	checkoutData.Release = releaseId

	checkoutData.User = profile.HyperUserId

	if paid {
		checkoutData.Payment_method = profile.StripeToken
	}

	if collectBilling {
		checkoutData.Billing_details.Address.City = profile.BillingAddress.City
		checkoutData.Billing_details.Address.Country = "US"
		checkoutData.Billing_details.Address.Line1 = profile.BillingAddress.Line1
		checkoutData.Billing_details.Address.Postal_code = profile.BillingAddress.PostCode
		checkoutData.Billing_details.Address.State = profile.BillingAddress.City
		if paid {
			checkoutData.Payment_method = profile.StripeToken2
		}
	}

	payload, _ := json.Marshal(checkoutData)

	fmt.Println(string(payload))

	fmt.Println(colors.TaskPrefix(id) + colors.Yellow("Submitting Checkout..."))

	fmt.Println(site + "ajax/checkouts")

	req, err := http.NewRequest("POST", site+"ajax/checkouts", bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to initiate Checkout!"))
		return
	}

	req.Header.Set("Cookie", "auhorization="+profile.DiscordSession)

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("hyper-account", hyperAccountId)

	req.Header.Set("x-amz-cf-id", bpToken)

	client := http.DefaultClient

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to execute Checkout!"))
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))

	type hyperResponseStruct struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	}

	rdata := new(hyperResponseStruct)
	json.Unmarshal([]byte(body), &rdata)
	if rdata.ID == "" {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Out of Stock!"))
		return
	}

	for {
		fmt.Println(colors.TaskPrefix(id) + colors.Yellow("Processing..."))
		req, err := http.NewRequest("GET", site+"ajax/checkouts/"+rdata.ID, nil)
		if err != nil {
			fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to check Success!"))
			return
		}
		req.Header.Set("Cookie", "auhorization="+profile.DiscordSession)
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36")
		req.Header.Set("hyper-account", hyperAccountId)
		resp, err = client.Do(req)
		if err != nil {
			fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to check Success!"))
			return
		}
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &rdata)
		if rdata.Status != "processing" {
			if rdata.Status == "succeeded" {
				stopTime := time.Now()
				diff := stopTime.Sub(beginTime)
				fmt.Println(colors.TaskPrefix(id) + colors.Green("Successfully checked out on profile \""+colors.White(profile.Name)+colors.Green("\"")))
				go utility.SendWebhook(userData.Webhook, utility.WebhookContentStruct{
					Speed:   diff.String(),
					Module:  "Hyper / Meta Labs",
					Site:    site,
					Profile: profile.Name,
				})
				payload, _ := json.Marshal(map[string]string{
					"site":     site,
					"module":   "Hyper / Meta Labs",
					"speed":    diff.String(),
					"mode":     "Normal",
					"password": "Unknown",
					"user":     userData.Username,
				})
				req, _ := http.NewRequest("POST", "http://ec2-13-52-240-112.us-west-1.compute.amazonaws.com:3000/checkouts", bytes.NewBuffer(payload))
				req.Header.Set("content-type", "application/json")
				client.Do(req)
				return
			}
			fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to check out on profile \""+colors.White(profile.Name)+colors.Red("\" Reason: "+rdata.Status+" (This means either OOS or already registered!)")))
			return
		}
		time.Sleep(time.Millisecond * 1000)
	}
}