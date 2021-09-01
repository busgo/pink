package node

import (
	"context"
	"errors"
	"fmt"
	"github.com/busgo/pink/executor"
	"github.com/busgo/pink/pkg/bus"
	"github.com/busgo/pink/pkg/etcd"
	"github.com/busgo/pink/pkg/log"
	"github.com/busgo/pink/pkg/protocol"
	"github.com/busgo/pink/pkg/util"
	"github.com/busgo/pink/schedule"
	"github.com/segmentio/ksuid"
	"go.etcd.io/etcd/client/v3/concurrency"
	"time"
)

type State int32

const (
	Follower State = iota + 1
	Leader
)

// elect
type ElectionState int32

const (
	ElectionReadyState ElectionState = iota
	ElectionDoingState
)

type PinkNode struct {
	id            string
	state         int
	electionState ElectionState
	electionPath  string
	electionTTL   int64
	instancePath  string
	instanceTTL   int64
	etcdCli       *etcd.Cli
	sch           *schedule.PinkScheduler
	pinkExecutor  *executor.PinkExecutor
	eventBus      *bus.EventBus
	sub           *bus.Sub
	closed        chan bool
}

func NewPinkNode(cli *etcd.Cli, eventBus *bus.EventBus) (*PinkNode, error) {
	return &PinkNode{
		id:            util.GetLocalIP() + "-" + ksuid.New().String(),
		state:         protocol.Follower,
		electionState: ElectionReadyState,
		electionPath:  protocol.NodeElectionPath,
		electionTTL:   protocol.NodeElectionTTL,
		instancePath:  protocol.NodeInstancePath,
		instanceTTL:   protocol.NodeInstanceTTL,
		etcdCli:       cli,
		eventBus:      eventBus,
		sub:           bus.NewSub(),
		closed:        make(chan bool),
	}, nil
}

func (n *PinkNode) validate() error {

	if n.id == "" {
		return errors.New("the pink node id is nil error")
	}
	log.Infof("the pink node instance is using [id:%s,election_path:%s,election_ttl:%d,instance_path:%s,instance_ttl:%d]", n.id, n.electionPath, n.electionTTL, n.instancePath, n.instanceTTL)
	if n.etcdCli == nil {
		return errors.New("the etcd cli is nil error")
	}
	return nil
}

// self register
func (n *PinkNode) Run() error {

	err := n.validate()
	if err != nil {
		return err
	}

	log.Infof("the pink node instance at %s run success", n.id)
	// elect loop
	go n.electLoop()
	n.tryElect()
	go n.lookup()
	return nil
}

func (n *PinkNode) Stop() {
	n.closed <- true
}

// start elect
func (n *PinkNode) electLoop() {
	ticker := time.Tick(time.Second * time.Duration(n.electionTTL))
	log.Infof("the pink node instance %s start elect loop....", n.id)
	for {
		select {
		case <-ticker:
			if n.electionState == ElectionReadyState {
				n.electionState = ElectionDoingState
				log.Warnf("the pink node instance %s start try elect loop....", n.id)
				n.tryElect()
			}

		}
	}
}

// try  elect
func (n *PinkNode) tryElect() {

	defer func() {
		n.electionState = ElectionReadyState
	}()
	ctx, _ := context.WithTimeout(context.TODO(), time.Second*3)
	id, err := n.etcdCli.Leader(ctx, n.electionPath)
	if err == nil {
		log.Infof("the pink node instance %s has leader is %s", n.id, id)
		if id == n.id {
			n.NotifyState(protocol.Leader)
		} else {
			n.NotifyState(protocol.Follower)
		}
		return
	}
	n.NotifyState(protocol.Follower)
	log.Infof("the pink node instance %s find leader fail:%+v", n.id, err)
	if !errors.Is(err, concurrency.ErrElectionNoLeader) {
		log.Warnf("the pink node %s get leader fail:%+v", n.id, err)
		return
	}

	log.Infof("the pink node instance %s start campaign  leader", n.id)
	err = n.etcdCli.Campaign(ctx, n.id, n.electionPath, n.electionTTL)
	if err == nil {
		log.Infof("the pink node instance %s campaign  leader success", n.id)
		n.NotifyState(protocol.Leader)
		return
	}

}

func (n *PinkNode) selfRegister(instance string, leaseId int64) int64 {
RETRY:

	log.Infof("the pink node instance %s self register to:%s", n.id, instance)
	if leaseId > 0 {
		_ = n.etcdCli.Revoke(context.Background(), leaseId)
	}
	leaseId, err := n.etcdCli.Keepalive(context.Background(), instance, n.id, n.instanceTTL)
	if err != nil {
		time.Sleep(time.Second)
		goto RETRY
	}
	log.Infof("the pink node instance %s self register to:%s ,leaseId %d success", n.id, instance, leaseId)
	return leaseId
}

func (n *PinkNode) lookup() {
	instance := fmt.Sprintf("%s%s", n.instancePath, n.id)
	leaseId := int64(0)
	leaseId = n.selfRegister(instance, leaseId)
	response := n.etcdCli.Watch(instance)
	log.Infof("the pink node instance %s self register watch to:%s", n.id, instance)
	for {
		select {
		case event := <-response.KeyChangeCh:

			switch event.Event {
			case etcd.KeyDeleteChangeEvent:
				log.Warnf("the pink node instance %s self register watch  to:%s key delete event :%+v", n.id, instance, event)
				leaseId = n.selfRegister(instance, leaseId)
			}
		}
	}
}

// notify state
func (n *PinkNode) NotifyState(state int) {
	if state == n.state {
		return
	}
	n.state = state
	err := n.eventBus.Publish(protocol.NodeStateChangeTopic, state)
	if err != nil {
		log.Errorf("the pink node publish the event bus fail, state %d,", state)
	}

	log.Infof("the pink node publish the event bus success, state %d,", state)
}
