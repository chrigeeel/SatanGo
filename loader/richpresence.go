package loader

import (
	"fmt"
	"time"

	"github.com/chrigeeel/satango/colors"
	"github.com/hugolgst/rich-go/client"
)


func UpdateRichPresence(version string) {
	fmt.Println(colors.Prefix() + colors.Yellow("Updating your Discord Rich Presence..."))
	err := client.Login("794918479881961523")
	if err != nil {
		fmt.Println(colors.Prefix() + colors.Red("Failed to update your Discord Rich Presence!"))
		return
	}

	startTime := time.Now()

	err = client.SetActivity(client.Activity{
		State:      "Sacrificing Stock",
		Details:    "Version " + version,
		LargeImage: "satanlogo",
		LargeText:  "@SatanBots",
		SmallImage: "face",
		SmallText:  "SatanBot",
		Timestamps: &client.Timestamps{
			Start: &startTime,
		},
		Buttons: []*client.Button{
			{
				Label: "Twitter",
				Url: "https://twitter.com/SatanBots",
			},
		},
	})
	
	if err != nil {
		fmt.Println(colors.Prefix() + colors.Red("Failed to update your Discord Rich Presence!"))
		return
	}
	fmt.Println(colors.Prefix() + colors.Green("Successfully updated your Discord Rich Presence!"))
}