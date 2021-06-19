package tldash

import "github.com/chrigeeel/satango/loader"

type siteStruct struct {
	DisplayName       string `json:"displayName"`
	BackendName       string `json:"backendName"`
	Stripe_public_key string `json:"stripe_public_key,omitempty"`
}

type taskStruct struct {
	Site        string               `json:"site"`
	Proxy       string               `json:"proxy"`
	Profile     loader.ProfileStruct `json:"profile"`
	StripeToken string               `json:"stripeToken"`
}