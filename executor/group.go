package executor

import (
	"context"
	"fmt"
	"github.com/busgo/pink/pkg/balance"
	"github.com/busgo/pink/pkg/etcd"
	"github.com/busgo/pink/pkg/protocol"
	"log"
	"strings"
	"sync"
	"time"
)

type PinkGroupManaged struct {
	cli    *etcd.Cli
	groups map[string]*PinkGroup
	sync.RWMutex
}

func NewPinkGroupManaged(cli *etcd.Cli) *PinkGroupManaged {

	managed := &PinkGroupManaged{
		cli:     cli,
		groups:  make(map[string]*PinkGroup),
		RWMutex: sync.RWMutex{},
	}
	go managed.lookup()
	return managed
}

// lookup the group change
func (managed *PinkGroupManaged) lookup() {

	keys, _, err := managed.cli.GetWithPrefix(context.Background(), protocol.GroupPath)
	if err != nil {
		log.Panicf("the pink group managed get group list fail:%+v", err)
	}

	for _, key := range keys {
		managed.addPinkGroup(getSuffixName(key), key)
	}
	log.Printf("the pink group managed watch for %s", protocol.GroupPath)
	response := managed.cli.WatchWithPrefix(protocol.GroupPath)
	for {
		select {
		case ch := <-response.KeyChangeCh:
			switch ch.Event {
			case etcd.KeyCreateChangeEvent:
				managed.addPinkGroup(getSuffixName(ch.Key), ch.Key)
			case etcd.KeyUpdateChangeEvent:
				managed.addPinkGroup(getSuffixName(ch.Key), ch.Key)
			case etcd.KeyDeleteChangeEvent:
				managed.deletePinkGroup(getSuffixName(ch.Key))
			}
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

// delete the pink group with name
func (managed *PinkGroupManaged) deletePinkGroup(name string) {
	group := managed.getPinkGroup(name)
	if group == nil {
		return
	}
	managed.Lock()
	group.close()
	delete(managed.groups, name)
	managed.Unlock()
	log.Printf("the pink group managed delete group:%s success", group.name)
}

// add group with name
func (managed *PinkGroupManaged) addPinkGroup(name, path string) {

	group := managed.getPinkGroup(name)
	if group == nil {
		group = NewPinkGroup(name, path, managed.cli)
		managed.Lock()
		managed.groups[name] = group
		managed.Unlock()
		log.Printf("the pink group managed add group:%s success", group.name)
	}

	group.fetch()

}

// get group with name
func (managed *PinkGroupManaged) getPinkGroup(name string) *PinkGroup {
	managed.RLock()
	defer managed.RUnlock()
	return managed.groups[name]
}

func (managed *PinkGroupManaged) UnParkPinkClient(name string) (string, error) {

	group := managed.getPinkGroup(name)
	if group == nil {
		return "", fmt.Errorf("the managed group not found pink group %s", name)
	}

	return group.unParkPinkClient()

}

type PinkGroup struct {
	name    string
	path    string
	cli     *etcd.Cli
	clients []*PinkClient
	sync.RWMutex
	closeCh chan bool
}

// new a pink group
func NewPinkGroup(name, path string, cli *etcd.Cli) *PinkGroup {

	g := &PinkGroup{
		name:    name,
		path:    path,
		cli:     cli,
		clients: make([]*PinkClient, 0),
		RWMutex: sync.RWMutex{},
		closeCh: make(chan bool),
	}

	go g.lookup()
	return g
}

func (g *PinkGroup) fetch() {
	path := fmt.Sprintf(protocol.ClientInstancePath, g.name)

	retries := 0
RETRY:
	log.Printf("the pink group %s start fetch clients for %s ...", g.name, path)
	keys, values, err := g.cli.GetWithPrefix(context.Background(), path)

	if err != nil {
		retries++
		if retries < 3 {
			time.Sleep(time.Second)
			goto RETRY
		}
	}

	if len(keys) == 0 {
		log.Printf("the pink group %s has no client for %s ...", g.name, path)
		return
	}

	for pos, path := range keys {
		name := values[pos]
		g.addPinkClient(name, path)
	}
	log.Printf("the pink group %s finish fetch clients for %s...", g.name, path)
}
func (g *PinkGroup) close() {
	g.closeCh <- true
}

// lookup the clients change
func (g *PinkGroup) lookup() {
	path := fmt.Sprintf(protocol.ClientInstancePath, g.name)
	log.Printf("the pink group %s watch for %s...", g.name, path)
	response := g.cli.WatchWithPrefix(path)

	for {
		select {
		case ch := <-response.KeyChangeCh:
			switch ch.Event {
			case etcd.KeyCreateChangeEvent:
				g.addPinkClient(ch.Value, ch.Key)
			case etcd.KeyUpdateChangeEvent:
				g.addPinkClient(ch.Value, ch.Key)
			case etcd.KeyDeleteChangeEvent:
				g.deletePinkClient(getSuffixName(ch.Key))

			}
		case <-g.closeCh:
			log.Printf("the pink group %s closed...", g.name)
			return
		}
	}
}

// add pink client with name
func (g *PinkGroup) addPinkClient(name, path string) {
	client, pos := g.getPinkClient(name)
	if pos >= 0 {
		log.Printf("the pink group %s has no  pink client %s ,after create", g.name, client.name)
		return
	}
	client = NewPinkClient(name, path)
	g.Lock()
	g.clients = append(g.clients, client)
	g.Unlock()
	log.Printf("the pink group %s add pink client %s success", g.name, name)
}

func (g *PinkGroup) deletePinkClient(name string) {
	_, pos := g.getPinkClient(name)
	if pos == -1 {
		log.Printf("the pink group %s not found pink client for name %s", g.name, name)
		return
	}
	g.Lock()
	g.clients = append(g.clients[:pos], g.clients[pos+1:]...)
	g.Unlock()
	log.Printf("the pink group %s has  delete pink client  %s", g.name, name)

}

// un park a pink client
func (g *PinkGroup) unParkPinkClient() (string, error) {
	g.RLock()
	defer g.RUnlock()

	if len(g.clients) == 0 {
		return "", fmt.Errorf("the pink group %s has no client", g.name)
	}

	ips := make([]string, len(g.clients))

	for pos, client := range g.clients {
		ips[pos] = client.name
	}
	return balance.Balance("round_robin", ips)
}

// get pink client with name
func (g *PinkGroup) getPinkClient(name string) (*PinkClient, int) {

	g.RLock()
	defer g.RUnlock()
	if len(g.clients) == 0 {
		log.Printf("the pink group %s has no pink client %s", g.name, name)
		return nil, -1
	}

	for pos, c := range g.clients {

		if c.name == name {
			return c, pos
		}
	}
	return nil, -1
}

type PinkClient struct {
	name string
	path string
}

func NewPinkClient(name, path string) *PinkClient {

	return &PinkClient{
		name: name,
		path: path,
	}
}
