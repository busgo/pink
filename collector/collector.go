package collector

import (
	"context"
	"github.com/busgo/pink/db/repository"
	"github.com/busgo/pink/pkg/bus"
	"github.com/busgo/pink/pkg/etcd"
	"github.com/busgo/pink/pkg/log"
	"github.com/busgo/pink/pkg/protocol"
	"strings"
	"time"
)

// PinkCollector  collect the execute snapshot
type PinkCollector struct {
	cli                       *etcd.Cli
	eventBus                  *bus.EventBus
	state                     int
	executeSnapshotRepository *repository.ExecuteSnapshotRepository
}

// new PinkCollector
func NewPinkCollector(executeSnapshotRepository *repository.ExecuteSnapshotRepository, cli *etcd.Cli, eventBus *bus.EventBus) *PinkCollector {

	c := &PinkCollector{
		cli:                       cli,
		eventBus:                  eventBus,
		state:                     protocol.Follower,
		executeSnapshotRepository: executeSnapshotRepository,
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
		log.Infof("the pink collector ")
		return
	}

	switch kc.Event {
	case etcd.KeyUpdateChangeEvent:
		log.Debugf("the pink collector receive the execute history snapshot update event %+v", kc)
		c.dealExecuteSnapshotChange(kc)
	}
}

// deal the execute snapshot change
func (c *PinkCollector) dealExecuteSnapshotChange(kc *etcd.KeyChange) {
	if strings.TrimSpace(kc.Value) == "" {
		_ = c.cli.Delete(context.Background(), kc.Key)
		return
	}
	es := new(protocol.ExecuteSnapshot).Decode(kc.Value)
	if es.State == protocol.ExecuteSnapshotFail || es.State == protocol.ExecuteSnapshotSuccess {

		c.executeSnapshotRepository.Insert(nil)
		log.Infof("the pink collector start  transfer the execute snapshot %+v", kc)
		ctx, _ := context.WithTimeout(context.Background(), time.Second*3)
		err := c.cli.Delete(ctx, kc.Key)
		if err != nil {
			log.Infof("the pink collector transfer the execute snapshot %+v fail", kc)
			return
		}
		log.Infof("the pink collector transfer the execute snapshot %+v success", kc)
	}
}
