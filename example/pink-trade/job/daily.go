package job

import (
	"log"
	"math/rand"
	"time"
)

type DailyTradeJob struct {
}

func (f *DailyTradeJob) Target() string {
	return "com.busgo.trade.job.DailyTradeJob"
}

// execute job
func (f *DailyTradeJob) Execute(param string) (string, error) {

	log.Printf("DailyTradeJob receive param:%s", param)
	n := rand.Int63() % 20
	time.Sleep(time.Second * time.Duration(n))
	return "调用成功", nil
}
