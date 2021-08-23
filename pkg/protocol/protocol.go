package protocol

import (
	"encoding/json"
	"time"
)
import "github.com/robfig/cron"

const (
	NodeStateChangeTopic = "NodeStateChangeTopic"
	NodeElectionTTL      = 3
	NodeInstanceTTL      = 10

	NodeElectionPath = "/pink/node/election/"
	NodeInstancePath = "/pink/node/instances/"

	JobConfPath                = "/pink/job/conf/"
	GroupPath                  = "/pink/group/instances/"
	ClientInstancePath         = "/pink/client/%s/instances/"
	ScheduleSnapshotPath       = "/pink/schedule/snapshots/"
	ExecuteSnapshotBasePath    = "/pink/execute/snapshots/"
	ExecuteSnapshotPath        = "/pink/execute/snapshots/%s/%s/%s"
	ExecuteSnapshotHistoryPath = "/pink/execute/history/snapshots/"
)

type JobChangeEventType int32

const (
	JobCreateChangeEventType JobChangeEventType = iota + 1
	JobUpdateChangeEventType
	JobDeleteChangeEventType
)

type JobState int32

const (
	JobNormalState JobState = iota + 1
	JobStopState
)

const (
	Follower = 1
	Leader   = 2
)

const (
	ExecuteSnapshotInit int32 = iota
	ExecuteSnapshotDoing
	ExecuteSnapshotSuccess
	ExecuteSnapshotFail
)

type JobConf struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Group      string `json:"group"`
	Cron       string `json:"cron"`
	Target     string `json:"target"`
	Param      string `json:"param"`
	State      int32  `json:"state"`
	Mobile     string `json:"mobile"`
	Remark     string `json:"remark"`
	Version    int32  `json:"version"`
	CreateTime int64  `json:"create_time"`
	UpdateTime int64  `json:"update_time"`
}

func (j *JobConf) Encode() string {
	content, _ := json.Marshal(j)
	return string(content)
}

func (j *JobConf) Decode(content string) *JobConf {
	_ = json.Unmarshal([]byte(content), j)
	return j
}

type JobChangeEvent struct {
	Event   JobChangeEventType
	Content *JobConf
}

type SchedulePlan struct {
	Id         string        `json:"id"`
	Name       string        `json:"name"`
	Group      string        `json:"group"`
	Cron       string        `json:"cron"`
	Target     string        `json:"target"`
	Param      string        `json:"param"`
	BeforeTime time.Time     `json:"before_time"`
	NextTime   time.Time     `json:"next_time"`
	Schedule   cron.Schedule `json:"-"`
	Mobile     string        `json:"mobile"`
	Version    int32         `json:"version"`
	Remark     string        `json:"remark"`
	CreateTime time.Time     `json:"create_time"`
	UpdateTime time.Time     `json:"update_time"`
}

type SchedulePlanSnapshot struct {
	Id           string `json:"id"`
	JobId        string `json:"job_id"`
	Name         string `json:"name"`
	Group        string `json:"group"`
	Cron         string `json:"cron"`
	Target       string `json:"target"`
	Param        string `json:"param"`
	BeforeTime   string `json:"before_time"`
	ScheduleTime string `json:"schedule_time"`
	Mobile       string `json:"mobile"`
	Version      int32  `json:"version"`
	Remark       string `json:"remark"`
}

func (snapshot *SchedulePlanSnapshot) Encode() string {
	content, _ := json.Marshal(snapshot)
	return string(content)
}

func (snapshot *SchedulePlanSnapshot) Decode(content string) *SchedulePlanSnapshot {
	_ = json.Unmarshal([]byte(content), snapshot)
	return snapshot
}

type ExecuteSnapshot struct {
	Id           string `json:"id"`
	JobId        string `json:"job_id"`
	Name         string `json:"name"`
	Group        string `json:"group"`
	Cron         string `json:"cron"`
	Target       string `json:"target"`
	Ip           string `json:"ip"`
	Param        string `json:"param"`
	State        int32  `json:"state"`
	BeforeTime   string `json:"before_time"`
	ScheduleTime string `json:"schedule_time"`
	StartTime    string `json:"start_time"`
	EndTime      string `json:"end_time"`
	Times        int64  `json:"times"`
	Mobile       string `json:"mobile"`
	Version      int32  `json:"version"`
	Remark       string `json:"remark"`
}

func (es *ExecuteSnapshot) Encode() string {

	content, _ := json.Marshal(es)
	return string(content)
}

func (es *ExecuteSnapshot) Decode(content string) *ExecuteSnapshot {

	_ = json.Unmarshal([]byte(content), es)
	return es
}
