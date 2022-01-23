package utility

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/chrigeeel/satango/structs"
)

type SuccessStruct struct {
	Site    string `json:"site"`
	Module  string `json:"module"`
	Mode    string `json:"mode"`
	Time    string `json:"time"`
	Profile string `json:"profile"`
}

func NewSuccess(wurl string, successData SuccessStruct) error {
	jsonData, err := json.Marshal(successData)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", "https://hardcore.astolfoporn.com/checkout", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("content-type", "application/json")
	http.DefaultClient.Do(req)

	if wurl == "" {
		return nil
	}

	webhookData := structs.Webhook{
		Username: "Satan Success",
		AvatarURL: "https://i.imgur.com/sNZJzJ2.png",
		Embeds: []*structs.Embed{
			{
				Title: ":smiling_imp: Successfully Sacrificed Stock :smiling_imp:",
				Color: 5294200,
				Fields: []*structs.Field{
					{
						Name: "Site",
						Value: successData.Site,
					},
					{
						Name: "Site Module",
						Value: successData.Module,
					},
					{
						Name: "Mode",
						Value: successData.Mode,
					},
					{
						Name: "Checkout Time",
						Value: successData.Time,
					},
					{
						Name: "Profile",
						Value: "||" + successData.Profile + "||",
					},
				},
			},
		},
	}

	jsonData, err = json.Marshal(webhookData)
	if err != nil {
		return err
	}
	req, err = http.NewRequest("POST", wurl, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("content-type", "application/json")
	http.DefaultClient.Do(req)

	return nil
}