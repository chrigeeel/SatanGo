package modules

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)


type fields struct {
	Name string `json:"name"`
	Value string `json:"value"`
}

type thumbnail struct {
	Url string `json:"url"`
}

type footer struct {
	Text string `json:"text"`
	Icon_url string `json:"icon_url"`
}

type embeds struct {
	Color string `json:"color"`
	Title string `json:"title"`
	Fields []fields `json:"fields"`
	Thumbnail thumbnail `json:"thumbnail"`
	Footer footer `json:"footer"`
}

type webhook struct {
	Username string `json:"username"`
	Avatar_url string `json:"avatar_url"`
	Embeds []embeds `json:"embeds"`
}

type WebhookContentStruct struct {
	Speed string `json:"speed"`
	Module string `json:"module"`
	Site string `json:"site"`
	Profile string `json:"profile"`
}

func SendWebhook(url string, content WebhookContentStruct) {
	if url == "" {
		return
	}
	webhookData := webhook{}
	webhookData.Username = "Satanbots Success"
	webhookData.Avatar_url = "https://imgur.com/dyj4onp.png"
	var fieldsData []fields
	fieldsData = append(fieldsData, fields{
		Name: "Checkout Speed",
		Value: content.Speed,
	})
	fieldsData = append(fieldsData, fields{
		Name: "Module",
		Value: content.Module,
	})
	fieldsData = append(fieldsData, fields{
		Name: "Site",
		Value: content.Site,
	})
	fieldsData = append(fieldsData, fields{
		Name: "Profile",
		Value: content.Profile,
	})
	var thumbnailData thumbnail
	thumbnailData.Url = "https://imgur.com/dyj4onp.png"
	var footerData footer
	footerData.Icon_url = "https://imgur.com/dyj4onp.png"
	footerData.Text = "SatanBot"
	webhookData.Embeds = append(webhookData.Embeds, embeds{
		Color: "9573391",
		Title: "Successfully Sacrificed Stock",
		Fields: fieldsData,
		Thumbnail: thumbnailData,
		Footer: footerData,
	})
	body, err := json.Marshal(webhookData)
	if err != nil {
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	return
}