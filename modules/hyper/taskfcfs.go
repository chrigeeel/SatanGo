package hyper

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
	"github.com/chrigeeel/satango/modules/utility"
)

func taskfcfs(wg *sync.WaitGroup, userData loader.UserDataStruct, id int, password string, paid bool, site string, hyperAccountId string, profile loader.ProfileStruct, bpToken string) {
	defer wg.Done()

	beginTime := time.Now()

	fmt.Println(colors.TaskPrefix(id) + colors.Yellow("Loading Release..."))

	req, err := http.NewRequest("GET", site+"purchase/?password="+password, nil)
	if err != nil {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to load release!"))
		return
	}

	client, err := utility.CoolClient("localhost")
	if err != nil {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to load release!"))
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to load release!"))
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to load release!"))
		return
	}
	r, _ := regexp.Compile("__NEXT_DATA__\" type=\"application\\/json\">({.*})")
	mdata := r.FindStringSubmatch(string(body))[1]

	page := new(hyperStruct)

	json.Unmarshal([]byte(mdata), &page)

	requireLogin := page.Props.Pageprops.Account.Settings.Payments.RequireLogin
	botProtection := page.Props.Pageprops.Account.Settings.BotProtection.Enabled
	collectBilling := page.Props.Pageprops.Account.Settings.Payments.CollectBillingAddress
	releaseId := page.Query.Release

	if releaseId == "" {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Wrong password or release OOS!"))
		return
	}

	checkoutData := hyperCheckoutStruct{}

	checkoutData.Billing_details.Email = profile.BillingAddress.Email
	checkoutData.Billing_details.Name = profile.BillingAddress.Name

	checkoutData.Release = releaseId

	if requireLogin {
		checkoutData.User = profile.HyperUserId
	}

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

	fmt.Println(colors.TaskPrefix(id) + colors.Yellow("Submitting Checkout..."))

	req, err = http.NewRequest("POST", site+"ajax/checkouts", bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to initiate Checkout!"))
		return
	}

	if requireLogin {
		req.Header.Set("Cookie", "auhorization="+profile.DiscordSession)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("hyper-account", hyperAccountId)

	if botProtection {
		fmt.Println(colors.TaskPrefix(id) + colors.Yellow("Bot Protection enabled..."))
		req.Header.Set("x-amz-cf-id", bpToken)
	}

	resp, err = client.Do(req)
	if err != nil {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to execute Checkout!"))
		return
	}
	defer resp.Body.Close()
	body, _ = ioutil.ReadAll(resp.Body)

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
		time.Sleep(time.Second * 1)
	}
}