package loader

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/chrigeeel/satango/colors"
	"github.com/mitchellh/go-ps"
	"github.com/shirou/gopsutil/process"
)

type user struct {
	Username string `json:"username"`
}

type authStruct struct {
	User user `json:"user"`
}

func debugChecker() {

	processes, _ := process.Pids()

	for i := 0; i < len(processes); i++ {
		pid := int(processes[i])
		p, err := ps.FindProcess(pid)
		if err != nil {
			fmt.Println(err)
		}
		if p != nil {
			name := p.Executable()
			name = strings.ToLower(name)
			pattern := "dnspy|httpdebuggersvc|fiddler|wireshark|charles|dragonfly|httpwatch|burpsuite|hxd|postman|http toolkit"
			match, _ := regexp.MatchString(pattern, name)
			if match == true {
				fmt.Println(colors.Prefix() + colors.White("Please stop any debuggers. Exiting..."))
				time.Sleep(3 * time.Second)
				os.Exit(3)
			}
		}

	}
}

func AuthKey(key string) authStruct {
	//debugChecker()
	fmt.Println(colors.Prefix() + colors.Yellow("Authenticating your key..."))
	resp, err := http.Get("http://50.16.47.99/auth/" + key)
	if err != nil {
		fmt.Println(colors.Prefix() + colors.Red("The Auth API is down atm! Please try again."))
	}
	defer resp.Body.Close()
	respBytes, _ := ioutil.ReadAll(resp.Body)
	var auth authStruct
	json.Unmarshal(respBytes, &auth)
	if auth.User.Username == "" {
		fmt.Println(colors.Prefix() + colors.Red("Authentication failed. Your key is invalid!"))
		time.Sleep(time.Second * 3)
		os.Exit(3)
	}
	fmt.Println(colors.Prefix() + colors.Green("Successfully authenticated your key!"))
	return auth
}
