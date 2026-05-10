package handlers

import (
	"log"
	"strings"
	"videonotebot/internal/Presentation_Layer/clients"
)

type TextHandler struct {
	client *clients.Client
}

func NewTextHandler(client *clients.Client) *TextHandler {
	return &TextHandler{
		client: client,
	}
}

func (h *TextHandler) Handle(msg *clients.Message) {
	text := strings.TrimSpace(msg.Text)
	chatID := msg.Chat.ChatID

	log.Printf("Text from %d: %s", chatID, text)

	switch {
	case text == "/start":
		reply := "Привет! Я CircleNoteBot.\n\nЯ умею превращать обычные видео в «кружочки» (Video Note).\nПросто пришли мне видео длительностью до 60 секунд, и я пришлю тебе готовый кружок.\n\nПоддерживаются форматы MP4, MOV, AVI и другие."
		h.client.SendMessage(chatID, reply)
	case text == "/thinkdifferent":
		reply := "Хвала безумцам. Бунтарям. Смутьянам. Неудачникам. Тем, кто всегда некстати и невпопад. Тем, кто видит мир иначе. Они не соблюдают правила. Они смеются над устоями. Их можно цитировать, спорить с ними, прославлять или проклинать их. Но только игнорировать их — невозможно. Ведь они несут перемены. Они толкают человечество вперёд. И пусть кто-то говорит: безумцы, мы говорим: гении. Ведь лишь безумец верит, что он в состоянии изменить мир, — и потому меняет его. (c) Стив Джобс"
		h.client.SendMessage(chatID, reply)
	default:
		reply := "Пожалуйста, пришлите видео. Я работаю только с видео и превращаю их в кружочки."
		h.client.SendMessage(chatID, reply)
	}

}
