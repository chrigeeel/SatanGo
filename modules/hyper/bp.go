package hyper

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/chrigeeel/satango/utility"
)

var (
	token string
	auth string
	key string = "KbPeShVmYq3t6w8z" 
)

func Solvebp(domain string) (string, error) {
	type bpStruct struct {
		UserAgent string `json:"userAgent"`
		Timestamp int64 `json:"timestamp"`
		IP string `json:"ip"`
		Token string `json:"token"`
		Auth string `json:"auth"`
		Challenge string `json:"challenge"`
		ChallengeFunction string `json:"challengeFunction"`
		CkHchKctl bool `json:"ck_hch_kctl"`
		CfCtl string `json:"cf_ctl"`
		Domain string `json:"domain"`
	}

	bp := bpStruct{
		UserAgent: "Smart vibrating horse cock",
		Timestamp: NowAsUnixMilli(),
		IP: "1.1.1.1",
		Token: token,
		Auth: auth,
		Challenge: "5",
		ChallengeFunction: "WkRKc2RWcEhPVE5NYkRsbVdUSm9hR0pIZUd4aWJXUnNTVVF3WjAxcFFYSkpSRTAz",
		CkHchKctl: true,
		CfCtl: "nemesisisaprettygoodfnfimo",
		Domain: domain,
	}

	payload, err := json.Marshal(bp)
	if err != nil {
		return "", errors.New("failed to solve bot protection")
	}

	bpToken, err := utility.AESEncrypt(string(payload), key)
	if err != nil {
		return "", errors.New("failed to solve bot protection")
	}
	
	return bpToken, nil
}

func Initbp() {
	type bpStruct struct {
		Token string `json:"token"`
		Auth string `json:"authenticity"`
	}
	resp, err := http.Get("https://backend-dot-hch-protection.uc.r.appspot.com/api/metalabs-lite/getChallenge")
	if err != nil {
		token = "5bda3f04-6b8b-4c09-536a-5136b6eaedb1"
		auth = "U2FsdGVkX18iupuiCE+qkVJOrmEGo1dGCwRZtPd2bXeC6dHyPdf8sgE45BB15ouy/1deTN0MoHKqtaYP5OJIm8slVxS2NiqQ9gvt+vNTxH5FXtxRDLq3au5k3qncxbk0fwNProTlWezhguUYbGcYm+SxRI6imgZgVM9SXBjvUKO3ygKaCmzCGmROrmXgnKA3"
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		token = "5bda3f04-6b8b-4c09-536a-5136b6eaedb1"
		auth = "U2FsdGVkX18iupuiCE+qkVJOrmEGo1dGCwRZtPd2bXeC6dHyPdf8sgE45BB15ouy/1deTN0MoHKqtaYP5OJIm8slVxS2NiqQ9gvt+vNTxH5FXtxRDLq3au5k3qncxbk0fwNProTlWezhguUYbGcYm+SxRI6imgZgVM9SXBjvUKO3ygKaCmzCGmROrmXgnKA3"
		return
	}
	var data bpStruct
	err = json.Unmarshal(body, &data)
	if err != nil {
		token = "5bda3f04-6b8b-4c09-536a-5136b6eaedb1"
		auth = "U2FsdGVkX18iupuiCE+qkVJOrmEGo1dGCwRZtPd2bXeC6dHyPdf8sgE45BB15ouy/1deTN0MoHKqtaYP5OJIm8slVxS2NiqQ9gvt+vNTxH5FXtxRDLq3au5k3qncxbk0fwNProTlWezhguUYbGcYm+SxRI6imgZgVM9SXBjvUKO3ygKaCmzCGmROrmXgnKA3"
		return
	}
	token, auth = data.Token, data.Auth
}

func NowAsUnixMilli() int64 {
    return time.Now().UnixNano() / 1e6
}