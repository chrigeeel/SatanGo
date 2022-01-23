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
	Version string = "0.5.93"
)

/*

	hello, although i wrote this code some time ago and this was my
	FIRST EVER go project, i hope you can still learn something from it
	and maybe use it for any of your go projects :)
	some of the modules might still work but most are probably broken...
	all apis are still online, even the TL captcha api!
	feel free to use it, it only works on TL captchas but is pretty decent.
	it's only running on a 2-core server so don't expect too much though!

	i will be adding more comments soon if there's enough demand



	cooldiscord is https://github.com/bwmarrin/discordgo but edited so WS-connections work with non-bot tokens.
	feel free to use it seperately
	

*/

func main() {
	go getpw.MonitorExtension()

	time.Sleep(time.Millisecond * 10)

	menus.CallClear()

	fmt.Println(colors.Prefix() + colors.White("Starting... | Version " + Version))

	userData := loader.LoadSettings()
	auth, err := loader.AuthKeyDirect(userData.Key, false)
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
	
	loader.LoadProfiles()
	loader.LoadProxies()
	token := loader.LoadMonitorToken()

	token = strings.ReplaceAll(token, "\"", "")

	fmt.Println(colors.Prefix() + colors.Green("Successfully loaded profiles, proxies and tokens!"))

	go loader.UpdateRichPresence(Version)

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

	menus.MainMenu(userData)
}
