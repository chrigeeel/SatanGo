package hyper

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
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
	fmt.Println(colors.Prefix() + colors.Red("(Y/N) Is the release you're going for paid?"))
	ans := utility.AskForSilent()
	var paid bool
	if strings.ToLower(ans)[0:1] == "y" {
		paid = true
		fmt.Println(colors.Prefix() + colors.White("Release is paid"))
	} else {
		paid = false
		fmt.Println(colors.Prefix() + colors.White("Release is not paid"))
	}
	fmt.Println(colors.Prefix() + colors.Yellow("Loading site..."))
	loadPage, err := load(site)
	if err != nil {
		fmt.Println(colors.Prefix() + colors.Red("Failed to load site "+site))
		return
	}
	fmt.Println(colors.Prefix() + colors.Green("Successfully loaded site"))

	fmt.Println(colors.Prefix() + colors.Yellow("Trying to solve Bot Protection..."))
	bpToken, err := Solvebp(site)
	if err != nil {
		fmt.Println(colors.Prefix() + colors.Red("Failed to solve Bot Protection. Please contact Chrigeeel or Shrek!"))
	} else {
		fmt.Println(colors.Prefix() + colors.Green("Successfully solved Bot Protection!"))
	}

	getpw.AskForPwShare()

	profiles = utility.AskForProfiles(profiles)

	if paid {
		fmt.Println(colors.Prefix() + colors.Yellow("Fetching stripe tokens for all Profiles..."))
		profiles = stripe(loadPage.Props.Pageprops.Account.Stripe_account, profiles)
	}
	if len(profiles) == 0 {
		fmt.Println(colors.Prefix() + colors.Red("You have no valid Profiles! Please check them and their corresponding Payment Info!"))
		time.Sleep(time.Second * 3)
		return
	}
	fmt.Println(colors.Prefix() + colors.Yellow("Logging in on all profiles..."))
	profiles = Login(loadPage.Props.Pageprops.Account.ID, site, profiles)
	if len(profiles) == 0 {
		fmt.Println(colors.Prefix() + colors.Red("You have no valid Profiles! Please check them and their corresponding Discord Tokens!"))
		time.Sleep(time.Second * 3)
		return
	}

	var taskLimit int = len(profiles) * 2

	fmt.Println(colors.Prefix() + colors.Red("How many tasks do you want to run? Your task limit is ") + colors.White(strconv.Itoa(taskLimit)) + colors.Red(" because you have ") + colors.White(strconv.Itoa(len(profiles))) + colors.Red(" valid profiles"))

	var taskAmount int
	for validAns := false; !validAns; {
		ans = utility.AskForSilent()
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
		releaseId := p.HyperInfo.ReleaseId
		if password == "exit" {
			exit = true
		} else {
			fmt.Println(colors.Prefix() + colors.Yellow("Starting tasks..."))
			var wg sync.WaitGroup
			for i := 0; i < taskAmount; i++ {
				if profileCounter+1 > len(profiles) {
					profileCounter = 0
				}
				if p.Mode == "share" {
					wg.Add(1)
					go taskfcfsShare(&wg, userData, i+1, p, profiles[profileCounter])
				} else {
					if releaseId != "" && (userData.Key == "SATAN-PRSF-JLE5-ICJU-M85H" || userData.Key == "GCGK-T824-E6CC-DUBG") {
						wg.Add(1)
						go taskfcfsCool(&wg, userData, i+1, releaseId, paid, false, site, loadPage.Props.Pageprops.Account.ID, profiles[profileCounter], bpToken)
						wg.Add(1)
						go taskfcfsCool(&wg, userData, i+1, releaseId, paid, true, site, loadPage.Props.Pageprops.Account.ID, profiles[profileCounter], bpToken)
					} else {
						wg.Add(1)
						go taskfcfs(&wg, userData, i+1, password, paid, site, loadPage.Props.Pageprops.Account.ID, profiles[profileCounter], bpToken)
					}
				}
				profileCounter++
			}
			wg.Wait()
		}
	}
}