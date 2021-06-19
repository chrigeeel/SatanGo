package getpw

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/chrigeeel/satango/colors"
)

type PWStruct struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Site string `json:"site"`
	ReleaseId string `json:"releaseId"`
}

var lookingForPw bool
var PWC chan PWStruct = make(chan PWStruct)

func GetPw2(site string) PWStruct {
	fmt.Println(colors.Prefix() + colors.Red("Waiting for Password"))
	fmt.Println(colors.Prefix() + colors.White("Copy the text \"") + colors.Red("exit") + colors.White("\" to exit"))
	var password string
	var releaseId string
	lookingForPw = true
	r := regexp.MustCompile("(?i)password=\n*?\\s*?(\\S+)")
	r2 := regexp.MustCompile("\\/purchase\\/([^\\?]*)")
	for {
		p := <-PWC
		if p.Site == site || p.Site == "clipboard" || p.Site == "link" {
			password = p.Password
			password = strings.ReplaceAll(password, " ", "")
			password = strings.ReplaceAll(password, "\n", "")
			password = strings.ReplaceAll(password, "\r", "")
			m := r.FindStringSubmatch(password)
			if len(m) == 2 {
				password = m[1]
			}
			m2 := r2.FindStringSubmatch(password)
			if len(m2) == 2 {
				releaseId = m2[1]
			}
			p.Password = password
			p.ReleaseId = releaseId
			fmt.Println(colors.Prefix() + colors.White("Detected password \"") + colors.Red(password) + colors.White("\""))
			lookingForPw = false
			return p
		}
	}
}