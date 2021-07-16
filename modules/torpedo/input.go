package torpedo

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/asaskevich/govalidator"
	"github.com/chrigeeel/satango/colors"
	"github.com/chrigeeel/satango/loader"
	"github.com/chrigeeel/satango/modules/getpw"
	"github.com/chrigeeel/satango/utility"
)

func Input(userData loader.UserDataStruct, profiles []loader.ProfileStruct) {

	profiles = utility.AskForProfiles(profiles)

	profiles = login(profiles)

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
		key := getpw.GetPw2("torpedo").Password
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
				go TorpedoKeyTask(&wg, userData, i+1, key, profiles[profileCounter])
				profileCounter++
			}
			wg.Wait()
		}
	}
}