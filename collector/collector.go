package collector

import (
	"context"
	"github.com/busgo/pink/db/model"
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
	cli                          *etcd.Cli
	eventBus                     *bus.EventBus
	state                        int
	executeSnapshotHisRepository *repository.ExecuteSnapshotHisRepository
}

// new PinkCollector
func NewPinkCollector(executeSnapshotHisRepository *repository.ExecuteSnapshotHisRepository, cli *etcd.Cli, eventBus *bus.EventBus) *PinkCollector {

	c := &PinkCollector{
		cli:                          cli,
		eventBus:                     eventBus,
		state:                        protocol.Follower,
		executeSnapshotHisRepository: executeSnapshotHisRepository,
	}

	go c.lookup()
	return c
}

// look up execute snapshot
func (c *PinkCollector) lookup() {
	resp := c.cli.WatchWithPrefix(protocol.ExecuteSnapshotHistoryPath)
	sub := c.eventBus.Subscribe(protocol.NodeStateChangeTopic)

	log.Infof("the pink collector lookup for %s", protocol.ExecuteSnapshotHistoryPath)
	for {
		select {
		case ch := <-resp.KeyChangeCh:
			c.handleExecuteSnapshotHisChange(ch)
		case rec := <-sub.Receive():
			c.state = rec.(int)
		}
	}
}

// handle execute snapshot change
func (c *PinkCollector) handleExecuteSnapshotHisChange(kc *etcd.KeyChange) {
	if c.state != protocol.Leader {
		log.Warnf("the pink collector found the pink node state is not leader")
		return
	}

	switch kc.Event {
	case etcd.KeyUpdateChangeEvent:
		log.Debugf("the pink collector receive the execute history snapshot update event %+v", kc)
		c.dealExecuteSnapshotHisChange(kc)
	case etcd.KeyCreateChangeEvent:
		log.Debugf("the pink collector receive the execute history snapshot delete event %+v", kc)
		c.dealExecuteSnapshotHisChange(kc)
	}
}

// deal the execute snapshot change
func (c *PinkCollector) dealExecuteSnapshotHisChange(kc *etcd.KeyChange) {
	if strings.TrimSpace(kc.Value) == "" {
		_ = c.cli.Delete(context.Background(), kc.Key)
		return
	}
	es := new(protocol.ExecuteSnapshot).Decode(kc.Value)
	if es.State == protocol.ExecuteSnapshotFail || es.State == protocol.ExecuteSnapshotSuccess {

		c.saveExecuteSnapshotHisToDB(es)
		ctx, _ := context.WithTimeout(context.Background(), time.Second*3)
		err := c.cli.Delete(ctx, kc.Key)
		if err != nil {
			log.Infof("the pink collector transfer the execute snapshot %+v fail", kc)
			return
		}
		log.Infof("the pink collector delete the execute snapshot  his %+v from etcd success", kc)
	}
}

//
func (c *PinkCollector) saveExecuteSnapshotHisToDB(snapshot *protocol.ExecuteSnapshot) {

	his := &model.ExecuteSnapshotHis{

		SnapshotId:   snapshot.Id,
		JobId:        snapshot.JobId,
		JobName:      snapshot.Name,
		Group:        snapshot.Group,
		Cron:         snapshot.Cron,
		Target:       snapshot.Target,
		Ip:           snapshot.Ip,
		Param:        snapshot.Param,
		State:        snapshot.State,
		BeforeTime:   snapshot.BeforeTime,
		ScheduleTime: snapshot.ScheduleTime,
		EndTime:      snapshot.EndTime,
		Times:        snapshot.Times,
		Mobile:       snapshot.Mobile,
		Remark:       snapshot.Remark,
	}
	if err := c.executeSnapshotHisRepository.Insert(his); err != nil {
		log.Errorf("the pink collector save  the execute snapshot his to db fail %+v fail %+v", his, err)
		return
	}
	log.Infof("the pink collector save  the execute snapshot his %+v to db success ", his)
}
