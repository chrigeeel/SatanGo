package shrey

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/chrigeeel/satango/colors"
	"github.com/chrigeeel/satango/loader"
	"github.com/chrigeeel/satango/modules/getpw"
	"github.com/chrigeeel/satango/utility"
)

func Input(userData loader.UserDataStruct, profiles []loader.ProfileStruct, proxies []string) {
	err := loader.AuthKeySilent(userData.Key)
	if err != nil {
		fmt.Println("")
		fmt.Println(colors.Prefix() + colors.Red("Failed to authenticate your key!"))
		fmt.Println(colors.Prefix() + colors.Red("Please contact staff!"))
		time.Sleep(time.Second * 10)
		os.Exit(3)
	}
	fmt.Println(colors.Prefix() + colors.Red("What site would you like to start tasks on?") + colors.White(" (example: \"dashboard.satanbots.com\")"))
	site := utility.AskForSilent()
	r := regexp.MustCompile(`[^\/]*\.[^\/]*\.?[^\/]*`)
	siteB := r.Find([]byte(site))
	if siteB == nil {
		fmt.Println(colors.Prefix() + colors.Red("Invalid site input!"))
		return
	}
	site = "https://" + string(siteB) + "/"

	getpw.AskForPwShare()

	profiles = utility.AskForProfiles(profiles)

	if len(profiles) == 0 {
		fmt.Println(colors.Prefix() + colors.Red("You have no valid Profiles! Please check them and their corresponding Payment Info!"))
		time.Sleep(time.Second * 3)
		return
	}
	fmt.Println(colors.Prefix() + colors.Yellow("Logging in on all profiles..."))
	profiles = login(profiles, site)
	if len(profiles) == 0 {
		fmt.Println(colors.Prefix() + colors.Red("You have no valid Profiles! Please check them and their corresponding Discord Tokens!"))
		time.Sleep(time.Second * 3)
		return
	}
	fmt.Println(colors.Prefix() + colors.Yellow("Attaching Payment Method on all profiles..."))
	profiles = stripe(profiles, site)

	taskLimit := len(profiles)

	fmt.Println(colors.Prefix() + colors.Red("How many tasks do you want to run? Your task limit is ") + colors.White(strconv.Itoa(taskLimit)) + colors.Red(" because you have ") + colors.White(strconv.Itoa(len(profiles))) + colors.Red(" valid profile(s)"))

	var taskAmount int
	for validAns := false; !validAns; {
		ans := utility.AskForSilent()
		validAns = true
		ansInt, _ := strconv.Atoi(ans)
		if !govalidator.IsInt(ans) {
			fmt.Println(colors.Prefix() + colors.Red("Wrong input!"))
			validAns = false
		} else if ansInt > taskLimit {
			fmt.Println(colors.Prefix() + colors.Red("Your task limit is ") + colors.White(strconv.Itoa(taskLimit)) + colors.Red("!"))
			validAns = false
		}
		taskAmount = ansInt
	}

	var profileCounter int

	for exit := false; !exit; {
		p := getpw.GetPw2(site)
		password := p.Password
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
				go taskfcfs(&wg, userData, i+1, p, profiles[profileCounter], site)
				profileCounter++
			}
			wg.Wait()
		}
	}
}