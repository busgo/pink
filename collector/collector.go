package collector

import (
	"context"
	"github.com/busgo/pink/pkg/bus"
	"github.com/busgo/pink/pkg/etcd"
	"github.com/busgo/pink/pkg/protocol"
	"log"
	"strings"
)

// PinkCollector  collect the execute snapshot
type PinkCollector struct {
	cli      *etcd.Cli
	eventBus *bus.EventBus
	state    int
}

// new PinkCollector
func NewPinkCollector(cli *etcd.Cli, eventBus *bus.EventBus) *PinkCollector {

	c := &PinkCollector{
		cli:      cli,
		eventBus: eventBus,
		state:    protocol.Follower,
	}

	go c.lookup()
	return c
}

// look up execute snapshot
func (c *PinkCollector) lookup() {
	resp := c.cli.WatchWithPrefix(protocol.ExecuteSnapshotHistoryPath)
	sub := c.eventBus.Subscribe(protocol.NodeStateChangeTopic)
	for {
		select {
		case ch := <-resp.KeyChangeCh:
			c.handleExecuteSnapshotChange(ch)
		case rec := <-sub.Receive():
			c.state = rec.(int)
		}
	}
}

// handle execute snapshot change
func (c *PinkCollector) handleExecuteSnapshotChange(kc *etcd.KeyChange) {
	if c.state != protocol.Leader {
		return
	}

	switch kc.Event {
	case etcd.KeyUpdateChangeEvent:
		log.Printf("the pink collector receive the execute history snapshot update event %+v", kc)
		c.dealExecuteSnapshotChange(kc)
	}
}

// deal the execute snapshot change
func (c *PinkCollector) dealExecuteSnapshotChange(kc *etcd.KeyChange) {
	if strings.TrimSpace(kc.Value) == "" {
		_ = c.cli.Delete(context.Background(), kc.Key)
		return
	}

	// TODO  transfer to MySQL or  Mongo db ES
	//es := new(protocol.ExecuteSnapshot).Decode(kc.Value)
	//if es.State == protocol.ExecuteSnapshotFail || es.State == protocol.ExecuteSnapshotSuccess {
	//	log.Printf("the pink collector start  transfer the execute snapshot %+v", kc)
	//	ctx, _ := context.WithTimeout(context.Background(), time.Second*3)
	//	err := c.cli.Transfer(ctx, kc.Key, fmt.Sprintf("%s%s", protocol.ExecuteSnapshotHistoryPath, es.Id), kc.Value)
	//	if err != nil {
	//		log.Printf("the pink collector transfer the execute snapshot %+v fail", kc)
	//		return
	//	}
	//	log.Printf("the pink collector transfer the execute snapshot %+v success", kc)
	//}
}
