package velo

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/chrigeeel/satango/colors"
	"github.com/chrigeeel/satango/loader"
	"github.com/chrigeeel/satango/modules/getpw"
	"github.com/chrigeeel/satango/utility"
)

type siteStruct struct {
	DisplayName       string `json:"displayName"`
	BackendName       string `json:"backendName"`
	Stripe_public_key string `json:"stripe_public_key,omitempty"`
}

var veloSites = []siteStruct{
	{
		DisplayName:       "ClipAIO",
		BackendName:       "dash.clipaio.com",
		Stripe_public_key: "pk_live_51JD4fxDAtQxFOfqOmxO8K2jg0guFi7zUx4jqK2TKzw7PfXqI9RJ2qYkFaQ75XPa063E6cCf4d9Gl8FB8rSMtGUOs00koDgBjMP",
	},
	{
		DisplayName:       "CryptoClub",
		BackendName:       "dash.cryptoclub.group",
		Stripe_public_key: "pk_live_51IibakGMt5G1CmqPqzOg6k2RKavZa2PF2BoQ0c5BM53GKANiJCGq7CZwL1uHAzjjcAYD4jaCYsYtsa6M3QV7ivMX00n8thUefS",
	},
	{
		DisplayName:       "Lumen",
		BackendName:       "dash.stellaraio.group",
		Stripe_public_key: "pk_live_51Idj8iGNqrK0CQYirqWCb6jK0ZlrW30HaUMlTg6YEm0lvUaNtktEerPM5DOSkLKcQcU345Ipj41D1VmlwXkWFwPi00Jz3GrJCS",
	},
	{
		DisplayName:       "Velo",
		BackendName:       "versusfnf.io",
		Stripe_public_key: "pk_live_Zztx90CWU6pPnsRBZcWm1DSE00s2fTsNbc",
	},
}

func Input(userData loader.UserDataStruct, profiles []loader.ProfileStruct, proxies []string) {
	err := loader.AuthKeySilent(userData.Key)
	if err != nil {
		fmt.Println("")
		fmt.Println(colors.Prefix() + colors.Red("Failed to authenticate your key!"))
		fmt.Println(colors.Prefix() + colors.Red("Please contact staff!"))
		time.Sleep(time.Second * 10)
		os.Exit(3)
	}
	fmt.Println(colors.Prefix() + colors.Red("What site would you like to start tasks on?"))
	for i := range veloSites {
		fmt.Println(colors.Prefix() + colors.White("["+strconv.Itoa(i+1)+"] "+veloSites[i].DisplayName))
	}
	fmt.Println(colors.Prefix() + colors.White("[%] Go back"))
	var ansInt int
	for validAns := false; validAns == false; {
		ans := utility.AskForSilent()
		if ans == "%" {
			return
		}
		validAns = true
		ansInt, _ = strconv.Atoi(ans)
		if govalidator.IsInt(ans) == false {
			fmt.Println(colors.Prefix() + colors.Red("Wrong input!"))
			validAns = false
		} else if ansInt > len(veloSites) {
			fmt.Println(colors.Prefix() + colors.Red("Wrong input!"))
			validAns = false
		}
	}
	site := veloSites[ansInt-1]

	getpw.AskForPwShare()

	profiles = utility.AskForProfiles(profiles)

	fmt.Println(colors.Prefix() + colors.Yellow("Loggin in on all profiles..."))
	profiles = login(profiles, site.BackendName)
	if len(profiles) == 0 {
		fmt.Println(colors.Prefix() + colors.Red("You have no valid Profiles! Please check them and their corresponding Discord Tokens!"))
		time.Sleep(time.Second * 3)
		return
	}
	fmt.Println(colors.Prefix() + colors.Yellow("Fetching Stripe on all profiles..."))
	profiles = stripe(profiles, site.Stripe_public_key)
	if len(profiles) == 0 {
		fmt.Println(colors.Prefix() + colors.Red("You have no valid Profiles! Please check them and their corresponding Discord Tokens!"))
		time.Sleep(time.Second * 3)
		return
	}

	var taskLimit int
	taskLimit = len(profiles) * 4

	fmt.Println(colors.Prefix() + colors.Red("How many tasks do you want to run? Your task limit is ") + colors.White(strconv.Itoa(taskLimit)) + colors.Red(" because you have ") + colors.White(strconv.Itoa(len(profiles))) + colors.Red(" valid profiles"))

	var taskAmount int
	for validAns := false; validAns == false; {
		ans := utility.AskForSilent()
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
		password := getpw.GetPw2(site.BackendName).Password
		if password == "exit" {
			exit = true
		} else {
			fmt.Println(colors.Prefix() + colors.Yellow("Starting tasks..."))
			var wg sync.WaitGroup
			for i := 0; i < taskAmount; i++ {
				if profileCounter+1 > len(profiles) {
					profileCounter = 0
				}
				wg.Add(1)
				go VeloTask(&wg, userData, i+1, site, password, profiles[profileCounter])
				profileCounter++
			}
			wg.Wait()
		}
	}
}