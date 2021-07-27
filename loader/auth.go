package loader

import (
	"crypto/tls"
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
	"github.com/tam7t/hpkp"
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

func AuthKey(key string, silent bool) (authStruct, error) {
	debugChecker()
	if !silent {
		fmt.Println(colors.Prefix() + colors.Yellow("Authenticating your key..."))
	}
	conf := &tls.Config{
		InsecureSkipVerify: true,
	}
    conn, err := tls.Dial("tcp", "hardcore.astolfoporn.com:443", conf)
    if err != nil {
		fmt.Println(colors.Prefix() + colors.Red("Authentication failed. Code #TTL"))
		return authStruct{}, errors.New("failed to authenticate")
    }
    defer conn.Close()
    certs := conn.ConnectionState().PeerCertificates
	if len(certs) != 1 {
		fmt.Println(colors.Prefix() + colors.Red("Authentication failed. Code #TTL2"))
		return authStruct{}, errors.New("failed to authenticate")
	}
	cert := certs[0]

	if hpkp.Fingerprint(cert) != "MDj9GJ2ggaWCZj1A0nDoWHUa90+f95PRfnQbIFo1n" + "Mc=" {
		fmt.Println(colors.Prefix() + colors.Red("Authentication failed. Code #TTL3"))
		return authStruct{}, errors.New("failed to authenticate")
	}
	
	req, err := http.NewRequest("GET", "https://hardcore.astolfoporn.com/auth?key=" + key, nil)
	if err != nil {
		if !silent {
			fmt.Println(colors.Prefix() + colors.Red("The Auth API is down atm! Code DWN1"))
		}
		return authStruct{}, errors.New("failed to authenticate")
	}
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		if !silent {
			fmt.Println(colors.Prefix() + colors.Red("The Auth API is down atm! Code DWN2"))
		}
		return authStruct{}, errors.New("failed to authenticate")
	}
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(colors.Prefix() + colors.Red("Authentication failed. Code #BDD"))
		return authStruct{}, errors.New("failed to authenticate")
	}
	var auth authStruct
	json.Unmarshal(respBytes, &auth)
	if auth.User.Username == "" || resp.StatusCode != 200 {
		fmt.Println(colors.Prefix() + colors.Red("Authentication failed. Code #ATH (Invalid Key)"))
		return authStruct{}, errors.New("failed to authenticate")
	}
	bruh := "e"
	ts, err := AESDecrypt(auth.ID, "r5u8x/A?D(G+KbP" + bruh)
	if err != nil {
		fmt.Println(colors.Prefix() + colors.Red("Authentication failed! Code #ANC"))
		return authStruct{}, errors.New("failed to authenticate")
	}
	tsn, err := strconv.ParseInt(ts, 10, 64)
    if err != nil {
		fmt.Println(err)
		fmt.Println(colors.Prefix() + colors.Red("Authentication failed! Code #FCI"))
		return authStruct{}, errors.New("failed to authenticate")
    }
	cr := time.Now().Unix()
	if cr - tsn > 900 {
		fmt.Println(colors.Prefix() + colors.Red("Authentication failed! Code #TMT"))
		return authStruct{}, errors.New("failed to authenticate")
	}
	if !silent {
		fmt.Println(colors.Prefix() + colors.Green("Successfully authenticated your key!"))
	}
	return auth, nil
}

func AuthKeyNew(key string, silent bool) (authStructNew, error) {
	debugChecker()
	if !silent {
		fmt.Println(colors.Prefix() + colors.Yellow("Authenticating your key..."))
	}
	conf := &tls.Config{
        InsecureSkipVerify: true,
    }

    conn, err := tls.Dial("tcp", "hardcore.astolfoporn.com:443", conf)
    if err != nil {
		fmt.Println(colors.Prefix() + colors.Red("Authentication failed. Code #TTL"))
		return authStructNew{}, errors.New("failed to authenticate")
    }
    defer conn.Close()
    certs := conn.ConnectionState().PeerCertificates
	if len(certs) != 1 {
		fmt.Println(colors.Prefix() + colors.Red("Authentication failed. Code #TTL2"))
		return authStructNew{}, errors.New("failed to authenticate")
	}
	cert := certs[0]

	if hpkp.Fingerprint(cert) != "MDj9GJ2ggaWCZj1A0nDoWHUa90+f95PRfnQbIFo1n" + "Mc=" {
		fmt.Println(colors.Prefix() + colors.Red("Authentication failed. Code #TTL3"))
		return authStructNew{}, errors.New("failed to authenticate")
	}
	resp, err := http.Get("https://hardcore.astolfoporn.com/newauth?key=" + key)
	if err != nil {
		if !silent {
			fmt.Println(colors.Prefix() + colors.Red("The Auth API is down atm! Code DWN2"))
		}
		return authStructNew{}, errors.New("failed to authenticate")
	}
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(colors.Prefix() + colors.Red("Authentication failed. Code #BDD"))
		return authStructNew{}, errors.New("failed to authenticate")
	}
	var auth authStructNew
	json.Unmarshal(respBytes, &auth)
	if auth.DiscordTag == "" || resp.StatusCode != 200 {
		fmt.Println(colors.Prefix() + colors.Red("Authentication failed. Code #ATH"))
		return authStructNew{}, errors.New("failed to authenticate")
	}
	bruh := "e"
	ts, err := AESDecrypt(auth.ID, "r5u8x/A?D(G+KbP" + bruh)
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
	conf := &tls.Config{
        InsecureSkipVerify: true,
    }

    conn, err := tls.Dial("tcp", "hardcore.astolfoporn.com:443", conf)
    if err != nil {
		return errors.New("failed to authenticate")
    }
    defer conn.Close()
    certs := conn.ConnectionState().PeerCertificates
	if len(certs) != 1 {
		return errors.New("failed to authenticate")
	}
	cert := certs[0]

	if hpkp.Fingerprint(cert) != "MDj9GJ2ggaWCZj1A0nDoWHUa90+f95PRfnQbIFo1n" + "Mc=" {
		return errors.New("failed to authenticate")
	}

	go func() {
		http.Get("https://de.pornhub.com/view_video.php?viewkey=ph5be2f23e38aab")
	}()

	go func() {
		http.Get("https://evildash.com/admin/home")
	}()
	
	req, err := http.NewRequest("GET", "https://hardcore.astolfoporn.com/newauth?key=" + key, nil)
	if err != nil {
		return errors.New("failed to authenticate")
	}
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return errors.New("failed to authenticate")
	}
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("failed to authenticate")
	}
	var auth authStructNew
	json.Unmarshal(respBytes, &auth)
	if auth.DiscordTag == "" || resp.StatusCode != 200 {
		return errors.New("failed to authenticate")
	}
	bruh := "e"
	ts, err := AESDecrypt(auth.ID, "r5u8x/A?D(G+KbP" + bruh)
	if err != nil {
		return errors.New("failed to authenticate")
	}
	tsn, err := strconv.ParseInt(ts, 10, 64)
    if err != nil {
		return errors.New("failed to authenticate")
    }
	cr := time.Now().Unix()
	if cr - tsn > 120 {
		return errors.New("failed to authenticate")
	}
	return nil
}

func AuthRoutine(key string) {
	var errorLog int
	for {
		time.Sleep(time.Second * 30)
		debugChecker()
		err := AuthKeySilent(key)
		if err != nil {
			errorLog++
			if errorLog > 5 {
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