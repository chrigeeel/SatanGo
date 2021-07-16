package getpw

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/chrigeeel/satango/colors"
)

type HyperInfo struct {
	ReleaseId string `json:"releaseId"`
	BpToken string `json:"bpToken"`
	CollectBilling bool `json:"collectBilling"`
	RequireLogin bool `json:"requireLogin"`
	AccountId string `json:"accountId"`
}
type PWStruct struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Site string `json:"site"`
	SiteType string `json:"siteType"`
	HyperInfo HyperInfo `json:"hyperInfo"`
	Mode string `json:"mode"`
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
	r2 := regexp.MustCompile(`\/purchase\/([^\?]*)`)
	for {
		p := <- PWC
		if p.Site == site || p.Mode == "clipboard" || p.Mode == "extension" {
			if p.Mode == "share" {
				fmt.Println(colors.Prefix() + colors.Green("You just received the password ") + colors.White("\"") + colors.Green(p.Password) + colors.White("\"") + colors.Green(" on the site ") + colors.White("\"") + colors.Green(p.Site) + colors.White("\""))
				fmt.Println(colors.Prefix() + colors.White("\"") + colors.Green(p.Username) + colors.White("\"") + colors.Green(" sent that to you, say thanks!"))
			}
			password = p.Password
			m := r.FindStringSubmatch(password)
			if len(m) == 2 {
				password = m[1]
			}
			m2 := r2.FindStringSubmatch(p.Password)
			if len(m2) == 2 {
				releaseId = m2[1]
			}
			password = strings.ReplaceAll(password, " ", "")
			password = strings.ReplaceAll(password, "\n", "")
			password = strings.ReplaceAll(password, "\r", "")
			p.Password = password
			if p.Mode != "extension" {
				p.HyperInfo.ReleaseId = releaseId
			}
			
			fmt.Println(colors.Prefix() + colors.White("Detected password \"") + colors.Red(password) + colors.White("\""))
			lookingForPw = false
			return p
		}
		if p.Mode == "discord" {
			content := p.Password
			m := r.FindStringSubmatch(content)
			if len(m) == 2 {
				password = m[1]
			} else {
				continue
			}
			r3 := regexp.MustCompile(`[^\/]*\.[^\/]*\.?[^\/]*`)
			siteB := string(r3.Find([]byte(site)))
			res4, err := regexp.MatchString(siteB, content)
			if err != nil || !res4 {
				continue
			}
			password = strings.ReplaceAll(password, " ", "")
			p.Password = password
			fmt.Println(colors.Prefix() + colors.White("Detected password \"") + colors.Red(password) + colors.White("\""))
			lookingForPw = false
			return p
		}
	}
}

func GetPwAfk() PWStruct {
	fmt.Println(colors.Prefix() + colors.Red("Waiting for Password from PW Sharing Server or Discord Monitoring"))
	fmt.Println(colors.Prefix() + colors.White("Copy the text \"") + colors.Red("exit") + colors.White("\" to exit"))
	r := regexp.MustCompile("(?i)password=\n*?\\s*?(\\S+)")
	var password string
	for {
		p := <- PWC
		if p.Mode == "share" && p.SiteType == "hyper" {
			fmt.Println(colors.Prefix() + colors.Green("You just received the password ") + colors.White("\"") + colors.Green(p.Password) + colors.White("\"") + colors.Green(" on the site ") + colors.White("\"") + colors.Green(p.Site) + colors.White("\""))
			fmt.Println(colors.Prefix() + colors.White("\"") + colors.Green(p.Username) + colors.White("\"") + colors.Green(" sent that to you, say thanks!"))
			return p
		}
		if p.Mode == "discord" {			
			content := p.Password
			m := r.FindStringSubmatch(content)
			if len(m) == 2 {
				password = m[1]
			} else {
				continue
			}
			r3 := regexp.MustCompile(`[^\/]*\.[^\/]*\.?[^\/]*`)
			siteB := r3.Find([]byte(content))
			if siteB == nil {
				continue
			}
			site := "https://" + string(siteB) + "/"
			password = strings.ReplaceAll(password, " ", "")
			password = strings.ReplaceAll(password, "Detected", "")
			p.Password = password
			p.Site = site
			fmt.Println(colors.Prefix() + colors.White("Detected password \"") + colors.Red(password) + colors.White("\" on the site \"") + colors.Red(site) + colors.White("\""))
			lookingForPw = false
			return p
		}
		if p.Password == "exit" {
			return p
		}
	}
}