package schedule

import (
	"context"
	"github.com/busgo/pink/executor"
	"github.com/busgo/pink/pkg/bus"
	"github.com/busgo/pink/pkg/etcd"
	"github.com/busgo/pink/pkg/protocol"
	"github.com/busgo/pink/pkg/protocol/builder"
	"github.com/robfig/cron"
	"log"
	"strings"
	"sync"
	"time"
)

type PinkScheduler struct {
	schedulePlans map[string]*protocol.SchedulePlan
	changeEventCh chan *protocol.JobChangeEvent
	pinkExecutor  *executor.PinkExecutor
	etcdCli       *etcd.Cli
	sync.RWMutex
	state    int
	eventBus *bus.EventBus
}

func NewPinkScheduler(etcdCli *etcd.Cli, eventBus *bus.EventBus, pinkExecutor *executor.PinkExecutor) *PinkScheduler {
	sch := &PinkScheduler{schedulePlans: make(map[string]*protocol.SchedulePlan),
		changeEventCh: make(chan *protocol.JobChangeEvent, 16),
		pinkExecutor:  pinkExecutor,
		etcdCli:       etcdCli,
		state:         protocol.Follower,
		eventBus:      eventBus,
		RWMutex:       sync.RWMutex{},
	}
	go sch.scheduleLoop()
	go sch.lookup()

	err := sch.loadJobConf()
	if err != nil {
		return nil
	}
	return sch
}

// load all job conf
func (sch *PinkScheduler) loadJobConf() error {

	_, values, err := sch.etcdCli.GetWithPrefix(context.Background(), protocol.JobConfPath)
	if err != nil {
		return err
	}
	if len(values) == 0 {
		return nil
	}

	for _, content := range values {

		sch.PushJobChangeEvent(&protocol.JobChangeEvent{
			Event:   protocol.JobCreateChangeEventType,
			Content: new(protocol.JobConf).Decode(content),
		})
	}
	return nil
}

func (sch *PinkScheduler) lookup() {

	response := sch.etcdCli.WatchWithPrefix(protocol.JobConfPath)
	for {
		select {
		case ch := <-response.KeyChangeCh:

			value := ch.Value

			content := new(protocol.JobConf)
			if value != "" {
				content = content.Decode(value)
			}

			jobConfChangeEvent := &protocol.JobChangeEvent{
				Content: content,
			}
			switch ch.Event {
			case etcd.KeyCreateChangeEvent:
				jobConfChangeEvent.Event = protocol.JobCreateChangeEventType
			case etcd.KeyUpdateChangeEvent:
				jobConfChangeEvent.Event = protocol.JobUpdateChangeEventType
			case etcd.KeyDeleteChangeEvent:
				jobConfChangeEvent.Event = protocol.JobDeleteChangeEventType
				content.Id = getSuffixName(ch.Key)
			}
			sch.changeEventCh <- jobConfChangeEvent
		}
	}
}

func getSuffixName(key string) string {

	pos := strings.LastIndex(key, "/")
	if pos == -1 {
		return key
	}
	return string(([]byte(key))[pos+1:])
}

func (sch *PinkScheduler) scheduleLoop() {

	sub := sch.eventBus.Subscribe(protocol.NodeStateChangeTopic)
	timer := time.NewTimer(time.Second)
	defer timer.Stop()

	for {

		select {
		case <-timer.C:
			duration := sch.trySchedule()
			log.Printf("the pink schedule plan %ds after to try schedule", duration/time.Second)
			timer.Reset(duration)
		case event := <-sch.changeEventCh:
			log.Printf("the pink scheduler handle job conf change event:%+v", event)
			sch.changeEventHandle(event)
			duration := sch.trySchedule()
			log.Printf("the pink scheduler plan %ds after to try schedule", duration/time.Second)
			timer.Reset(duration)
		case state := <-sub.Receive():
			log.Printf("the pink scheduler receive the pink node state change event %d", state)
			sch.state = state.(int)
		}
	}
}

// try schedule the schedule plains
func (sch *PinkScheduler) trySchedule() time.Duration {

	sch.Lock()
	defer sch.Unlock()
	if len(sch.schedulePlans) == 0 {
		return time.Second
	}

	now := time.Now()
	var leastTime *time.Time
	for _, plan := range sch.schedulePlans {
		nextTime := plan.NextTime
		if now.After(nextTime) {
			if sch.state == protocol.Leader {
				executePlan := builder.NewSchedulePlanSnapshot(plan, now)
				sch.pinkExecutor.PushExecutePlan(executePlan)
				//	log.Printf("the schedule plan start execute .........plan:%+v", plan)
			}
			plan.BeforeTime = nextTime
			scheduleTime := plan.Schedule.Next(now)
			plan.NextTime = scheduleTime
		}
		if leastTime == nil {
			leastTime = &plan.NextTime
		}

		if leastTime.After(plan.NextTime) {
			leastTime = &plan.NextTime
		}
	}

	if leastTime == nil || leastTime.Before(now) {
		return time.Second
	}
	return leastTime.Sub(now)
}

// job change event handle
func (sch *PinkScheduler) changeEventHandle(event *protocol.JobChangeEvent) {

	switch event.Event {

	case protocol.JobCreateChangeEventType:
		sch.createEventHandle(event.Content)
	case protocol.JobUpdateChangeEventType:
		sch.updateEventHandle(event.Content)
	case protocol.JobDeleteChangeEventType:
		sch.deleteEventHandle(event.Content)
	}
}

// update event handle
func (sch *PinkScheduler) deleteEventHandle(j *protocol.JobConf) {

	plan := sch.getSchedulePlan(j.Id)
	if plan == nil {
		log.Printf("the pink scheduler found the schedule  plan  is nil,details:%+v", j)
		return
	}
	sch.Lock()
	log.Printf("the pink scheduler start delete the schedule plan:%+v", plan)
	delete(sch.schedulePlans, j.Id)
	sch.Unlock()

}

// update event handle
func (sch *PinkScheduler) updateEventHandle(j *protocol.JobConf) {

	plan := sch.getSchedulePlan(j.Id)
	if plan == nil {
		if protocol.JobState(j.State) == protocol.JobNormalState {
			log.Printf("the pink scheduler found schedule plan is nil,sudo create create event plan:%+v", j)
			sch.createEventHandle(j)
			return
		}
		log.Printf("the pink scheduler found schedule plan is nil, the job conf state  is stop %+v", j)
		return
	}

	if j.State != int32(protocol.JobNormalState) {

		log.Printf("the job state is stop start delete schedule plan,details:%+v", j)
		sch.deleteEventHandle(j)
		return
	}

	if plan.Version > j.Version {
		log.Printf("the pink scheduler found schedule plan version:%d > job version:%d", plan.Version, j.Version)
		return
	}

	schedule, err := cron.Parse(j.Cron)
	if err != nil {
		log.Printf("the pink scheduler parse the job cron is error,details:%+v,err:%+v", j, err)
		return
	}

	schedulePlan := builder.NewSchedulePlan(j, schedule)
	schedulePlan.BeforeTime = plan.BeforeTime
	log.Printf("the pink scheduler start handle update schedule plan:%+v", schedulePlan)
	sch.Lock()
	sch.schedulePlans[j.Id] = schedulePlan
	sch.Unlock()

}

// create event handle
func (sch *PinkScheduler) createEventHandle(j *protocol.JobConf) {

	plan := sch.getSchedulePlan(j.Id)
	if plan != nil {
		log.Printf("the schedule plan has exists ,j:%+v,plan:%+v", j, plan)
		return
	}

	if j.State != int32(protocol.JobNormalState) {
		log.Printf("the job state is stop,details:%+v", j)
		return
	}

	schedule, err := cron.Parse(j.Cron)
	if err != nil {
		log.Printf("parse the job cron is error,details:%+v,err:%+v", j, err)
		return
	}

	schedulePlan := builder.NewSchedulePlan(j, schedule)
	log.Printf("the create event handle create schedule plan:%+v", schedulePlan)
	sch.Lock()
	sch.schedulePlans[j.Id] = schedulePlan
	sch.Unlock()

}

// push job change event
func (sch *PinkScheduler) PushJobChangeEvent(event *protocol.JobChangeEvent) {
	sch.changeEventCh <- event
}

// get schedule plan by id
func (sch *PinkScheduler) getSchedulePlan(id string) *protocol.SchedulePlan {
	sch.RLock()
	defer sch.RUnlock()
	return sch.schedulePlans[id]
}

// get add schedule plan
func (sch *PinkScheduler) GetAllSchedulePlan() []*protocol.SchedulePlan {

	sch.RLock()
	defer sch.RUnlock()

	if len(sch.schedulePlans) == 0 {
		return make([]*protocol.SchedulePlan, 0)
	}
	plans := make([]*protocol.SchedulePlan, 0)
	for _, plan := range sch.schedulePlans {
		plans = append(plans, plan)
	}
	return plans
}
