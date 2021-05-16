package modules

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/chrigeeel/satango/colors"
	tls "github.com/refraction-networking/utls"
	"github.com/x04/cclient"
	"golang.design/x/clipboard"
)

func CoolClient(proxy string) http.Client {
	prxSlices := strings.Split(proxy, ":")
	var proxyFormatted string
	if len(prxSlices) == 4 {
		proxyFormatted = "http://" + prxSlices[2] + ":" + prxSlices[3] + "@" + prxSlices[0] + ":" + prxSlices[1]
	} else if len(prxSlices) == 2 {
		proxyFormatted = "http://" + prxSlices[0] + ":" + prxSlices[1]
	} else if proxy == "localhost" {
		proxyFormatted = ""
	} else {
		fmt.Println(colors.Prefix() + colors.White("Invalid Proxy: "+proxy))
		client, err := cclient.NewClient(tls.HelloRandomizedNoALPN)
		if err != nil {
			log.Fatal(err)
		}
		return client
	}
	if proxyFormatted != "" {
		client, err := cclient.NewClient(tls.HelloRandomizedNoALPN, proxyFormatted)
		if err != nil {
			log.Fatal(err)
		}
		return client
	}
	client, err := cclient.NewClient(tls.HelloRandomizedNoALPN)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func GetPw() string {
	fmt.Println(colors.Prefix() + colors.Red("Waiting for Password"))
	clipOldB := clipboard.Read(clipboard.FmtText)
	clipOld := string(clipOldB)
	r := regexp.MustCompile("password=(.*)")
	for gotPw := false; gotPw == false; {
		clipNewB := clipboard.Read(clipboard.FmtText)
		clipNew := string(clipNewB)
		if clipNew != clipOld && clipNew != "" {
			clipOld = clipNew
			gotPw = true
		}
		time.Sleep(time.Microsecond * 10)
	}
	m := r.FindStringSubmatch(clipOld)
	if len(m) == 2 {
		clipOld = m[1]
	}
	fmt.Println(colors.Prefix() + colors.White("Detected password \"") + colors.Red(clipOld) + colors.White("\""))
	return clipOld
}
