package tldash

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/chrigeeel/satango/colors"
	"github.com/chrigeeel/satango/loader"
	"github.com/chrigeeel/satango/utility"
)

func inputraffle(userData loader.UserDataStruct, profiles []loader.ProfileStruct, proxies []string, discordLogin bool, site siteStruct, solveIp string) {
	fmt.Println(colors.Prefix() + colors.White("You have ") + colors.Red(strconv.Itoa(len(proxies))+colors.White(" Proxies including localhost!")))
	if len(profiles) >= len(proxies) {
		fmt.Println(colors.Prefix() + colors.White("That means you can only run ") + colors.Red(strconv.Itoa(len(proxies))) + colors.White(" Profiles! Please enter more Proxies in proxies.txt"))
		var newProfiles []loader.ProfileStruct
		var toPrint string
		toPrint = colors.Prefix() + colors.Red("Only running the profiles: \n") + colors.Prefix()
		for i := range proxies {
			newProfiles = append(newProfiles, profiles[i])
			toPrint = toPrint + colors.White("\"") + colors.Red(profiles[i].Name) + colors.White("\", ")
		}
		profiles = newProfiles
		fmt.Println(toPrint)
	} else {
		fmt.Println(colors.Prefix() + colors.Green("That means you can run all profiles!"))
	}
	taskAmount := len(profiles)
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
		fmt.Println(colors.Prefix() + colors.Red("Please enter the raffle password:"))
		fmt.Println(colors.Prefix() + colors.Red("Enter ") + colors.White("\"") + colors.Red("exit") + colors.White("\"") + colors.Red(" to exit"))
		password := utility.AskForSilent()
		if password == "exit" {
			exit = true
		} else {
			fmt.Println(colors.Prefix() + colors.Yellow("Starting tasks..."))
			var wg sync.WaitGroup
			for i := range tasks {
				wg.Add(1)
				go taskraffle(&wg, userData, i+1, password, solveIp, tasks[i], false)
			}
			wg.Wait()
		}
	}
}