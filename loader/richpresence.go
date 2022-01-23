package loader

import (
	"time"

	"github.com/hugolgst/rich-go/client"
)


func UpdateRichPresence(version string) {
	err := client.Login("794918479881961523")
	if err != nil {
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
		return
	}
}