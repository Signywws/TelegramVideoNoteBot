package telegram

type Update struct { // Событие, которое Telegram присылает в getUpdates
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
	ID int `json:"id"`
}

type Video struct {
	FileID   string `json:"file_id"`
	Duration int    `json:"duration"`
}
