package polling

import (
	"context"
	"log"
	"time"
	"videonotebot/internal/Presentation_Layer/clients"
	"videonotebot/internal/Presentation_Layer/dispatcher"
)

type Poller struct {
	client     *clients.Client
	dispatcher dispatcher.Dispatcher
	offset     int
	timeout    int // таймаут long polling
	updates    chan *clients.Update
}

func NewPoller(client *clients.Client, dispatcher *dispatcher.Dispatcher, timeout int) *Poller {
	return &Poller{
		client:     client,
		dispatcher: *dispatcher,
		offset:     0,
		timeout:    timeout,
		updates:    make(chan *clients.Update, 1000),
	}
}

func (p *Poller) Start(ctx context.Context) {
	log.Println("Poller Started!")

	go func() {
		select {
		case <-ctx.Done():
			log.Println("Poller stopped")
			return
		default:
		}
		for update := range p.updates {
			p.dispatcher.HandleUpdate(update)
		}
		log.Println("Update consumer stopped")
	}()

	go func() {
		defer close(p.updates)
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
				time.Sleep(time.Second)
				continue
			}

			for _, update := range updates {
				if update.Message == nil {
					continue
				}
				upd := update
				p.updates <- &upd // безопасный указатель
				p.offset = update.UpdateID + 1
			}

		}

	}()
}
