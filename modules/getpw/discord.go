package getpw

import (
	"fmt"
	"strings"
	"sync"

	"github.com/chrigeeel/satango/colors"
	"github.com/chrigeeel/satango/cooldiscord"
)

var (
	userId string
)

func MonitorDiscord(wg *sync.WaitGroup, token string) {
	dg, err := cooldiscord.New(token)
	if err != nil {
		fmt.Println(colors.Prefix() + colors.Red("Failed to login with Discord Monitoring token!"))
		return
	}

	u, err := dg.User("@me")
	if err != nil {
		fmt.Println(colors.Prefix() + colors.Red("Failed to login with Discord Monitoring token!"))
		return
	}
	userId = u.ID

	dg.AddHandler(messageHandler)

	err = dg.Open()
	if err != nil {
		fmt.Println(colors.Prefix() + colors.Red("Failed to login with Discord Monitoring token!"))
		return
	}

	fmt.Println(colors.Prefix() + colors.Green("Successfully logged in as ") + colors.White("\"") + colors.Green(u.Username + "#" + u.Discriminator) + colors.White("\"!"))
	wg.Done()
}

func messageHandler(s *cooldiscord.Session, m *cooldiscord.MessageCreate) {
	var def *cooldiscord.MessageCreate

	if (m == def) {
		return
	}

	var content string
	content = m.Content
	if len(m.Embeds) > 0 {
		for _, embed := range m.Embeds {
			content = content + embed.Description
			if len(embed.Fields) > 0 {
				for _, field := range embed.Fields {
					content = content + field.Name
					content = content + field.Value
				}
			}
		}
	}
	content = strings.ReplaceAll(content, "\n", "")
	content = strings.ReplaceAll(content, "\r", "")
	content = strings.ReplaceAll(content, "Join", "")
	content = strings.ReplaceAll(content, "Checkout", "")
	p := PWStruct{
		Password: content,
		Mode: "discord",
	}
	PWC <- p
}