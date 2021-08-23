package job

import (
	"log"
	"math/rand"
	"time"
)

type UserSignInJob struct {
}

func (f *UserSignInJob) Target() string {
	return "com.busgo.user.job.UserSignInJob"
}

// execute job
func (f *UserSignInJob) Execute(param string) (string, error) {

	log.Printf("UserSignInJob receive param:%s", param)
	n := rand.Int63() % 20
	time.Sleep(time.Second * time.Duration(n))
	return "调用成功", nil
}
