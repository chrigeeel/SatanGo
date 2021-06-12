package main

import (
	"fmt"
	"time"

	"github.com/chrigeeel/satango/colors"
	"github.com/chrigeeel/satango/loader"
	"github.com/chrigeeel/satango/menus"
	"github.com/chrigeeel/satango/modules"
)

var (
	Version string = "0.5.25"
)

func main() {

	go modules.Server()

	time.Sleep(time.Millisecond * 10)

	menus.CallClear()

	fmt.Println(colors.Prefix() + colors.White("Starting... | Version " + Version))

	userData := loader.LoadSettings()
	auth := loader.AuthKey(userData.Key)
	username := auth.User.Username
	userData.Username = username
	userData.DiscordId = auth.User.DiscordId
	userData.Version = Version
	fmt.Println(colors.Prefix() + colors.White("Welcome back, ") + colors.Red(username) + "!")
	loader.LoadProfiles()

	loader.CheckForUpdate(userData)

	fmt.Println(colors.Prefix() + colors.Yellow("Loading your data..."))
	
	profiles := loader.LoadProfiles()
	proxies := loader.LoadProxies()

	fmt.Println(colors.Prefix() + colors.Green("Successfully loaded profiles, proxies and tokens"))
	
	go modules.PwSharingReceive()
	time.Sleep(time.Millisecond * 500)

	menus.MainMenu(userData, profiles, proxies)
}
