package keyclaimer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/chrigeeel/satango/colors"
	"github.com/chrigeeel/satango/loader"
	"github.com/chrigeeel/satango/modules/getpw"
	"github.com/chrigeeel/satango/modules/wrath"
	"github.com/chrigeeel/satango/utility"
)

type siteStruct struct {
	DisplayName string `json:"displayName"`
	ClientID string `json:"clientId"`
	URL string `json:"url"`
}

func Input(userData loader.UserDataStruct, profiles []loader.ProfileStruct) {
	fmt.Println(colors.Prefix() + colors.Yellow("Loading sites..."))
	err := loader.AuthKeySilent(userData.Key)
	if err != nil {
		fmt.Println("")
		fmt.Println(colors.Prefix() + colors.Red("Failed to authenticate your key!"))
		fmt.Println(colors.Prefix() + colors.Red("Please contact staff!"))
		time.Sleep(time.Second * 10)
		os.Exit(3)
	}
	resp, err := http.Get("https://hardcore.astolfoporn.com/api/sites/freddy")
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
	fmt.Println(colors.Prefix() + colors.White("["+strconv.Itoa(len(sites)+1)+"] "+"WrathAIO"))
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
			if ansInt == len(sites) + 1 {
				wrath.Input(userData, profiles)
				return
			} 
			fmt.Println(colors.Prefix() + colors.Red("Wrong input!"))
			validAns = false
		}
	}
	site := sites[ansInt-1]

	profiles = utility.AskForProfiles(profiles)

	profiles = login(profiles, site)

	var taskLimit int
	taskLimit = len(profiles)
	if taskLimit > 15 {
		taskLimit = 15
	}
	fmt.Println(colors.Prefix() + colors.Red("How many tasks do you want to run? You have ") + colors.White(strconv.Itoa(len(profiles))) + colors.Red(" profiles and your task limit is ") + colors.White(strconv.Itoa(taskLimit)) + colors.Red("."))
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
		key := getpw.GetPw2(site.DisplayName).Password
		if key == "exit" {
			exit = true
		} else {
			fmt.Println(colors.Prefix() + colors.Yellow("Starting tasks..."))
			var wg sync.WaitGroup
			for i := 0; i < taskAmount; i++ {
				if profileCounter+1 > len(profiles) {
					profileCounter = 0
				}
				wg.Add(1)
				go taskfcfs(&wg, userData, i+1, key, profiles[profileCounter], site)
				profileCounter++
			}
			wg.Wait()
		}
	}
}