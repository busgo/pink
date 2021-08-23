package balance

import (
	"fmt"
	"sync/atomic"
)

type PinkBalancer interface {
	// balance
	Balance(clients []string) (string, error)
}
type PinkBalanceManaged struct {
	balancers map[string]PinkBalancer
}

var bm = PinkBalanceManaged{balancers: make(map[string]PinkBalancer)}

func (bm *PinkBalanceManaged) register(name string, balancer PinkBalancer) {
	bm.balancers[name] = balancer
}

func RegisterBalancer(name string, balancer PinkBalancer) {
	bm.register(name, balancer)
}

// balance
func Balance(name string, clients []string) (string, error) {
	balancer := bm.balancers[name]
	if balancer == nil {
		return "", fmt.Errorf("not found the balancer with name %s", name)
	}
	return balancer.Balance(clients)
}

// round_robin balancer
type RoundRobinBalancer struct {
	pos int32
}

func init() {
	RegisterBalancer("round_robin", &RoundRobinBalancer{pos: 0})
}

func (rr *RoundRobinBalancer) Balance(clients []string) (string, error) {

	client := rr.getClient(clients)
	if client == "" {
		return "", fmt.Errorf("the round_robin balancer the clients is nil")
	}
	return client, nil

}

func (rr *RoundRobinBalancer) getClient(clients []string) string {

	size := len(clients)
	if size == 0 {
		return ""
	}

	p := atomic.AddInt32(&rr.pos, 1)
	if p >= int32(size) {
		p = 0
		atomic.StoreInt32(&rr.pos, 0)
	}
	return clients[p]
}
