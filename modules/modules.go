package modules

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/chrigeeel/satango/colors"
	"github.com/chrigeeel/satango/loader"
	"github.com/gofiber/fiber/v2"
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

//initializes newLink variable for use in GetPw)_ & 
var newLink string
var lookingForPw bool

func GetPw() string {
	fmt.Println(colors.Prefix() + colors.Red("Waiting for Password"))
	fmt.Println(colors.Prefix() + colors.White("Copy the text \"") + colors.Red("exit") + colors.White("\" to exit"))
	lookingForPw = true
	clipOldB := clipboard.Read(clipboard.FmtText)
	clipOld := string(clipOldB)
	oldLink := newLink
	r := regexp.MustCompile("password=(.*)")
	var password string
	for gotPw := false; gotPw == false; {
		clipNewB := clipboard.Read(clipboard.FmtText)
		clipNew := string(clipNewB)
		if clipNew != clipOld && clipNew != "" {
			clipOld = clipNew
			password = clipNew
			m := r.FindStringSubmatch(password)
			if len(m) == 2 {
				password = m[1]
			}
			gotPw = true
		}
		if newLink != oldLink {
			oldLink = newLink
			password = newLink
			m := r.FindStringSubmatch(password)
			if len(m) == 2 {
				password = m[1]
				gotPw = true
			}
		}
		time.Sleep(time.Microsecond * 10)
	}
	lookingForPw = false
	fmt.Println(colors.Prefix() + colors.White("Detected password \"") + colors.Red(password) + colors.White("\""))
	return password
}

func removeIndex(profiles []loader.ProfileStruct, s int) []loader.ProfileStruct {
	return append(profiles[:s], profiles[s+1:]...)
}

func askForSilent() string {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(colors.Prefix() + colors.White("> "))
	scanner.Scan()
	return scanner.Text()
}


func Server() {

	app := fiber.New()

	app.Post("/sendpass", handleLink)
	app.Listen(":5000")
}

func handleLink(c *fiber.Ctx) error {
	type linkStruct struct {
		Link string `json:"link"`
	}

	link := new(linkStruct)

	if err := c.BodyParser(link); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	newLink = link.Link
	if lookingForPw == true {
		fmt.Println(colors.Prefix() + colors.White("You opened the link \"") + colors.Red(link.Link) + colors.White("\""))
	}
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"success": true,
	})
}