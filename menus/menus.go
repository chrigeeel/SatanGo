package menus

import (
	"bufio"
	"fmt"
	"os"

	"github.com/chrigeeel/satango/colors"
	"github.com/chrigeeel/satango/loader"
	"github.com/chrigeeel/satango/modules"
)

func MainMenu(userData loader.UserDataStruct, profiles []loader.ProfileStruct, proxies []string) {
	fmt.Println(colors.Red("------------------------------------------------------------------"))
	fmt.Println(colors.Prefix() + colors.Red("What would you like to do?"))
	fmt.Println(colors.Prefix() + colors.White("[1] Start the TL Module"))
 	fmt.Println(colors.Prefix() + colors.White("[2] Start the Hyper/Meta Labs Module"))
	fmt.Println(colors.Prefix() + colors.White("[%] Create a new Profile"))
	ans := askForSilent()
	proxies = loader.LoadProxies()
	profiles = loader.LoadProfiles()
	switch ans {
	case "1":
		fmt.Println(colors.Prefix() + colors.Red("Would you like to start the Raffle or FCFS module?"))
		fmt.Println(colors.Prefix() + colors.White("[1] Start the TL FCFS Module"))
		fmt.Println(colors.Prefix() + colors.White("[2] Start the TL Raffle Module"))
		ans = askForSilent()
		switch ans {
		case "1":
			modules.TLInput(userData, profiles, proxies, "FCFS")
			MainMenu(userData, profiles, proxies)
		case "2":
			modules.TLInput(userData, profiles, proxies, "RAFFLE")
			MainMenu(userData, profiles, proxies)
		default:
			fmt.Println(colors.Prefix() + colors.Red("Invalid answer!"))
			MainMenu(userData, profiles, proxies)
		}
	case "2":
		modules.HyperInput(userData, profiles, proxies)
		MainMenu(userData, profiles, proxies)
	case "%":
		loader.CreateProfile(profiles)
		MainMenu(userData, profiles, proxies)
	default:
		fmt.Println(colors.Prefix() + colors.Red("Invalid answer!"))
		MainMenu(userData, profiles, proxies)
	}

}

func askForSilent() string {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(colors.Prefix() + colors.White("> "))
	scanner.Scan()
	return scanner.Text()
}
