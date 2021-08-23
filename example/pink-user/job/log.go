package job

import (
	"log"
	"math/rand"
	"time"
)

type UserLogJob struct {
}

func (f *UserLogJob) Target() string {
	return "com.busgo.user.job.UserLogJob"
}

// execute job
func (f *UserLogJob) Execute(param string) (string, error) {

	log.Printf("UserLogJob receive param:%s", param)
	n := rand.Int63() % 15
	time.Sleep(time.Second * time.Duration(n))
	return "调用成功", nil
}
