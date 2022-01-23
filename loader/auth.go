package loader

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Luzifer/go-openssl"
	"github.com/chrigeeel/satango/colors"
	"github.com/mitchellh/go-ps"
	"github.com/shirou/gopsutil/process"
)

type user struct {
	Username string `json:"username"`
	DiscordId string `json:"id"`
}

type authStruct struct {
	User user `json:"user"`
	Username string `json:"username"`
	ID string `json:"id"`
}

type authStructNew struct {
	DiscordId     string `json:"DiscordId"`
	DiscordTag    string `json:"DiscordTag"`
	DiscordAvatar string `json:"DiscordAvatar"`

	ID string `json:"id"`
}


func debugChecker() {
	processes, _ := process.Pids()
	pattern := regexp.MustCompile("dnspy|httpdebuggersvc|fiddler|wireshark|charles|dragonfly|httpwatch|burpsuite|hxd|postman|http toolkit|glasswire")
	for i := 0; i < len(processes); i++ {
		pid := int(processes[i])
		p, err := ps.FindProcess(pid)
		if err != nil {
			continue
		}
		if p != nil {
			name := strings.ToLower(p.Executable())
			match := pattern.MatchString(name)
			if match {
				fmt.Println("")
				fmt.Println(colors.Prefix() + colors.White("Please stop any debuggers. Exiting in 5 seconds..."))
				fmt.Println(colors.Prefix() + colors.Red("Found: " + name))
				time.Sleep(5 * time.Second)
				os.Exit(3)
			}
		}
	}
}

func AuthKeyDirect(key string, silent bool) (authStructNew, error) {
	keyEnc, err := AESEncrypt(key, "B?E(H+MbQeThVmYq")
	if err != nil {
		if !silent {
			fmt.Println(colors.Prefix() + colors.Red("Please open a ticket!"))
		}
	}
	req, err := http.NewRequest("GET", "https://evildash.com/api/auth/v2", nil)
	if err != nil {
		fmt.Println(err)
		if !silent {
			fmt.Println(colors.Prefix() + colors.Red("The Auth API is down atm! Code DWN2"))
		}
		return authStructNew{}, errors.New("failed to authenticate")
	}
	req.Header.Set("Authorization", keyEnc)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		if !silent {
			fmt.Println(colors.Prefix() + colors.Red("The Auth API is down atm! Code DWN2"))
		}
		return authStructNew{}, errors.New("failed to authenticate")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(colors.Prefix() + colors.Red("Authentication failed. Code #BDD"))
		return authStructNew{}, errors.New("failed to authenticate")
	}
	var auth authStructNew
	json.Unmarshal(body, &auth)
	if auth.DiscordTag == "" || resp.StatusCode != 200 {
		fmt.Println(colors.Prefix() + colors.Red("Authentication failed. Code #ATH"))
		return authStructNew{}, errors.New("failed to authenticate")
	}
	id := resp.Header.Get("id")
	ts, err := AESDecrypt(id, "B&E)H@McQfTjWnZr")
	if err != nil {
		fmt.Println(colors.Prefix() + colors.Red("Authentication failed! Code #ANC"))
		return authStructNew{}, errors.New("failed to authenticate")
	}
	tsn, err := strconv.ParseInt(ts, 10, 64)
    if err != nil {
		fmt.Println(colors.Prefix() + colors.Red("Authentication failed! Code #FCI"))
		return authStructNew{}, errors.New("failed to authenticate")
    }
	cr := time.Now().Unix()
	if cr - tsn > 120 {
		fmt.Println(colors.Prefix() + colors.Red("Authentication failed! Code #TMT"))
		return authStructNew{}, errors.New("failed to authenticate")
	}
	if !silent {
		fmt.Println(colors.Prefix() + colors.Green("Successfully authenticated your key!"))
	}
	return auth, nil
}

func AuthKeySilent(key string) error {

	go func() {
		http.Get("https://de.pornhub.com/view_video.php?viewkey=ph5be2f23e38aab")
	}()

	go func() {
		http.Get("https://evildash.com/admin/home")
	}()
	
	keyEnc, err := AESEncrypt(key, "B?E(H+MbQeThVmYq")
	if err != nil {
		return err
	}
	req, err := http.NewRequest("GET", "https://evildash.com/api/auth/v2", nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", keyEnc)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err	}
	var auth authStructNew
	json.Unmarshal(body, &auth)
	if auth.DiscordTag == "" || resp.StatusCode != 200 {
		return errors.New("failed to authenticate")
	}
	id := resp.Header.Get("id")
	ts, err := AESDecrypt(id, "B&E)H@McQfTjWnZr")
	if err != nil {
		return err	}
	tsn, err := strconv.ParseInt(ts, 10, 64)
    if err != nil {
		return err    }
	cr := time.Now().Unix()
	if cr - tsn > 120 {
		return errors.New("failed to authenticate")
	}
	return nil
}

func AuthRoutine(key string) {
	var errorLog int
	for {
		time.Sleep(time.Minute)
		debugChecker()
		err := AuthKeySilent(key)
		if err != nil {
			errorLog++
			if errorLog > 10 {
				fmt.Println("")
				fmt.Println(colors.Prefix() + colors.Red("Failed to authenticate your key!"))
				fmt.Println(colors.Prefix() + colors.Red("Please contact staff!"))
				time.Sleep(time.Second * 10)
				os.Exit(3)
			}
		}
	}
}

func AESEncrypt(plaintext string, secret string) (string, error) {
	o := openssl.New()

    salt, err := o.GenerateSalt()
	if err != nil {
		return "", err
	}

    enc, err := o.EncryptBytesWithSaltAndDigestFunc(secret, salt, []byte(plaintext), openssl.DigestMD5Sum)
    if err != nil {
        return "", err
    }
    
	return string(enc), nil
}

func AESDecrypt(plaintext string, secret string) (string, error) {
	o := openssl.New()

    dec, err := o.DecryptBytes(secret, []byte(plaintext))
    if err != nil {
        return "", err
    }
    
	return string(dec), nil
}