package polling

import (
	"context"
	"log"
	"time"
	"videonotebot/internal/Presentation_Layer/clients"
	"videonotebot/internal/Presentation_Layer/dispatcher"
	"videonotebot/internal/pool"
)

type Poller struct {
	client     *clients.Client
	dispatcher dispatcher.Dispatcher
	offset     int
	timeout    int // таймаут long polling
	pool       *pool.Pool
}

func NewPoller(client *clients.Client, dispatcher *dispatcher.Dispatcher, timeout int, pool *pool.Pool) *Poller {
	return &Poller{
		client:     client,
		dispatcher: *dispatcher,
		offset:     0,
		timeout:    timeout,
		pool:       pool,
	}
}

func (p *Poller) Start(ctx context.Context) {
	log.Println("Poller Started!")
	p.pool.Submit(func() {
		for {
			updates, err := p.client.GetUpdate(p.offset, p.timeout)
			if err != nil {
				log.Println("GetUpdate error: ", err)
				time.Sleep(time.Second)
				return
			}

			for _, update := range updates {
				if update.Message != nil {
					p.dispatcher.HandleUpdate(&update)
					p.offset = update.UpdateID + 1
				}
				continue
			}

		}

	})

}
