package hyper

type HyperStruct struct {
	Props struct {
		Pageprops struct {
			Account struct {
				Settings struct {
					Payments struct {
						CollectBillingAddress bool `json:"collect_billing_address"`
						RequireLogin          bool `json:"require_login"`
					} `json:"payments"`
					BotProtection struct {
						Enabled bool `json:"enabled"`
					} `json:"bot_protection"`
				} `json:"settings"`
				ID             string `json:"id"`
				Stripe_account string `json:"stripe_account"`
			} `json:"account"`
			Release struct {
				OutOfStock bool `json:"out_of_stock"`
			}
		} `json:"pageProps"`
	} `json:"props"`
	Query struct {
		Token   string `json:"token"`
		Release string `json:"release"`
	} `json:"query"`
}
type HyperCheckoutStruct struct {
	Billing_details struct {
		Address struct {
			City        string `json:"city,omitempty"`
			Country     string `json:"country,omitempty"`
			Line1       string `json:"line1,omitempty"`
			Line2       string `json:"line2,omitempty"`
			Postal_code string `json:"postal_code,omitempty"`
			State       string `json:"state,omitempty"`
		} `json:"address,omitempty"`
		Email string `json:"email"`
		Name  string `json:"name"`
	} `json:"billing_details,omitempty"`
	Payment_method string `json:"payment_method,omitempty"`
	User           string `json:"user,omitempty"`
	Release        string `json:"release"`
}