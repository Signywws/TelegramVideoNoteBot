package dispatcher

import (
	"videonotebot/internal/Presentation_Layer/clients"
	"videonotebot/internal/Presentation_Layer/handlers"
	"videonotebot/internal/Repository_Layer/storage"
	"videonotebot/internal/Service_Layer/service"
	"videonotebot/internal/pool"
)

type Dispatcher struct {
	Text  *handlers.TextHandler
	Video *handlers.VideoHandler
	Other *handlers.OtherHandler
}

func NewDispatcher(client *clients.Client, pool *pool.Pool, processor *service.VideoProcessor,
	filestore storage.FileStorage) *Dispatcher {
	return &Dispatcher{
		Text:  handlers.NewTextHandler(client),
		Video: handlers.NewVideoHandler(client, pool, processor),
		Other: handlers.NewOtherHandler(client),
	}
}

func (d *Dispatcher) HandleUpdate(update *clients.Update) {
	if update.Message == nil {
		return
	}

	msg := update.Message

	switch {
	case msg.Text != "":
		d.Text.Handle(msg)
	case msg.Video != nil:
		d.Video.Handle(msg)
	default:
		d.Other.Handle(msg)
	}
}
