package main

import (
	"fmt"
	"time"

	"github.com/chrigeeel/satango/colors"
	"github.com/chrigeeel/satango/loader"
	"github.com/chrigeeel/satango/menus"
	getpw "github.com/chrigeeel/satango/modules/getpw"
)

var (
	Version string = "0.5.29"
)

func main() {

	go getpw.MonitorClipboard()
	go getpw.MonitorExtension()

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
	
	go getpw.PWShareConnectReceive(userData)

	time.Sleep(time.Millisecond * 500)

	menus.MainMenu(userData, profiles, proxies)
}
