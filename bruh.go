package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	http "github.com/zMrKrabz/fhttp"
)

func main() {

  url := "https://dashboard.purepings.eu/password"
  method := "POST"

  payload := strings.NewReader("authenticity_token=sEmbJAWHTES9gxnpcPokjKcFDmEPaQjN934MLkJC9EjqN%2BaMLFBPhsp%2BFJXEqgxgYO9QlJlfHRscenxpRpxcJg%3D%3D&password=PurePings2021")

  client := &http.Client{
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	},
}
  req, err := http.NewRequest(method, url, payload)

  if err != nil {
    fmt.Println(err)
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
  req.Header.Add("Cookie", "_shreyauth_session=ys6bbk%2BJXGEwJF6vvmSYtItWe69oblD0MVta4CcxcFONtGYp%2FyyM7Od6amBfOgCN%2FGf7536MpoCLyBs9Xt5Wr1Pm5eh3vSNJ8cehTENmQbDL3JzX0Ih6Hht%2F8OexhTmKjINUG0a4bDuYCHvFzGBJ6litxM3R44gPbKrifyBpzq9RMzNZgDu3uo9edeC1fcEllmlwiBDnKo8cUShCP8ADymX4%2FAk0ZoUmb957GUQMHaNWwJKD5LhwQool8uejgol8xAPwYqObqA7Q2WiUfA9V647S8A%3D%3D--cw%2BhzhAGmgXW%2FkuT--39yU4eKAdYRiddxt%2BKtt7A%3D%3D")

  res, err := client.Do(req)
  if err != nil {
    fmt.Println(err)
    return
  }
  defer res.Body.Close()

  body, err := ioutil.ReadAll(res.Body)
  if err != nil {
    fmt.Println(err)
    return
  }
  fmt.Println(string(body))
  fmt.Println(res.Cookies())
  fmt.Println(res.Header)
  fmt.Println(res.Header.Get("Set-Cookie"))
}