package loader

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cavaliercoder/grab"
	"github.com/chrigeeel/satango/colors"
)


type versionStruct struct {
	Version string `json:"version"`
}

func CheckForUpdate(userData UserDataStruct) {
	fmt.Println(colors.Prefix() + colors.Yellow("Checking for Updates..."))
	key := userData.Key
	resp1, err := http.Get("http://50.16.47.99/update/info")
	if err != nil {
		fmt.Println(colors.Prefix() + colors.Red("Failed to check for Updates! Exiting..."))
		time.Sleep(time.Second * 5)
		os.Exit(3)
		return
	}
	defer resp1.Body.Close()
	body, _ := ioutil.ReadAll(resp1.Body)
	var version versionStruct
	json.Unmarshal(body, &version)

	newVersion := version.Version
	newVersionInt, err := strconv.Atoi(strings.ReplaceAll(newVersion, ".", ""))
	if err != nil {
		fmt.Println(colors.Prefix() + colors.Red("Failed to check for Updates! Exiting..."))
		time.Sleep(time.Second * 5)
		os.Exit(3)
		return
	}

	oldVersion := userData.Version
	oldVersionInt, err := strconv.Atoi(strings.ReplaceAll(oldVersion, ".", ""))
	if err != nil {
		fmt.Println(colors.Prefix() + colors.Red("Failed to check for Updates! Exiting..."))
		time.Sleep(time.Second * 5)
		os.Exit(3)
		return
	}

	if oldVersionInt >= newVersionInt {
		fmt.Println(colors.Prefix() + colors.Green("You are on the latest Version!"))
		return
	}
	fmt.Println(colors.Prefix() + colors.Green("Found a new Version: " + newVersion))
	fmt.Println(colors.Prefix() + colors.Yellow("Starting Download..."))

	client := grab.NewClient()
	req, _ := grab.NewRequest("satango - " + newVersion + ".exe", "http://50.16.47.99/update/" + key)

	resp := client.Do(req)

	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()

Loop:
	for {
		select {
		case <-t.C:
			fmt.Printf("\r" + colors.Prefix() + colors.Yellow("Downloaded %v / %v bytes (%.2f%%)"),
				resp.BytesComplete(),
				resp.Size,
				100*resp.Progress())

		case <-resp.Done:
			// download is complete
			break Loop
		}
	}

	// check for errors
	if err := resp.Err(); err != nil {
		fmt.Println(colors.Prefix() + colors.Red("Download failed! Please try again or contact Staff"))
		time.Sleep(time.Second * 5)
		os.Exit(1)
		return
	}
	fmt.Println("\n")
	fmt.Println(colors.Prefix() + colors.Green("Successfully downloaded new Version. Please start ") + colors.White("satango - " + newVersion + ".exe!"))
	fmt.Println(colors.Prefix() + colors.Red("Please delete the old satango! Exiting..."))
	time.Sleep(time.Second * 5)
	os.Exit(3)
}