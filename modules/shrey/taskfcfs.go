package shrey

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/chrigeeel/satango/colors"
	"github.com/chrigeeel/satango/loader"
	"github.com/chrigeeel/satango/modules/getpw"
	"github.com/chrigeeel/satango/utility"
)

func taskfcfs(wg *sync.WaitGroup, userData loader.UserDataStruct, id int, p getpw.PWStruct, profile loader.ProfileStruct, site string) {
	defer wg.Done()

	beginTime := time.Now()

	payload := strings.NewReader("password=" + p.Password + "&authenticity_token=" + url.QueryEscape(profile.ShreyCSRF))

	rurl := site + "password"

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, err := http.NewRequest("POST", rurl, payload)
  
	if err != nil {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to load release!"))
		return
	}
	req.Header.Add("sec-ch-ua", "\" Not;A Brand\";v=\"99\", \"Google Chrome\";v=\"91\", \"Chromium\";v=\"91\"")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("Upgrade-Insecure-Requests", "1")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Add("Sec-Fetch-Site", "cross-site")
	req.Header.Add("Sec-Fetch-Mode", "navigate")
	req.Header.Add("Sec-Fetch-User", "?1")
	req.Header.Add("Sec-Fetch-Dest", "document")
	req.Header.Add("Cookie", "_shreyauth_session=" + profile.DiscordSession)
  
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to load release!"))
		return
	}
	defer resp.Body.Close()
	var shreyPassCookie string
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "pass" {
			shreyPassCookie = cookie.Value
		}
		if cookie.Name == "_shreyauth_session" {
			profile.DiscordSession = cookie.Value
		}
	}
	if shreyPassCookie == "" {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Wrong password or release OOS!"))
		return
	}

	rurl = site
	req, err = http.NewRequest("GET", rurl, nil)
	if err != nil {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to load release!"))
		return
	}
	req.Header.Add("Upgrade-Insecure-Requests", "1")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Add("Sec-Fetch-Site", "cross-site")
	req.Header.Add("Sec-Fetch-Mode", "navigate")
	req.Header.Add("Sec-Fetch-User", "?1")
	req.Header.Add("Sec-Fetch-Dest", "document")
	req.Header.Add("sec-ch-ua", "\" Not;A Brand\";v=\"99\", \"Google Chrome\";v=\"91\", \"Chromium\";v=\"91\"")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("Cookie", "_shreyauth_session="+profile.DiscordSession)
	resp, err = client.Do(req)
	if err != nil {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to load release!"))
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to load release!"))
		return
	}
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "_shreyauth_session" {
			profile.DiscordSession = cookie.Value
		}
	}
	r := regexp.MustCompile(`release_id=(\d*)`)
	match := r.FindStringSubmatch(string(body))
	if len(match) != 2 {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Wrong password or release OOS!"))
		return
	}
	releaseId := match[1]
	r = regexp.MustCompile(`csrf-token" content="([^"]*)`)
	match = r.FindStringSubmatch(string(body))
	if len(match) != 2 {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Wrong password or release OOS!"))
		return
	}
	profile.ShreyCSRF = match[1]

	fmt.Println(colors.TaskPrefix(id) + colors.Green("Successfully loaded Release " + releaseId + "!"))
	fmt.Println(colors.TaskPrefix(id) + colors.Yellow("Submitting Checkout..."))

	rurl = site + "subscribe"

	payload = strings.NewReader("release_id=" + releaseId + "&coupon=&stripeToken=&guess=&authenticity_token=" + url.QueryEscape(profile.ShreyCSRF))

	req, err = http.NewRequest("POST", rurl, payload)
	if err != nil {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to submit checkout!"))
		return
	}
	req.Header.Add("sec-ch-ua", "\" Not;A Brand\";v=\"99\", \"Google Chrome\";v=\"91\", \"Chromium\";v=\"91\"")
	req.Header.Add("Accept", "text/javascript, application/javascript, application/ecmascript, application/x-ecmascript, */*; q=0.01")
	req.Header.Add("X-CSRF-Token", "Qz6LiMeltxXOVh6jhoU4lAgRboBeZCV7sOQx5vpwnJbPO3S9ogAeu7YJA8msdzuUB11chD7lJmw4h5bTsl/jjw==")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("Sec-Fetch-Site", "cross-site")
	req.Header.Add("Sec-Fetch-Mode", "cors")
	req.Header.Add("Sec-Fetch-Dest", "empty")
	req.Header.Add("Cookie", "pass=" + shreyPassCookie + "; _shreyauth_session=" + profile.DiscordSession)
	resp, err = client.Do(req)
	if err != nil {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to submit checkout!"))
		return
	}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to submit checkout!"))
		return
	}
	r = regexp.MustCompile(`\?cid=([^"]*)`)
	match = r.FindStringSubmatch(string(body))
	if len(match) != 2 {
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to submit checkout! (Could be OOS)"))
		return
	}

	rurl = site + "subscriptions/poll?cid=" + match[1]

	stopTime := time.Now()
	diff := stopTime.Sub(beginTime)

	for i := 0; i < 100; i++ {
		fmt.Println(colors.TaskPrefix(id) + colors.Yellow("Processing..."))
		req, err := http.NewRequest("GET", rurl, nil)
		if err != nil {
			fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to get processing status!"))
			continue
		}
		req.Header.Add("sec-ch-ua", "\" Not;A Brand\";v=\"99\", \"Google Chrome\";v=\"91\", \"Chromium\";v=\"91\"")
		req.Header.Add("Accept", "text/javascript, application/javascript, application/ecmascript, application/x-ecmascript, */*; q=0.01")
		req.Header.Add("X-CSRF-Token", profile.ShreyCSRF)
		req.Header.Add("sec-ch-ua-mobile", "?0")
		req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
		req.Header.Add("X-Requested-With", "XMLHttpRequest")
		req.Header.Add("Sec-Fetch-Site", "cross-site")
		req.Header.Add("Sec-Fetch-Mode", "cors")
		req.Header.Add("Sec-Fetch-Dest", "empty")
		req.Header.Add("Cookie", "pass=" + shreyPassCookie + "; _shreyauth_session=" + profile.DiscordSession)
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to get processing status!"))
			continue
		}
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to get processing status!"))
			continue
		}
		r = regexp.MustCompile(`.html\(("Success)`)
		match = r.FindStringSubmatch(string(body))
		if len(match) == 2 {
			fmt.Println(colors.TaskPrefix(id) + colors.Green("Successfully checked out on profile ") + colors.White("\"") + colors.Green(profile.Name) + colors.White("\""))
			go utility.NewSuccess(userData.Webhook, utility.SuccessStruct{
				Site: site,
				Module: "Shrey",
				Mode: "Normal",
				Time: diff.String(),
				Profile: profile.Name,
			})
			return
		}
		r = regexp.MustCompile(`\.text\("([^"]*)`)
		match = r.FindStringSubmatch(string(body))
		if len(match) != 2 {
			time.Sleep(time.Second)
			continue
		}
		reason := match[1]
		fmt.Println(colors.TaskPrefix(id) + colors.Red("Failed to check out with reason: ") + colors.White("\"") + colors.Red(reason) +  colors.White("\""))
		return
	}	
}