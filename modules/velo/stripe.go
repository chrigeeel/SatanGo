package velo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"github.com/chrigeeel/satango/colors"
	"github.com/chrigeeel/satango/loader"
	"github.com/chrigeeel/satango/modules/utility"
)

func stripe(profiles []loader.ProfileStruct, stripeKey string) []loader.ProfileStruct {

	var wg sync.WaitGroup

	tokenLocal := func(wg *sync.WaitGroup, id int) {
		defer wg.Done()
		type tokenStruct struct {
			ID string `json:"id"`
		}

		client := &http.Client{}
		url := "https://api.stripe.com/v1/payment_methods"
		if len(profiles[id].PaymentDetails.CardExpYear) < 2 {
			profiles[id].StripeToken = "error"
			fmt.Println(colors.Prefix() + colors.Red("Stripe rejected your profile ") + colors.White(profiles[id].Name) + colors.Red(" - Not running tasks on this profile"))
			return
		}
		payload := strings.NewReader(
			`type=card` +
				`&card[number]=` + profiles[id].PaymentDetails.CardNumber +
				`&card[cvc]=` + profiles[id].PaymentDetails.CardCvv +
				`&card[exp_month]=` + profiles[id].PaymentDetails.CardExpMonth +
				`&card[exp_year]=` + profiles[id].PaymentDetails.CardExpYear[len(profiles[id].PaymentDetails.CardExpYear)-2:] +
				`&billing_details[name]=` + profiles[id].BillingAddress.Name +
				`&billing_details[email]=` + profiles[id].BillingAddress.Email +
				`&billing_details[address][country]=` + "US" +
				`&billing_details[address][line1]=` + profiles[id].BillingAddress.Line1 +
				`&billing_details[address][city]=` + profiles[id].BillingAddress.City +
				`&billing_details[address][postal_code]=` + profiles[id].BillingAddress.PostCode +
				`&key=` + stripeKey)

		req, err := http.NewRequest("POST", url, payload)
		if err != nil {
			profiles[id].StripeToken = "error"
			fmt.Println(colors.Prefix() + colors.Red("Stripe rejected your profile ") + colors.White(profiles[id].Name) + colors.Red(" - Not running tasks on this profile"))
			return
		}

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := client.Do(req)
		if err != nil {
			profiles[id].StripeToken = "error"
			fmt.Println(colors.Prefix() + colors.Red("Stripe rejected your profile ") + colors.White(profiles[id].Name) + colors.Red(" - Not running tasks on this profile"))
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			profiles[id].StripeToken = "error"
			fmt.Println(colors.Prefix() + colors.Red("Stripe rejected your profile ") + colors.White(profiles[id].Name) + colors.Red(" - Not running tasks on this profile"))
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
		fmt.Println(colors.Prefix() + colors.Green("Successfully fetched Stripe token one for profile ") + colors.White(profiles[id].Name))
	}

	for id := range profiles {
		wg.Add(1)
		go tokenLocal(&wg, id)
	}
	wg.Wait()

	for i := len(profiles) - 1; i >= 0; {
		if profiles[i].StripeToken == "error" {
			profiles = utility.RemoveIndex(profiles, i)
		} else if profiles[i].StripeToken2 == "error" {
			profiles = utility.RemoveIndex(profiles, i)
		}
		i--
	}

	return profiles
}