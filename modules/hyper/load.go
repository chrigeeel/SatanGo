package hyper

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	"github.com/chrigeeel/satango/colors"
)

func load(site string) (hyperStruct, error) {
	loadPage := hyperStruct{}

	client := http.Client{Timeout: 10 * time.Second}
	url := site + "purchase"
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println(colors.Prefix() + colors.Red("Failed to load site "+site))
		return loadPage, errors.New("failed to load site")
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	r, _ := regexp.Compile("__NEXT_DATA__\" type=\"application\\/json\">({.*})")
	mdata := r.FindStringSubmatch(string(body))[1]

	json.Unmarshal([]byte(mdata), &loadPage)
	return loadPage, nil
}