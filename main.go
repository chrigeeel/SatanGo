package main

import (
	"fmt"

	"github.com/chrigeeel/satango/colors"
	"github.com/chrigeeel/satango/loader"
	"github.com/chrigeeel/satango/menus"
)

func main() {

	menus.CallClear()

	fmt.Println(colors.Prefix() + colors.White("Starting..."))

	userData := loader.LoadSettings()
	auth := loader.AuthKey(userData.Key)
	username := auth.User.Username
	userData.Username = username

	fmt.Println(colors.Prefix() + colors.White("Welcome back, ") + colors.Red(username) + "!")
	loader.LoadProfiles()
	fmt.Println(colors.Prefix() + colors.Yellow("Loading your data..."))
	
	profiles := loader.LoadProfiles()
	proxies := loader.LoadProxies()

	fmt.Println(colors.Prefix() + colors.Green("Successfully loaded profiles, proxies and tokens"))
	menus.MainMenu(userData, profiles, proxies)

	/*
		//client := modules.CoolClient("estr2.premium.maskedproxy.xyz:5806:6xbEh5Wn:ZCgP3OCKzxBF0U1WhwpOBpaKFKecXSX86yOQKa72C4IyWNCe9lmJMgHDXvUhn6peDPji7-THg5jSkBjt")
		var tokens []string
		tokens = append(tokens, "mfa.Ihmo-BKnwS8tDKx3zKNG0HY5__ZjAyQzwFszPoIxhYVCr6cnuEW3R73NQioQg7JCUKJY-a5i3PhQ4K4Kshxn")
		tokens = append(tokens, "asdasd")
		modules.TlLogin(tokens)
	*/
	/*
		req2, err := http.NewRequest("GET", "https://button-backend.tldash.ai/api/register/demo/demo", nil)
		if err != nil {
			log.Fatal(err)
		}
		req2.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.93 Safari/537.36")
		resp2, err := client.Do(req2)
		if err != nil {
			log.Fatal(err)
		}
		defer resp2.Body.Close()

		b2, _ := ioutil.ReadAll(resp2.Body)
		//fmt.Println(b)
		bs2 := string(b2[:])
		fmt.Println(bs2)
	*/
}
