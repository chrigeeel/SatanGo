package getpw

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/chrigeeel/satango/colors"
	"github.com/chrigeeel/satango/loader"
	"github.com/chrigeeel/satango/utility"
	"github.com/gorilla/websocket"
)

var LastPass string
var Share bool = false

var addr = "52.45.120.253:8080"

var PWShare chan PWStruct

func PWShareConnectReceive(wg *sync.WaitGroup, userData loader.UserDataStruct) {

	for i := 0; i != -1; i++ {
		username := userData.Username

		u := url.URL{Scheme: "ws", Host: addr, Path: "/ws"}
		fmt.Println(colors.Prefix() + colors.Yellow("Connecting to PWSharing Server..."))
	
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			if i == 0 {
				wg.Done()
			}
			fmt.Println(colors.Prefix() + colors.Red("Failed connecting to PWSharing Server!"))
			return
		}
		defer c.Close()
		fmt.Println(colors.Prefix() + colors.Green("Successfully connected to PWSharing Server!"))
		if i == 0 {
			wg.Done()
		}
		done := make(chan struct{})
		PWShare = make(chan PWStruct)
	
		go func() {
			defer close(done)
			for {
				_, message, err := c.ReadMessage()
				if err != nil {
					fmt.Println(colors.Prefix() + colors.Red("Failed to read message!"))
					return
				}
				var m PWStruct
				err = json.Unmarshal([]byte(message), &m)
				if m.Username != username && err == nil {
					PWC <- m
				}
			}
		}()
	
		ticker := time.NewTicker(time.Second * 10)
		defer ticker.Stop()

		reconnect := false
	
		go func() {
			for {
				select {
				case m := <-PWShare:
					fmt.Println(colors.Prefix() + colors.Yellow("Sending password to PWSharing Server..."))
					jsonContent, err := json.Marshal(m)
					if err != nil {
						fmt.Println(colors.Prefix() + colors.Red("Invalid Password!"))
					}
					err = c.WriteMessage(websocket.TextMessage, jsonContent)
					if err != nil {
						fmt.Println(colors.Prefix() + colors.Red("Failed to Send password to PWSharing Server!"))
						return
					}
					fmt.Println(colors.Prefix() + colors.Green("Successfully Sent password to PWSharing Server!"))
				case <-done:
					return
				case <-ticker.C:
					err := c.WriteMessage(websocket.TextMessage, []byte("ping"))
					if err != nil {
						fmt.Println("")
						fmt.Println(colors.Prefix() + colors.Red("Failed to Ping PWSharing Server, trying again..."))
						err = c.WriteMessage(websocket.TextMessage, []byte("ping"))
						if err != nil {
							fmt.Println(colors.Prefix() + colors.Red("Failed to Ping PWSharing Server again!"))
							c.Close()
							close(done)
							fmt.Println(colors.Prefix() + colors.Yellow("Reconnecting to PWSharing Server..."))
							reconnect = true
							return
						}

					}
				}
			}
		}()
		for !reconnect {
			time.Sleep(time.Second * 1)
		}
	}
}

func PWSharingSend(userData loader.UserDataStruct, password string, site string, siteType string) {
	if password != LastPass && Share == true {
		LastPass = password
		m := PWStruct{
			Username: userData.Username,
			Password: password,
			Site: site,
			SiteType: siteType,
			Mode: "share",
		}
		PWShare <- m
	}
}

func PWSharingSend2(p PWStruct) {
	if p.Password != LastPass && Share == true {
		LastPass = p.Password
		p.Mode = "share"
		PWShare <- p
	}
}

func AskForPwShare() {
	fmt.Println(colors.Prefix() + colors.Red("(Y/N) Would you like to turn on PWSharing? (You will only receive passwords from other people if turned on, highly recommended!)"))
	ans := utility.AskForSilent()
	if strings.ToLower(ans)[0:1] == "y" {
		Share = true
		fmt.Println(colors.Prefix() + colors.White("Turned PWSharing on!"))
	} else {
		Share = false
		fmt.Println(colors.Prefix() + colors.White("Turned PWSharing off!"))
	}
}