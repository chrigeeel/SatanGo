package utility

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/chrigeeel/satango/colors"
	tls "github.com/refraction-networking/utls"
	"github.com/x04/cclient"
)

func CoolClient(proxy string) (http.Client, error) {
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
			return client, err
		}
		return client, errors.New("invalid proxy")
	}
	if proxyFormatted != "" {
		client, err := cclient.NewClient(tls.HelloRandomizedNoALPN, proxyFormatted)
		if err != nil {
			return client, err
		}
		return client, nil
	}
	client, err := cclient.NewClient(tls.HelloRandomizedNoALPN)
	if err != nil {
		return client, err
	}
	return client, nil
}