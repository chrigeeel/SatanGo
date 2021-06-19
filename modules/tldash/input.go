package tldash

import (
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
	"github.com/chrigeeel/satango/modules/getpw"
	"github.com/chrigeeel/satango/modules/utility"
)

func Input(userData loader.UserDataStruct, profiles []loader.ProfileStruct, proxies []string, mode string) {
	fmt.Println(colors.Prefix() + colors.Yellow("Loading sites..."))
	resp, err := http.Get("http://50.16.47.99/sites/tl")
	if err != nil {
		fmt.Println(colors.Prefix() + colors.Red("Error loading sites! exiting..."))
		time.Sleep(time.Second * 3)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(colors.Prefix() + colors.Red("Error loading sites! exiting..."))
		time.Sleep(time.Second * 3)
		return
	}
	var sites []siteStruct
	json.Unmarshal(body, &sites)
	fmt.Println(colors.Prefix() + colors.Red("What site would you like to start tasks on?"))
	for i := range sites {
		fmt.Println(colors.Prefix() + colors.White("["+strconv.Itoa(i+1)+"] "+sites[i].DisplayName))
	}
	fmt.Println(colors.Prefix() + colors.White("[%] Go back"))
	var ansInt int
	for validAns := false; !validAns; {
		ans := utility.AskForSilent()
		if ans == "%" {
			return
		}
		validAns = true
		ansInt, _ = strconv.Atoi(ans)
		if !govalidator.IsInt(ans) {
			fmt.Println(colors.Prefix() + colors.Red("Wrong input!"))
			validAns = false
		} else if ansInt > len(sites) {
			fmt.Println(colors.Prefix() + colors.Red("Wrong input!"))
			validAns = false
		}
	}
	site := sites[ansInt-1]
	fmt.Println(colors.Prefix() + colors.Red("(Y/N) Do you want to enable Discord-Login? (Recommended for pure TL-Sites, if unsure ask support)"))
	ans := utility.AskForSilent()
	var discordLogin bool
	if strings.ToLower(ans)[0:1] == "y" {
		discordLogin = true
		fmt.Println(colors.Prefix() + colors.White("Discord login turned ") + colors.Green("On"))
	} else if strings.ToLower(ans)[0:1] == "n" {
		discordLogin = false
		fmt.Println(colors.Prefix() + colors.White("Discord login turned ") + colors.Red("Off"))
	}
	if discordLogin {
		profiles = login(profiles)
	}
	if len(profiles) == 0 {
		fmt.Println(colors.Prefix() + colors.Red("You have no valid Profiles! Please check them and their corresponding Discord Tokens!"))
		time.Sleep(time.Second * 3)
		return
	}
	if site.Stripe_public_key != "" {
		profiles = stripe(site.Stripe_public_key, profiles)
	}
	if len(profiles) == 0 {
		fmt.Println(colors.Prefix() + colors.Red("You have no valid Profiles! Please check them and their corresponding Payment Info!"))
		time.Sleep(time.Second * 3)
		return
	}

	profiles = utility.AskForProfiles(profiles)

	fmt.Println(colors.Prefix() + colors.Yellow("Locating best TL API Server..."))
	solveIp := findApi()

	if mode == "RAFFLE" {
		inputraffle(userData, profiles, proxies, discordLogin, site, solveIp)
		return
	}

	var taskLimit int
	if len(proxies) > 20 {
		taskLimit = 40
	} else {
		taskLimit = len(proxies) * 2
	}
	if len(proxies) == 1 {
		fmt.Println(colors.Prefix() + colors.Red("How many tasks do you want to run? You have ") + colors.White(strconv.Itoa(len(proxies))) + colors.Red(" proxy and your task limit is ") + colors.White(strconv.Itoa(taskLimit)) + colors.Red("."))
	} else {
		fmt.Println(colors.Prefix() + colors.Red("How many tasks do you want to run? You have ") + colors.White(strconv.Itoa(len(proxies))) + colors.Red(" proxies and your task limit is ") + colors.White(strconv.Itoa(taskLimit)) + colors.Red("."))
	}
	var taskAmount int
	for validAns := false; !validAns; {
		ans = utility.AskForSilent()
		validAns = true
		ansInt, _ = strconv.Atoi(ans)
		if !govalidator.IsInt(ans) {
			fmt.Println(colors.Prefix() + colors.Red("Wrong input!"))
			validAns = false
		} else if ansInt > taskLimit {
			fmt.Println(colors.Prefix() + colors.Red("Your task limit is ") + colors.White(strconv.Itoa(taskLimit)) + colors.Red("!"))
			validAns = false
		}
		taskAmount = ansInt
	}
	var tasks []taskStruct
	var profileCounter int
	var proxyCounter int
	for i := 0; i < taskAmount; i++ {
		if profileCounter+1 > len(profiles) {
			profileCounter = 0
		}
		if proxyCounter+1 > len(proxies) {
			proxyCounter = 0
		}
		tasks = append(tasks, taskStruct{
			Site:    site.BackendName,
			Proxy:   proxies[proxyCounter],
			Profile: profiles[profileCounter],
		})
		proxyCounter++
		profileCounter++
	}
	for exit := false; !exit; {
		password := getpw.GetPw2(site.BackendName).Password
		if password == "exit" {
			exit = true
		} else {
			fmt.Println(colors.Prefix() + colors.Yellow("Starting tasks..."))
			//go StartLucifer(userData, password, solveIp, tasks[0].Site, tasks[0].Proxy, tasks[0].Profile, tasks[0].Profile.StripeToken)
			//go StartLucifer(userData, password, solveIp, tasks[0].Site, tasks[0].Proxy, tasks[0].Profile, tasks[0].Profile.StripeToken)
			var wg sync.WaitGroup
			for i := range tasks {
				if len(tasks) > i {
					wg.Add(1)
					go taskfcfs(&wg, userData, i+1, password, solveIp, tasks[i])
				}
			}
			wg.Wait()
		}
	}
}