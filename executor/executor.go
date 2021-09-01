package executor

import (
	"context"
	"fmt"
	"github.com/busgo/pink/pkg/bus"
	"github.com/busgo/pink/pkg/etcd"
	"github.com/busgo/pink/pkg/log"
	"github.com/busgo/pink/pkg/protocol"
	"github.com/busgo/pink/pkg/util"

	"time"
)

type PinkExecutor struct {
	etcdCli                *etcd.Cli
	scheduleSnapshotPlanCh chan *protocol.SchedulePlanSnapshot
	groupManaged           *PinkGroupManaged
	state                  int
	eventBus               *bus.EventBus
}

func NewPinkExecutor(etcdCli *etcd.Cli, eventBus *bus.EventBus, groupManaged *PinkGroupManaged) *PinkExecutor {
	pe := &PinkExecutor{etcdCli: etcdCli,
		scheduleSnapshotPlanCh: make(chan *protocol.SchedulePlanSnapshot, 32),
		groupManaged:           groupManaged,
		state:                  protocol.Follower,
		eventBus:               eventBus,
	}
	go pe.lookup()
	return pe
}

func (pe *PinkExecutor) PushExecutePlan(snapshot *protocol.SchedulePlanSnapshot) {
	pe.scheduleSnapshotPlanCh <- snapshot
}

func (pe *PinkExecutor) lookup() {

	ticker := time.Tick(time.Second * 30)
	sub := pe.eventBus.Subscribe(protocol.NodeStateChangeTopic)
	for {
		select {
		case snapshot := <-pe.scheduleSnapshotPlanCh:
			_ = pe.dealScheduleSnapshot(snapshot)
		case <-ticker:
			pe.scanScheduleSnapshot()
		case state := <-sub.Receive():
			log.Infof("the pink executor receive the pink node state change event %d", state)
			pe.state = state.(int)
		}
	}
}

// deal schedule snapshot
func (pe *PinkExecutor) dealScheduleSnapshot(snapshot *protocol.SchedulePlanSnapshot) error {
	ip, err := pe.groupManaged.UnParkPinkClient(snapshot.Group)
	if err != nil {
		log.Errorf("the pink executor deal snapshot %+v unpark the pink client fail %+v", snapshot, err)
		return pe.addScheduleSnapshot(snapshot)
	}
	log.Debugf("execute ip ----> %s", ip)
	return pe.execute(ip, snapshot)
}

// scan schedule snapshot
func (pe *PinkExecutor) scanScheduleSnapshot() {

	if pe.state != protocol.Leader {
		return
	}
	log.Infof("the pink executor start scan schedule snapshot")
	ctx, _ := context.WithTimeout(context.Background(), time.Second*3)
	keys, values, err := pe.etcdCli.GetWithPrefix(ctx, protocol.ScheduleSnapshotPath)
	if err != nil {
		log.Errorf("the pink executor start scan schedule snapshot  fail %+v", err)
		return
	}

	if len(values) == 0 {
		log.Warnf("the pink executor start scan schedule snapshot the retry schedule snapshot not found")
		return
	}

	for pos, value := range values {
		snapshot := new(protocol.SchedulePlanSnapshot).Decode(value)
		if pe.dealScheduleSnapshot(snapshot) == nil {
			log.Warnf("the pink executor start delete schedule snapshot %s", value)
			_ = pe.etcdCli.Delete(ctx, keys[pos])
		}
	}

}

// add schedule snapshot
func (pe *PinkExecutor) addScheduleSnapshot(snapshot *protocol.SchedulePlanSnapshot) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*3)
	path := fmt.Sprintf("%s%s", protocol.ScheduleSnapshotPath, snapshot.Id)
	err := pe.etcdCli.PutWithNotExist(ctx, path, snapshot.Encode())
	if err != nil {
		log.Errorf("the pink executor put schedule snapshot %+v fail %+v", snapshot, err)
		return err
	}
	return nil
}

// execute schedule snapshot
func (pe *PinkExecutor) execute(ip string, snapshot *protocol.SchedulePlanSnapshot) error {

	executeSnapshot := &protocol.ExecuteSnapshot{
		Id:           util.Generate(),
		JobId:        snapshot.Id,
		Name:         snapshot.Name,
		Group:        snapshot.Group,
		Cron:         snapshot.Cron,
		Target:       snapshot.Target,
		Ip:           ip,
		Param:        snapshot.Param,
		State:        protocol.ExecuteSnapshotInit,
		BeforeTime:   snapshot.BeforeTime,
		ScheduleTime: snapshot.ScheduleTime,
		StartTime:    "",
		EndTime:      "",
		Times:        0,
		Mobile:       snapshot.Mobile,
		Version:      0,
		Remark:       snapshot.Remark,
	}

	path := fmt.Sprintf(protocol.ExecuteSnapshotPath, snapshot.Group, ip, executeSnapshot.Id)
	content := executeSnapshot.Encode()
	ctx, _ := context.WithTimeout(context.Background(), time.Second*3)
	err := pe.etcdCli.PutWithNotExist(ctx, path, content)
	if err != nil {
		log.Errorf("the pink executor create execute plan snapshot %s  for ip %s fail err %+v", snapshot.Encode(), ip, err)
		return pe.addScheduleSnapshot(snapshot)
	}
	log.Infof("the pink executor create execute plan snapshot %s  for ip %s success", content, ip)
	return nil
}
