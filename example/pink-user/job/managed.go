package job

import (
	client "github.com/busgo/pink-go"
	"log"
)

type Managed struct {
	client *client.PinkClient
	closed chan bool
}

func NewManaged(client *client.PinkClient) *Managed {

	return &Managed{client: client, closed: make(chan bool, 0)}
}

func (m *Managed) AddJob(job client.Job) {
	m.client.Subscribe(job.Target(), job)
}

func (m *Managed) Run() error {

	log.Printf("the trade job managed has start......")
	<-m.closed
	return nil
}

func (m *Managed) Stop() {

	log.Printf("the trade job managed has stop......")
	m.closed <- true

}
