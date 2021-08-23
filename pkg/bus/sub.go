package bus

type Sub struct {
	out chan interface{}
}

func NewSub() *Sub {
	return &Sub{out: make(chan interface{}, 10)}
}

func (sub *Sub) receive(msg interface{}) {
	sub.out <- msg
}

func (sub *Sub) Out() interface{} {
	return <-sub.out
}

func (sub *Sub) Receive() <-chan interface{} {
	return sub.out
}
