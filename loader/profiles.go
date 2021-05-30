package loader

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/chrigeeel/satango/colors"
)

type BillingAddressStruct struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Line1    string `json:"line1"`
	Line2 string `json:"line2"`
	City string `json:"city"`
	PostCode string `json:"postCode"`
	Country  string `json:"country"`
}

type PaymentDetailsStruct struct {
	NameOnCard   string `json:"nameOnCard"`
	CardNumber   string `json:"cardNumber"`
	CardExpMonth string `json:"cardExpMonth"`
	CardExpYear  string `json:"cardExpYear"`
	CardCvv      string `json:"cardCvv"`
}

type ProfileStruct struct {
	Name           string `json:"name"`
	DiscordToken   string `json:"discordToken"`
	DiscordSession string `json:"discordSession,omitempty"`
	HyperUserId string`json:"hyperUserId,omitempty"`
	ShreyCookie string `json:"shreyCookie"`
	StripeToken string`json:"stripeToken,omitempty"`
	StripeToken2 string `json:"stripetoken2,omitempty"`
	BillingAddress BillingAddressStruct `json:"billingAddress"`
	PaymentDetails PaymentDetailsStruct `json:"paymentDetails"`
}

func LoadProfiles() []ProfileStruct {
	profilesFile, err := os.Open("./settings/profiles.json")
	if err != nil {
		fmt.Println(colors.Prefix() + colors.Red("No profiles set up!"))
		var profiles []ProfileStruct
		CreateProfile(profiles)
	}
	defer profilesFile.Close()
	profilesBytes, _ := ioutil.ReadAll(profilesFile)
	var profiles []ProfileStruct
	json.Unmarshal(profilesBytes, &profiles)
	return profiles
}

func CreateProfile(profiles []ProfileStruct) {
	fmt.Println(colors.Prefix() + colors.Yellow("Creating a new profile..."))
	var profile ProfileStruct
	profile.Name = askFor("Please input the name of the profile")
	profile.DiscordToken = askFor("Please input the Discord Token for this profile (Press Enter to skip)")
	profile.BillingAddress.Name = askFor("Please input your name")
	profile.BillingAddress.Email = askFor("Please input your email")
	profile.BillingAddress.Line1 = askFor("Please input your address")
	profile.BillingAddress.Line2 = askFor("Please input your address line two")
	profile.BillingAddress.City = askFor("Please input your city")
	profile.BillingAddress.PostCode = askFor("Please input your zip-/postcode")
	profile.BillingAddress.Country = askFor("Please input your country")
	profile.PaymentDetails.NameOnCard = askFor("Please input the name on your card")
	profile.PaymentDetails.CardNumber = askFor("Please input your card number")
	profile.PaymentDetails.CardExpMonth = askFor("Please input your card expiry month (format: \"02\")")
	profile.PaymentDetails.CardExpYear = askFor("Please input your card expiry year (format: \"2022\")")
	profile.PaymentDetails.CardCvv = askFor("Please input your CVV")
	profiles = append(profiles, profile)
	jsonData, _ := json.MarshalIndent(profiles, "", "    ")
	_ = ioutil.WriteFile("./settings/profiles.json", jsonData, 0644)
}

func askFor(question string) string {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println(colors.Prefix() + colors.White(question))
	fmt.Print(colors.Prefix() + colors.White("> "))
	scanner.Scan()
	return scanner.Text()
}
