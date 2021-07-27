package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/chrigeeel/satango/colors"
	"github.com/chrigeeel/satango/loader"
	"github.com/chrigeeel/satango/menus"
	"github.com/chrigeeel/satango/modules/getpw"
)

var (
	Version string = "0.5.70"
)

func main() {
	go getpw.MonitorExtension()

	time.Sleep(time.Millisecond * 10)

	menus.CallClear()

	fmt.Println(colors.Prefix() + colors.White("Starting... | Version " + Version))

	userData := loader.LoadSettings()
	auth, err := loader.AuthKeyNew(userData.Key, false)
	if err != nil {
		fmt.Println(colors.Prefix() + colors.Red("Exiting in 10 seconds..."))
		time.Sleep(time.Second * 10)
		os.Exit(3)
	}
	go func() {
		time.Sleep(time.Second * 15)
		go loader.AuthRoutine(userData.Key)
		time.Sleep(time.Second * 5)
		go loader.AuthRoutine(userData.Key)
		time.Sleep(time.Second * 5)
		go loader.AuthRoutine(userData.Key)
		time.Sleep(time.Second * 5)
	}()
	username := auth.DiscordTag
	userData.Username = username
	userData.DiscordId = auth.DiscordId
	userData.Version = Version
	
	fmt.Println(colors.Prefix() + colors.White("Welcome back, ") + colors.Red(username) + "!")
	loader.CheckForUpdate(userData)

	fmt.Println(colors.Prefix() + colors.Yellow("Loading your data..."))
	
	profiles := loader.LoadProfiles()
	proxies := loader.LoadProxies()
	token := loader.LoadMonitorToken()

	token = strings.ReplaceAll(token, "\"", "")

	fmt.Println(colors.Prefix() + colors.Green("Successfully loaded profiles, proxies and tokens!"))

	loader.UpdateRichPresence(Version)

	if token != "" {
		var wg sync.WaitGroup
		fmt.Println(colors.Prefix() + colors.Yellow("Logging into Discord with Monitoring token..."))
		wg.Add(1)
		go getpw.MonitorDiscord(&wg, token)
		wg.Wait()
	}

	go getpw.MonitorClipboard()
	
	var wg sync.WaitGroup
	wg.Add(1)
	go getpw.PWShareConnectReceive(&wg, userData)
	wg.Wait()

	menus.MainMenu(userData, profiles, proxies)
}
