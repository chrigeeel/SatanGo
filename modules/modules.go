package modules

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
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
	"golang.org/x/net/websocket"
)

type Message struct {
	Username string `json:"username"`
	Password string `json:"password"`
	ReleaseId string `json:"releaseId"`
	Site string `json:"site"`
}

//initializes variables for use in GetPw & other
var (
	newLink string
	newShare Message
	lookingForPw bool
	WS *websocket.Conn
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

func GetPw(site string) string {
	fmt.Println(colors.Prefix() + colors.Red("Waiting for Password"))
	fmt.Println(colors.Prefix() + colors.White("Copy the text \"") + colors.Red("exit") + colors.White("\" to exit"))
	lookingForPw = true
	clipOldB := clipboard.Read(clipboard.FmtText)
	clipOld := string(clipOldB)
	newLink = "bruh"
	oldLink := newLink
	oldShare := newShare
	r := regexp.MustCompile("password=(.*)")
	var password string
	password = "bruh"
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
		if (newShare != Message{}) && (oldShare != Message{}) {
			if (newShare.Password != oldShare.Password) && (newShare.Site == site) {
				fmt.Println(colors.Prefix() + colors.Green("Received password via pw sharing: ") + colors.White("\"") + colors.Green(newShare.Password) + colors.White("\"!"))
				fmt.Println(colors.Prefix() + colors.Green("The password was sent to you by the user ") + colors.White("\"") + colors.Green(newShare.Username) + colors.White("\", go say thanks!"))
				password = newShare.Password
				m := r.FindStringSubmatch(password)
				if len(m) == 2 {
					password = m[1]
				}
				gotPw = true
			}
		}
		time.Sleep(time.Microsecond * 5)
	}
	lookingForPw = false
	password = strings.ReplaceAll(password, " ", "")
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

func PwSharingReceive() {
	fmt.Println(colors.Prefix() + colors.Yellow("Connecting to PW-Sharing server..."))
	WS, err := connect()
	if err != nil {
	  fmt.Println(colors.Red("Failed to connect to PW-Sharing server!"))
	  return
	}
	defer WS.Close()
	if WS != nil {
		fmt.Println(colors.Prefix() + colors.Green("Successfully connected to PW-Sharing server!"))
	}
	var m Message
	for {
		err := websocket.JSON.Receive(WS, &m)
		if err != nil {
			fmt.Println(colors.Prefix()  + colors.Red("Error receiving password!"))
			fmt.Println(colors.Prefix() + colors.Green("PLEASE DM CHRIGEEEL"))
			time.Sleep(time.Second * 1000)
			return
		}
		if lookingForPw == true {
			newShare = m
		}
	}
}

func PwSharingSend(password string, username string, site string) {
	m := Message{
		Username: username,
		Password: password,
		Site: site,
	}
	err := websocket.JSON.Send(WS, m)
	if err != nil {
		fmt.Println(colors.Prefix() + colors.Red("Failed sending password to PW-Sharing!"))
		return
	}
	fmt.Println(colors.Prefix() + colors.Yellow("Successfully sent password to PW-Sharing"))
}

func connect() (*websocket.Conn, error) {
	return websocket.Dial(fmt.Sprintf("ws://52.72.153.196"), "", mockedIP())
}

func mockedIP() string {
	var arr [4]int
	for i := 0; i < 4; i++ {
		rand.Seed(time.Now().UnixNano())
		arr[i] = rand.Intn(256)
	}
	return fmt.Sprintf("http://%d.%d.%d.%d", arr[0], arr[1], arr[2], arr[3])
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