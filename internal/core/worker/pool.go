package worker

import (
	"log"

	"github.com/apt-tool/apt-core/internal/config/ftp"
	"github.com/apt-tool/apt-core/internal/core/ai"
	"github.com/apt-tool/apt-core/pkg/client"
	"github.com/apt-tool/apt-core/pkg/models"
)

type Pool struct {
	cfg    ftp.Config
	client client.HTTPClient
	models *models.Interface

	capacity int
	inuse    int
	channel  chan int
	done     chan int
}

func New(cfg ftp.Config, client client.HTTPClient, models *models.Interface, capacity int) *Pool {
	return &Pool{
		cfg:      cfg,
		client:   client,
		models:   models,
		capacity: capacity,
		inuse:    0,
		channel:  make(chan int),
		done:     make(chan int),
	}
}

func (p *Pool) update() {
	for {
		id := <-p.done

		p.inuse--

		log.Printf("[worker.update] finished process for id=%d\n", id)
	}
}

func (p *Pool) Register() {
	aiInstance := ai.AI{}

	for i := 0; i < p.capacity; i++ {
		go func() {
			err := worker{
				ai:      &aiInstance,
				cfg:     p.cfg,
				client:  p.client,
				models:  p.models,
				channel: p.channel,
				done:    p.done,
			}.work()
			if err != nil {
				log.Printf("[worker.Register] failed to start worker error=%v\n", err)
			}
		}()
	}

	go p.update()

	log.Printf("[worker.Register] started %d workers\n", p.capacity)
}

func (p *Pool) Do(id int) bool {
	if p.inuse == p.capacity {
		return false
	}

	p.inuse++

	p.channel <- id

	log.Printf("[worker.Do] start process for id=%d\n", id)

	return true
}
