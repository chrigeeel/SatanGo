package afk

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/chrigeeel/satango/colors"
	"github.com/chrigeeel/satango/loader"
	"github.com/chrigeeel/satango/modules/getpw"
	"github.com/chrigeeel/satango/modules/hyper"
	"github.com/chrigeeel/satango/utility"
)


func Input(userData loader.UserDataStruct, profiles []loader.ProfileStruct) {
	profiles = utility.AskForProfiles(profiles)
	fmt.Println(colors.Prefix() + colors.Yellow("Logging in on all profiles..."))
	profiles = hyper.Login("RQyuOz8FmunrZMrdkwBaC", "https://botsandmonitors.metalabs.gg/", profiles)
	if len(profiles) == 0 {
		fmt.Println(colors.Prefix() + colors.Red("You have no valid Profiles! Please check them and their corresponding Discord Tokens!"))
		time.Sleep(time.Second * 3)
		return
	}
	var taskLimit int = len(profiles) * 2

	fmt.Println(colors.Prefix() + colors.Red("How many AFK tasks do you want to run? Your task limit is ") + colors.White(strconv.Itoa(taskLimit)) + colors.Red(" because you have ") + colors.White(strconv.Itoa(len(profiles))) + colors.Red(" valid profiles"))
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

	for {
		p := getpw.GetPwAfk()
		if p.Password == "exit" {
			return
		}
		if p.SiteType == "hyper" {
			fmt.Println(colors.Prefix() + colors.Yellow("Starting tasks..."))
			var wg sync.WaitGroup
			for i := 0; i < taskAmount; i++ {
				if profileCounter+1 > len(profiles) {
					profileCounter = 0
				}
				wg.Add(1)
				go taskfcfs(&wg, userData, i+1, p, profiles[profileCounter])
				profileCounter++
			}
			wg.Wait()
		}
		if p.Mode == "discord" {
			fmt.Println(colors.Prefix() + colors.Yellow("Starting tasks..."))
			var wg sync.WaitGroup
			for i := 0; i < taskAmount; i++ {
				if profileCounter+1 > len(profiles) {
					profileCounter = 0
				}
				wg.Add(1)
				go taskfcfsunk(&wg, userData, i+1, p.Password, p.Site, profiles[profileCounter])
				profileCounter++
			}
			wg.Wait()
		}
	}
}