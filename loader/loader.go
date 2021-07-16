package loader

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/chrigeeel/satango/colors"
)


type UserDataStruct struct {
	Key     string `json:"key"`
	Webhook string `json:"webhook"`
	Username string `json:"username"`
	DiscordId string `json:"discordId"`
	Version string `json:"version"`
}

func CreateSettings() {
	scanner := bufio.NewScanner(os.Stdin)
	err := os.Mkdir("./settings", 0755)
	if err != nil {
		fmt.Println(colors.Prefix() + colors.White("Setup failed. Please contact staff or try again."))
	}
	fmt.Println(colors.Prefix() + colors.White("Please Input your key"))
	fmt.Print(colors.Prefix() + "> ")
	scanner.Scan()
	keyData := scanner.Text()
	fmt.Println(colors.Prefix() + colors.White("Please Input your Webhook URL"))
	fmt.Print(colors.Prefix() + "> ")
	scanner.Scan()
	webhookData := scanner.Text()
	data := UserDataStruct{
		Key:     keyData,
		Webhook: webhookData,
	}
	jsonData, _ := json.MarshalIndent(data, "", "    ")
	_ = ioutil.WriteFile("./settings/settings.json", jsonData, 0644)
}

func LoadSettings() UserDataStruct {
	if _, err := os.Stat("./settings"); os.IsNotExist(err) {
		fmt.Println(colors.Prefix() + colors.White("No settings folder detected. Setting up..."))
		CreateSettings()
	}
	settingsFile, err := os.Open("./settings/settings.json")
	if err != nil {
		fmt.Println(colors.Prefix() + colors.White("Setup failed. Please contact staff or try again."))
	}
	defer settingsFile.Close()
	settingsBytes, _ := ioutil.ReadAll(settingsFile)
	var settings UserDataStruct
	json.Unmarshal(settingsBytes, &settings)
	if settings.Key == "" {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println(colors.Prefix() + colors.White("Please Input your key"))
		fmt.Print(colors.Prefix() + "> ")
		scanner.Scan()
		keyData := scanner.Text()
		data := UserDataStruct{
			Key: keyData,
		}
		jsonData, _ := json.MarshalIndent(data, "", "    ")
		_ = ioutil.WriteFile("./settings/settings.json", jsonData, 0644)
		settings.Key = keyData
	}
	return settings
}

func LoadProxies() []string {
	proxiesFile, err := os.Open("./settings/proxies.txt")
	if err != nil {
		data := []byte("localhost")
		err := ioutil.WriteFile("./settings/proxies.txt", data, 0644)
		if err != nil {
			fmt.Println(colors.Prefix() + colors.Red("You don't have any valid proxies in /settings/proxies.txt! Please create the file."))
		}
	}
	defer proxiesFile.Close()

	scanner := bufio.NewScanner(proxiesFile)
	var proxies []string
	for scanner.Scan() {
		proxies = append(proxies, scanner.Text())
	}

	return proxies
}

func LoadMonitorToken() string {
	monitorFile, err := os.Open("./settings/monitortoken.txt")
	if err != nil {
		data := askFor("Please input your Discord Monitoring token (Should be safe to use main account, press [ENTER] to skip)")
		err := ioutil.WriteFile("./settings/monitortoken.txt", []byte(data), 0644)
		if err != nil {
			fmt.Println(colors.Prefix() + colors.Red("Your Discord Monitoring token is not set up in ./settings/monitortoken.txt ! Please create the file"))
		}
		monitorFile, _ = os.Open("./settings/monitortoken.txt")
	}
	defer monitorFile.Close()

	token, err := ioutil.ReadAll(monitorFile)
	if err != nil {
		return ""
	}
	return string(token)
}