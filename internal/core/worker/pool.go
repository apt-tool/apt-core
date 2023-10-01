package worker

import (
	"log"

	"github.com/ptaas-tool/base-api/internal/config/ftp"
	"github.com/ptaas-tool/base-api/internal/core/ai"
	"github.com/ptaas-tool/base-api/pkg/client"
	"github.com/ptaas-tool/base-api/pkg/models"
)

type Pool struct {
	cfg    ftp.Config
	client client.HTTPClient
	models *models.Interface

	capacity int
	inuse    int
	channel  chan int
	reruns   chan int
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
		reruns:   make(chan int),
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
				rerun:   p.reruns,
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

func (p *Pool) Do(id int, rerun bool) bool {
	if p.inuse == p.capacity {
		return false
	}

	p.inuse++

	if rerun {
		p.reruns <- id
	} else {
		p.channel <- id
	}

	log.Printf("[worker.Do] start process for id=%d\n", id)

	return true
}
