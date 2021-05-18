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
	fmt.Println(colors.Red("----------------------------------------------------------------"))
	fmt.Println(colors.Prefix() + colors.Red("What would you like to do?"))
	fmt.Println(colors.Prefix() + colors.White("[1] Start the TL Module"))
	//fmt.Println(colors.Prefix() + colors.White("[2] Start the Hyper/Meta Labs Module"))
	fmt.Println(colors.Prefix() + colors.White("[%] Create a new Profile"))
	ans := askForSilent()
	switch ans {
	case "1":
		modules.TLInput(userData, profiles, proxies)
		MainMenu(userData, profiles, proxies)
		/*
	case "2":
		modules.HyperInput(userData, profiles, proxies)
		MainMenu(userData, profiles, proxies)
		*/
	case "%":
		loader.CreateProfile(profiles)
		MainMenu(userData, profiles, proxies)
	default:
		fmt.Println("Wrong bruh")
	}

}

func askForSilent() string {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(colors.Prefix() + colors.White("> "))
	scanner.Scan()
	return scanner.Text()
}
