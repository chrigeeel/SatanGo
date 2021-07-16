package getpw

import (
	"fmt"
	"regexp"

	"github.com/chrigeeel/satango/colors"
	"github.com/gofiber/fiber/v2"
)

func MonitorExtension() {
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

	nLink := link.Link
	if lookingForPw == false {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
		})
	}

	r := regexp.MustCompile("(?i)password=\n*?\\s*?(\\S+)")
	r2 := regexp.MustCompile("\\/purchase\\/([^\\?]*)")
	fmt.Println(colors.Prefix() + colors.White("You opened the link \"") + colors.Red(nLink) + colors.White("\""))
	
	var password string
	var releaseId string
	m := r.FindStringSubmatch(nLink)
	if len(m) == 2 {
		password = m[1]
	}
	m2 := r2.FindStringSubmatch(nLink)
	if len(m2) == 2 {
		releaseId = m2[1]
	}
	if password == "" && releaseId == "" {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
		})
	}
	hyperInfo := HyperInfo{
		ReleaseId: releaseId,
	}
	p := PWStruct{
		Password: password,
		Mode: "extension",
		HyperInfo: hyperInfo,
			
	}
	PWC <- p
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
	})
}