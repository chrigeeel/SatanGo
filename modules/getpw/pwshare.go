package getpw

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/chrigeeel/satango/colors"
	"github.com/chrigeeel/satango/loader"
	"github.com/gorilla/websocket"
)

var LastPass string

var addr = "44.193.110.231:8080"

var PWShare chan PWStruct

func PWShareConnectReceive(userData loader.UserDataStruct) {

	username := userData.Username

	u := url.URL{Scheme: "ws", Host: addr, Path: "/ws"}
	fmt.Println(colors.Prefix() + colors.Yellow("Connecting to PWSharing Server..."))

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println(colors.Prefix() + colors.Red("Failed connecting to PWSharing Server!"))
		return
	}
	defer c.Close()
	fmt.Println(colors.Prefix() + colors.Green("Successfully connected to PWSharing Server!"))

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
				fmt.Println(colors.Prefix() + colors.Green("You just received the password ") + colors.White("\"") + colors.Green(m.Password) + colors.White("\"") + colors.Green(" on the site ") + colors.White("\"") + colors.Green(m.Site) + colors.White("\""))
				fmt.Println(colors.Prefix() + colors.White("\"") + colors.Green(m.Username) + colors.White("\"") + colors.Green(" sent that to you, say thanks!"))
				PWC <- m
			}
		}
	}()

	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

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
					fmt.Println(colors.Prefix() + colors.Red("Failed to Ping PWSharing Server!"))
					return
				}
			}
		}
	}()
	for {
		time.Sleep(time.Second * 1)
	}
}

func PWSharingSend(userData loader.UserDataStruct, password string, site string) {
	if password != LastPass {
		LastPass = password
		m := PWStruct{
			Username: userData.Username,
			Password: password,
			Site: site,
		}
		PWShare <- m
	}
}