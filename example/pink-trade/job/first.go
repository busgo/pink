package job

import (
	"log"
	"math/rand"
	"time"
)

type FirstTradeJob struct {
}

func (f *FirstTradeJob) Target() string {
	return "com.busgo.trade.job.FirstTradeJob"
}

// execute job
func (f *FirstTradeJob) Execute(param string) (string, error) {

	log.Printf("FirstTradeJob receive param:%s", param)
	n := rand.Int63() % 20
	time.Sleep(time.Second * time.Duration(n))
	return "调用成功", nil
}
