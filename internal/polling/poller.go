package polling

import (
	"context"
	"log"
	"time"
	"videonotebot/internal/telegram"
)

type Poller struct {
	client  *telegram.Client
	offset  int
	timeout int
}

func New(client *telegram.Client, timeout int) *Poller {
	return &Poller{
		client:  client,
		offset:  0,
		timeout: timeout,
	}
}

func (p *Poller) Start(ctx context.Context) {
	log.Print("Poller started!")
	for {
		select {
		case <-ctx.Done():
			log.Println("Poller stoped")
			return
		default:
		}

		updates, err := p.client.GetUpdates(p.offset, p.timeout)
		if err != nil {
			log.Printf("GetUpdates error: %v", err)
			time.Sleep(2 * time.Second) // пауза перед повторной попыткой
			continue
		}

		for _, upd := range updates {
			if upd.Message == nil {
				continue
			}
			if upd.Message.Text != "" {
				log.Printf("[%d] %s", upd.Message.Chat.ID, upd.Message.Text)
				err := p.client.SendMessage(upd.Message.Chat.ID, "Привет! Я бот для создания круглых видео. Отправь мне видео до 60 секунд.")
				if err != nil {
					log.Printf("SendMessage error: %v", err)
				}
			}
			p.offset = upd.UpdateID + 1
		}
	}
}
