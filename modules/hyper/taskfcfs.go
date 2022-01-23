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
	"github.com/chrigeeel/satango/modules/getpw"
	"github.com/chrigeeel/satango/utility"
)

func taskfcfs(wg *sync.WaitGroup, userData loader.UserDataStruct, id int, password string, paid bool, site string, hyperAccountId string, profile loader.ProfileStruct, bpToken string) {
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
	r := regexp.MustCompile("__NEXT_DATA__\" type=\"application\\/json\">({.*})")
	
	var mdata string
	if len(r.FindStringSubmatch(string(body))) > 0 {
		mdata = r.FindStringSubmatch(string(body))[1]
	} else {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to load release!"))
		return
	}

	page := new(HyperStruct)

	json.Unmarshal([]byte(mdata), &page)

	requireLogin := page.Props.Pageprops.Account.Settings.Payments.RequireLogin
	botProtection := page.Props.Pageprops.Account.Settings.BotProtection.Enabled
	collectBilling := page.Props.Pageprops.Account.Settings.Payments.CollectBillingAddress
	releaseId := page.Query.Release
	if releaseId == "" { releaseId = page.Query.Link }
	oos := page.Props.Pageprops.Release.OutOfStock

	if oos || releaseId == "" {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Wrong password or release OOS!"))
		return
	}

	hyperData := getpw.HyperInfo{
		ReleaseId: releaseId,
		BpToken: bpToken,
		CollectBilling: collectBilling,
		RequireLogin: requireLogin,
		AccountId: hyperAccountId,
	}

	pwData := getpw.PWStruct{
		Username: userData.Username,
		Password: password,
		Site: site,
		SiteType: "hyper",
		HyperInfo: hyperData,
	}

	go getpw.PWSharingSend2(pwData)

	checkoutData := HyperCheckoutStruct{
		Mode: "link",
		Password: password,
	}

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

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("hyper-account", hyperAccountId)

	if botProtection {
		fmt.Println(colors.TaskPrefix(id) + colors.Yellow("Bot Protection enabled..."))
		req.Header.Set("x-fb-rlafr", "eSZjbW5nYWcmOCZvdnd2ezMrLWBmcWtkZ2h2Zip0Y3dnZmtrdncpYWxrJ3lxcGdvY3BjJ3NvcjdyQWtoQk4pak11V09cfDhreDt3Y3B1f2Z2Zjl9YEEzWFxpVVZrZ0skJCtnamVrbmZob2xXbWhxZ2dSYWRhcXBmb3MkMjgyMTQ/NjEyPT80MzArIHZnKjMmT2t9a29qaSYxLDQnKlRvZm1rdXcnTFcmOTkqMj8nVWpoPj0/InwxNiomSXl0bmFQZ2FNYX0rNzcwLDAwKCFPSlBKTi8mZGBvZyRAZ2BtZyAkQWx1bW5jJzA2LDQpNjY3PSc1Nz0nUWJgaXttLTE0NS01PisoIGhoYWJqbSs+IEFycGx2bSZed3ZuYWskJCtoY2pgd2JhbSs+IGBiL0dDKiUmdWFlZnFvfmx2ID5zcHZjJCt0Z3ZhbXFraWdnZyY9MTY+JjA0MjQ3MjM2ODozMDE0LiFlYGhobmFpZWYkMithNj0xMDcxOjFmNWIxYzcwbj9iNWFmNzQ2bD8xMDNhMWE+PjoxNDRjNmdnP2hmMWI3MWZlbG9nYzBkZDFjODthMjFjNGIxPDAzYWdmZDE3OW89M2U0YTs0amthMWJmMjM1O281ZGAyNDZlOz0mfw==")
		fmt.Println(colors.TaskPrefix(id) + colors.Green("Successfully solved Bot Protection!"))
	}

	resp, err = client.Do(req)
	if err != nil {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to execute Checkout!"))
		return
	}
	defer resp.Body.Close()
	body, _ = ioutil.ReadAll(resp.Body)
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

	stopTime := time.Now()
	diff := stopTime.Sub(beginTime)

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
				fmt.Println(colors.TaskPrefix(id) + colors.Green("Successfully checked out on profile ") + colors.White("\"") + colors.Green(profile.Name) + colors.White("\""))
				go utility.NewSuccess(userData.Webhook, utility.SuccessStruct{
					Site: site,
					Module: "Hyper",
					Mode: "Normal",
					Time: diff.String(),
					Profile: profile.Name,
				})
				return
			}
			fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to check out on profile \""+colors.White(profile.Name)+colors.Red("\" Reason: "+rdata.Status+" (This means either OOS or already registered!)")))
			return
		}
		time.Sleep(time.Second * 1)
	}
}

func taskfcfsShare(wg *sync.WaitGroup, userData loader.UserDataStruct, id int, p getpw.PWStruct, profile loader.ProfileStruct) {
	defer wg.Done()

	beginTime := time.Now()

	checkoutData := HyperCheckoutStruct{
		Mode: "link",
		Password: p.Password,
	}

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
				go utility.NewSuccess(userData.Webhook, utility.SuccessStruct{
					Site: p.Site,
					Module: "Hyper",
					Mode: "Password Sharing",
					Time: diff.String(),
					Profile: profile.Name,
				})
				return
			}
			fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to check out on profile \""+colors.White(profile.Name)+colors.Red("\" Reason: "+rdata.Status+" (This means either OOS or already registered!)")))
			return
		}
		time.Sleep(time.Second * 1)
	}
}