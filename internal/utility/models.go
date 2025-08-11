package utility

var (
	BotToken    string
	ChatID      string
	PhoneNumber string
	ApiKey      string
	XApiKey     string
)

type HolodexScraper struct {
	VideoInfos []APIVideoInfo
}

type Channel struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Org    string `json:"org"`
	Suborg string `json:"suborg"`
}

type APIVideoInfo struct {
	ID             string  `json:"id"` // It's also the youtube link
	Title          string  `json:"title"`
	Type           string  `json:"type"`     // "stream" or "placeholder"
	TopicID        string  `json:"topic_id"` // e.g. "singing"
	Duration       int     `json:"duration"` // seconds
	Status         string  `json:"status"`   // "upcoming", "live", etc.
	StartScheduled string  `json:"start_scheduled"`
	StartActual    string  `json:"start_actual"`
	Channel        Channel `json:"channel"`
}
