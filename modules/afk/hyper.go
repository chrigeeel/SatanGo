package afk

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
	"github.com/chrigeeel/satango/modules/getpw"
	"github.com/chrigeeel/satango/modules/hyper"
	"github.com/chrigeeel/satango/utility"
)

func taskfcfs(wg *sync.WaitGroup, userData loader.UserDataStruct, id int, p getpw.PWStruct, profile loader.ProfileStruct) {
	defer wg.Done()

	beginTime := time.Now()

	checkoutData := hyper.HyperCheckoutStruct{}

	checkoutData.Billing_details.Email = profile.BillingAddress.Email
	checkoutData.Billing_details.Name = profile.BillingAddress.Name

	checkoutData.Release = p.HyperInfo.ReleaseId

	if p.HyperInfo.RequireLogin {
		checkoutData.User = profile.HyperUserId
	}

	if p.HyperInfo.CollectBilling {
		checkoutData.Billing_details.Address.City = profile.BillingAddress.City
		checkoutData.Billing_details.Address.Country = "US"
		checkoutData.Billing_details.Address.Line1 = profile.BillingAddress.Line1
		checkoutData.Billing_details.Address.Postal_code = profile.BillingAddress.PostCode
		checkoutData.Billing_details.Address.State = profile.BillingAddress.City
	}

	payload, _ := json.Marshal(checkoutData)

	fmt.Println(colors.TaskPrefix(id) + colors.Yellow("Submitting Checkout..."))

	req, err := http.NewRequest("POST", p.Site + "ajax/checkouts", bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to initiate Checkout!"))
		return
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("hyper-account", p.HyperInfo.AccountId)

	req.Header.Set("x-amz-cf-id", p.HyperInfo.BpToken)

	client := http.DefaultClient

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to execute Checkout!"))
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

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
		req, err := http.NewRequest("GET", p.Site + "ajax/checkouts/"+rdata.ID, nil)
		if err != nil {
			fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to check Success!"))
			return
		}
		req.Header.Set("Cookie", "auhorization="+profile.DiscordSession)
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36")
		req.Header.Set("hyper-account", p.HyperInfo.AccountId)
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
					Site:    p.Site,
					Profile: profile.Name,
				})
				payload, _ := json.Marshal(map[string]string{
					"site":     p.Site,
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

func taskfcfsunk(wg *sync.WaitGroup, userData loader.UserDataStruct, id int, password string, site string, profile loader.ProfileStruct) {
	defer wg.Done()

	beginTime := time.Now()

	fmt.Println(colors.TaskPrefix(id) + colors.Yellow("Loading Release..."))
	req, err := http.NewRequest("GET", site+"purchase?password="+password, nil)
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
	
	var mdata string
	if len(r.FindStringSubmatch(string(body))) > 0 {
		mdata = r.FindStringSubmatch(string(body))[1]
	} else {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to load release!"))
		return
	}

	page := new(hyper.HyperStruct)

	json.Unmarshal([]byte(mdata), &page)

	hyperAccountId := page.Props.Pageprops.Account.ID
	requireLogin := page.Props.Pageprops.Account.Settings.Payments.RequireLogin
	botProtection := page.Props.Pageprops.Account.Settings.BotProtection.Enabled
	collectBilling := page.Props.Pageprops.Account.Settings.Payments.CollectBillingAddress
	releaseId := page.Query.Release
	oos := page.Props.Pageprops.Release.OutOfStock
	if releaseId == "" || oos {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Wrong password or release OOS!"))
		return
	}
	fmt.Println(colors.TaskPrefix(id) + colors.Green("Release instock!"))

	checkoutData := hyper.HyperCheckoutStruct{}

	checkoutData.Billing_details.Email = profile.BillingAddress.Email
	checkoutData.Billing_details.Name = profile.BillingAddress.Name

	checkoutData.Release = releaseId

	if requireLogin {
		checkoutData.User = profile.HyperUserId
	}

	if collectBilling {
		checkoutData.Billing_details.Address.City = profile.BillingAddress.City
		checkoutData.Billing_details.Address.Country = "US"
		checkoutData.Billing_details.Address.Line1 = profile.BillingAddress.Line1
		checkoutData.Billing_details.Address.Postal_code = profile.BillingAddress.PostCode
		checkoutData.Billing_details.Address.State = profile.BillingAddress.City
	}

	payload, _ := json.Marshal(checkoutData)

	fmt.Println(colors.TaskPrefix(id) + colors.Yellow("Submitting Checkout..."))

	req, err = http.NewRequest("POST", site+"ajax/checkouts", bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to initiate Checkout!"))
		return
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("hyper-account", hyperAccountId)

	if botProtection {
		fmt.Println(colors.TaskPrefix(id) + colors.Yellow("Bot Protection enabled..."))
		bpToken, err := hyper.Solvebp(site)
		if err != nil {
			fmt.Println(colors.Prefix() + colors.Red("Failed to solve Bot Protection. Please contact Chrigeeel or Shrek!"))
		} else {
			req.Header.Set("x-amz-cf-id", bpToken)
			fmt.Println(colors.TaskPrefix(id) + colors.Green("Successfully solved Bot Protection!"))
		}
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