package bus

import "sync"

type SubNode struct {
	subs []*Sub
	rw   sync.RWMutex
}

func NewSubNode() *SubNode {

	return &SubNode{
		subs: make([]*Sub, 0),
		rw:   sync.RWMutex{},
	}
}

func (n *SubNode) Len() int {
	return len(n.subs)
}

func (n *SubNode) delete(sub *Sub) {

	n.rw.Lock()
	defer n.rw.Unlock()

	pos := n.findSubPos(sub)
	if pos == -1 {
		return
	}
	n.subs = append(n.subs[:pos], n.subs[pos+1:]...)
}

func (n *SubNode) findSubPos(sub *Sub) int {

	if len(n.subs) == 0 {
		return -1
	}
	for pos, s := range n.subs {
		if s == sub {
			return pos
		}
	}

	return -1
}
