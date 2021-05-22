package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)


func main() {
	type LoadPageStruct struct {
		Props struct {
			Pageprops struct {
				Account   struct {
					Settings         struct {
						Payments struct {
							CollectBillingAddress bool   `json:"collect_billing_address"`
							RequireLogin          bool   `json:"require_login"`
						} `json:"payments"`
						BotProtection struct {
							Enabled bool `json:"enabled"`
						} `json:"bot_protection"`
					} `json:"settings"`
					ID string `json:"id"`
				} `json:"account"`
			} `json:"pageProps"`
		} `json:"props"`
		Query struct {
			Token string `json:"token"`
			Release string `json:"release"`
		} `json:"query"`
	}

	paid := false

	client := http.DefaultClient
	req, _ := http.NewRequest("GET", "https://dashboard.satanbots.com/purchase", nil)
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	r, _ := regexp.Compile("__NEXT_DATA__\" type=\"application\\/json\">({.*})")
	mdata := r.FindStringSubmatch(string(body))[1]

	fmt.Println(mdata)

	page := LoadPageStruct{}
	json.Unmarshal([]byte(mdata), &page)
	accountId := page.Props.Pageprops.Account.ID

	type loginStruct struct {
		Location string `json:"location"`
	}

	dUrl := "https://discord.com/api/v9/oauth2/authorize?client_id=648234176805470248&response_type=code&redirect_uri=https%3A%2F%2Fapi.hyper.co%2Fportal%2Fauth%2Fdiscord%2Fcallback&scope=identify%20email%20guilds%20guilds.join&state=%7B%22account%22%3A%22" + accountId + "%22%7D"
	token := "mfa.BUaom7Crxxf-UX9Ii6SBrS22fiXxtKGHX-yq-HCki4Tq1Q4ctx6SZKdmkCYhvq06_C51oXmB_y1VcIfkj6Iu"

	payload, _ := json.Marshal(map[string]string{
		"permssions": "0",
		"authorize": "true",
	})

	req, _ = http.NewRequest("POST", dUrl, bytes.NewBuffer(payload))
	req.Header.Set("content-type", "application/json")
	req.Header.Set("authorization", token)
	resp, _ = client.Do(req)
	defer resp.Body.Close()

	body, _ = ioutil.ReadAll(resp.Body)

	login := loginStruct{}
	json.Unmarshal([]byte(body), &login)
	code := login.Location
	r, _ = regexp.Compile("\\?code=(\\w*)")
	code = r.FindStringSubmatch(code)[1]
	fmt.Println(code)

	callBackUrl := "https://api.hyper.co/portal/auth/discord/callback?code=" + code + "&state=%7B%22account%22%3A%22" + accountId + "%22%7D"
	resp, _ = client.Get(callBackUrl)
	defer resp.Body.Close()
	body, _ = ioutil.ReadAll(resp.Body)

	r, _ = regexp.Compile("__NEXT_DATA__\" type=\"application\\/json\">({.*})")
	mdata = r.FindStringSubmatch(string(body))[1]

	json.Unmarshal([]byte(mdata), &page)
	hyperToken := page.Query.Token
	fmt.Println(hyperToken)

	type userStruct struct {
		ID string `json:"id"`
	}

	req, _ = http.NewRequest("GET", "https://dashboard.satanbots.com/ajax/user", nil)

	req.Header.Set("cookie", "authorization=" + hyperToken)
	req.Header.Set("hyper-account", accountId)

	resp, _ = client.Do(req)
	defer resp.Body.Close()

	user := userStruct{}

	body, _ = ioutil.ReadAll(resp.Body)
	json.Unmarshal([]byte(body), &user)

	hyperId := user.ID
	fmt.Println(hyperId)
	fmt.Println(paid)
	
	req, _ = http.NewRequest("GET", "https://dashboard.satanbots.com/purchase/?password=x6HhvwBebCIYBALLSOMGBALLS", nil)
	
	resp, _  = client.Do(req)
	defer resp.Body.Close()
	body, _ = ioutil.ReadAll(resp.Body)
	r, _ = regexp.Compile("__NEXT_DATA__\" type=\"application\\/json\">({.*})")
	mdata = r.FindStringSubmatch(string(body))[1]

	json.Unmarshal([]byte(mdata), &page)
	requireLogin := page.Props.Pageprops.Account.Settings.Payments.RequireLogin
	botProtection := page.Props.Pageprops.Account.Settings.BotProtection.Enabled
	collectBilling := page.Props.Pageprops.Account.Settings.Payments.CollectBillingAddress
	releaseId := page.Query.Release
	fmt.Println(releaseId, requireLogin, botProtection, collectBilling)

	type checkoutStruct struct {
		Billing_details struct {
			Address struct {
				City string `json:"city,omitempty"`
				Country string `json:"country,omitempty"`
				Line1 string `json:"line1,omitempty"`
				Postal_code string `json:"postal_code,omitempty"`
				State string `json:"state,omitempty"`
			} `json:"address,omitempty"`
			Email string `json:"email"`
			Name string `json:"name"`
		} `json:"billing_details"`
		Payment_method string `json:"payment_method"`
		Release string `json:"release"`
		User string `json:"user"`
	}

	checkoutData := checkoutStruct{}

	if paid == true {
		checkoutData.Payment_method = "Bruh"
	}
	if collectBilling == true {
		fmt.Println("bruh")
	}
	checkoutData.Release = releaseId
	checkoutData.User = hyperId

	checkoutData.Billing_details.Email = "chriguuuul@gmail.com"
	checkoutData.Billing_details.Name = "Christian Tognazza"

	payload, _ = json.Marshal(checkoutData)

	fmt.Println(string(payload))

	req, _ = http.NewRequest("POST", "https://dashboard.satanbots.com/ajax/checkouts", bytes.NewBuffer(payload))

	req.Header.Set("cookie", "authorization=" + hyperToken)
	if botProtection == true {
		req.Header.Set("x-amz-cf-id", "U2FsdGVkX1/oevSvIQbzfLKlXhAKdJATjezc0HURuWJo58if5XL7QCmp7OvG0vS57O3rxfmD0v9kBIU+Nb129dWlHoCg/S4mQAi1F4rD6YNi7o8ACcpW1tEGTO4Z7u+Cv0DDCxk6par1+IPIm3cIVWOZqrIQAMrxCeL9M/t8s3GGB+Y49eYwkqiHAcm/zkxbSFLrNv4DCiWXxrlfw/SqvbAaiyxXPivVNtHYaEMjlME1hOge8NrUZJCQM7cRWk5/EX+ghb8GWdXbHh+rSSGhyBbW+9GvPYm37BWUrzYxF7ktK9Xkf2wnNTUG0iX5jHVMNYORvnXHInTk4HySNMsznw8pI7fZenuGYlphw7j/kt14Rlr+VP4aNnquhgqu52GkDTUoNRwVwZKYcCQ1I74fmbTOewq0yQqyxs2pdMeV41QTQiVBgh/Xu6wzcp38XGSE806FC29tHzLupoeiLIVDOmjnrkeHsIKzgUzwTTdFV3Of2CpeGD+VAjptnJt4dSMnttc4/gXohyyoCMlB/5bC6BxHzO86gUy5oO6UwuwVRirvnW/1cE/n2+6Kg5uT8gW2ov7F0mizHPRX8woRprae7zZmgeAuZOjNJjBC6N8EkRoNCeZn3Zp/wHjFcJ7a5CAH68RtehNVlHupwDuW07NzOZNXCpi2aDbePaK/NW8+fOwpagRc4SeikQKq7HYdwBb5s/WD1oGwELkvyY886uHeaCVg8/tB4BMDcwP1IwVcGD6uLoLNVpURV3Fm77hGeEC1Xaet6o5MQXvB4SRn5NCwPrkx2PYqRK5LpTl6Aba3kQW7AgtOY1Ws+gt6SUtKSAYZyMgLkSdSyM60ekxbs84qlK254Bx3iI0Wm+2umev3SvH+bIyE5PPps3jsseaZRDM62wTrH70An9MlgTAgskV3ze/yZhB+fDXg+6lhabXc7pFOJp7QZwESlSbLwEU5qSMD")
	}
	req.Header.Set("hyper-account", accountId)
	req.Header.Set("content-type", "application/json")

	resp, _ = client.Do(req)
	defer resp.Body.Close()
	body, _ = ioutil.ReadAll(resp.Body)

	fmt.Println(string(body))
}