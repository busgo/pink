package util

import (
	"fmt"
	"os"
	"sync/atomic"
	"time"
)

var pos int64

func Generate() string {

	now := time.Now()
	str := now.Format("060102150405")
	m := now.UnixNano()/1e6 - now.UnixNano()/1e9*1e3
	ms := fmt.Sprintf("%03d", m)
	p := os.Getpid() % 1000
	pid := fmt.Sprintf("%03d", p)

	i := atomic.AddInt64(&pos, 1)
	r := i % 10000
	rs := fmt.Sprintf("%04d", r)
	return fmt.Sprintf("%s%s%s%s", str, ms, pid, rs)

}
