package polling

import (
	"context"
	"log"
	"videonotebot/internal/Presentation_Layer/clients"
	"videonotebot/internal/Presentation_Layer/dispatcher"
)

type Poller struct {
	client     *clients.Client
	dispatcher dispatcher.Dispatcher
	offset     int
	timeout    int // таймаут long polling
}

func NewPoller(client *clients.Client, dispatcher *dispatcher.Dispatcher, timeout int) *Poller {
	return &Poller{
		client:     client,
		dispatcher: *dispatcher,
		offset:     0,
		timeout:    timeout,
	}
}

func (p *Poller) Start(ctx context.Context) {
	log.Println("Poller Started!")

	for {
		select {
		case <-ctx.Done():
			log.Println("Poller stopped")
			return
		default:
		}

		updates, err := p.client.GetUpdate(p.offset, p.timeout)
		if err != nil {
			log.Printf("GetUpdates error: %v", err)
		}

		for _, update := range updates {
			if update.Message == nil {
				continue
			}

			p.dispatcher.HandleUpdate(&update)
			// обновляеем до последнего обработанного
			p.offset = update.UpdateID + 1
		}

	}
}
