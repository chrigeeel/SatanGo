package structs

// DISCORD WEBHOOK STRUCTS ----

type Author struct {
	Name    string `json:"name,omitempty"`
	URL     string `json:"url,omitempty"`
	IconURL string `json:"icon_url,omitempty"`
}

type Field struct {
	Name   string `json:"name,omitempty"`
	Value  string `json:"value,omitempty"`
	Inline bool   `json:"inline,omitempty"`
}

type Thumbnail struct {
	URL string `json:"url,omitempty"`
}

type Image struct {
	URL string `json:"url,omitempty"`
}

type Footer struct {
	Text    string `json:"text,omitempty"`
	IconURL string `json:"icon_url,omitempty"`
}

type Embed struct {
	Author      *Author    `json:"author,omitempty"`
	Title       string     `json:"title,omitempty"`
	URL         string     `json:"url,omitempty"`
	Description string     `json:"description,omitempty"`
	Color       int        `json:"color,omitempty"`
	Fields      []*Field   `json:"fields,omitempty"`
	Thumbnail   *Thumbnail `json:"thumbnail,omitempty"`
	Image       *Image     `json:"image,omitempty"`
	Footer      *Footer    `json:"footer,omitempty"`
}

type Webhook struct {
	Username  string   `json:"username,omitempty"`
	AvatarURL string   `json:"avatar_url,omitempty"`
	Content   string   `json:"content,omitempty"`
	Embeds    []*Embed `json:"embeds,omitempty"`
}

// DISCORD WEBHOOK STRUCTS ----