package menus

import (
	"fmt"

	"github.com/chrigeeel/satango/colors"
	"github.com/chrigeeel/satango/loader"
	"github.com/chrigeeel/satango/modules/afk"
	"github.com/chrigeeel/satango/modules/hyper"
	"github.com/chrigeeel/satango/modules/keyclaimer"
	"github.com/chrigeeel/satango/modules/shinobi"
	"github.com/chrigeeel/satango/modules/shrey"
	"github.com/chrigeeel/satango/modules/tldash"
	"github.com/chrigeeel/satango/modules/velo"
	"github.com/chrigeeel/satango/utility"
)

func MainMenu(userData loader.UserDataStruct) {
	fmt.Println(colors.Red("------------------------------------------------------------------"))
	fmt.Println(colors.Prefix() + colors.Red("What would you like to do?"))
	fmt.Println(colors.Prefix() + colors.White("[1] Start the TL Module"))
 	fmt.Println(colors.Prefix() + colors.White("[2] Start the Hyper/Meta Labs Module"))
	fmt.Println(colors.Prefix() + colors.White("[3] Start the Shrey Module"))
	fmt.Println(colors.Prefix() + colors.White("[4] Start the Velo Module"))
	fmt.Println(colors.Prefix() + colors.White("[5] Start the Key Claimer Modules"))
	fmt.Println(colors.Prefix() + colors.White("[6] Start the Shinobi Module"))
	fmt.Println(colors.Prefix() + colors.White("[7] Start FREE AFK Module"))
	fmt.Println(colors.Prefix() + colors.White("[%] Create a new Profile"))
	ans := utility.AskForSilent()
	proxies := loader.LoadProxies()
	profiles := loader.LoadProfiles()
	hyper.Initbp()
	switch ans {
	case "1":
		fmt.Println(colors.Prefix() + colors.Red("Would you like to start the Raffle or FCFS module?"))
		fmt.Println(colors.Prefix() + colors.White("[1] Start the TL FCFS module"))
		fmt.Println(colors.Prefix() + colors.White("[2] Start the TL Raffle module"))
		ans = utility.AskForSilent()
		switch ans {
		case "1":
			tldash.Input(userData, profiles, proxies, "FCFS")
			MainMenu(userData)
		case "2":
			tldash.Input(userData, profiles, proxies, "RAFFLE")
			MainMenu(userData)
		default:
			fmt.Println(colors.Prefix() + colors.Red("Invalid answer!"))
			MainMenu(userData)
		}
	case "2":
		hyper.Input(userData, profiles, proxies)
		MainMenu(userData)
	case "3":
		shrey.Input(userData, profiles, proxies)
		MainMenu(userData)
	case "4":
		velo.Input(userData, profiles, proxies)
		MainMenu(userData)
	case "5":
		keyclaimer.Input(userData, profiles)
		MainMenu(userData)
	case "6":
		shinobi.Input(userData, profiles, proxies)
		MainMenu(userData)
	case "7":
		afk.Input(userData, profiles)
		MainMenu(userData)
	case "%":
		loader.CreateProfile(profiles)
		MainMenu(userData)
	default:
		fmt.Println(colors.Prefix() + colors.Red("Invalid answer!"))
		MainMenu(userData)
	}

}
