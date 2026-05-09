package clients

type Update struct {
	UpdateID int      `json:"update_id"`
	Message  *Message `json:"message"`
}

type Message struct {
	MessageID int    `json:"message_id"`
	Chat      Chat   `json:"chat"`
	Text      string `json:"text"`
	Video     *Video `json:"video"`
}

type Chat struct {
	ChatID int `json:"id"`
}

type Results struct {
	Ok      bool     `json:"ok"`
	Results []Update `json:"result"`
}

type Video struct {
	FileID   string `json:"file_id"`
	Duration int    `json:"duration"`
}

type File struct {
	FileID   string `json:"file_id"`
	FilePath string `json:"file_path"` // путь относительно https://api.telegram.org/file/bot<token>/
}
