package builder

import (
	"github.com/busgo/pink/http/model"
	"github.com/busgo/pink/pkg/protocol"
	"github.com/busgo/pink/pkg/util"
	"github.com/robfig/cron"
	"time"
)

func NewJobConf(request *model.JobConfAddRequest) *protocol.JobConf {

	return &protocol.JobConf{
		Id:         util.Generate(),
		Name:       request.Name,
		Group:      request.Group,
		Cron:       request.Cron,
		Target:     request.Target,
		Param:      request.Param,
		State:      request.State,
		Mobile:     request.Mobile,
		Remark:     request.Remark,
		Version:    1,
		CreateTime: time.Now().Unix(),
		UpdateTime: time.Now().Unix(),
	}
}

func NewSchedulePlan(j *protocol.JobConf, sch cron.Schedule) *protocol.SchedulePlan {

	return &protocol.SchedulePlan{
		Id:         j.Id,
		Name:       j.Name,
		Group:      j.Group,
		Cron:       j.Cron,
		Target:     j.Target,
		Param:      j.Param,
		NextTime:   sch.Next(time.Now()),
		Schedule:   sch,
		Mobile:     j.Mobile,
		Remark:     j.Remark,
		Version:    j.Version,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
}
func NewSchedulePlanSnapshot(plan *protocol.SchedulePlan, scheduleTime time.Time) *protocol.SchedulePlanSnapshot {

	return &protocol.SchedulePlanSnapshot{
		Id:           util.Generate(),
		JobId:        plan.Id,
		Name:         plan.Name,
		Group:        plan.Group,
		Cron:         plan.Cron,
		Target:       plan.Target,
		Param:        plan.Param,
		BeforeTime:   plan.BeforeTime.Format("2006-01-02 15:04:05"),
		ScheduleTime: scheduleTime.Format("2006-01-02 15:04:05"),
		Mobile:       plan.Mobile,
		Remark:       plan.Remark,
		Version:      plan.Version,
	}
}
