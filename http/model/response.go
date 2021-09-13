package model

import "encoding/json"

// 调度计划详情
type SchedulePlanDetails struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Group      string `json:"group"`
	Cron       string `json:"cron"`
	Target     string `json:"target"`
	Param      string `json:"param"`
	BeforeTime string `json:"before_time"`
	NextTime   string `json:"next_time"`
	Mobile     string `json:"mobile"`
	Remark     string `json:"remark"`
}

type JobConfDetails struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Group      string `json:"group"`
	Cron       string `json:"cron"`
	Target     string `json:"target"`
	Param      string `json:"param"`
	State      int32  `json:"state"`
	Mobile     string `json:"mobile"`
	Remark     string `json:"remark"`
	CreateTime string `json:"create_time"`
	UpdateTime string `json:"update_time"`
	Version    int64  `json:"version"`
}

func (d *JobConfDetails) Encode() string {
	content, _ := json.Marshal(d)
	return string(content)
}
func (d *JobConfDetails) Decode(content string) *JobConfDetails {
	_ = json.Unmarshal([]byte(content), d)
	return d
}

type GroupDetails struct {
	Name    string   `json:"name"`
	Remark  string   `json:"remark"`
	Clients []string `json:"clients"`
}

func (g *GroupDetails) Encode() string {
	content, _ := json.Marshal(g)
	return string(content)
}
func (g *GroupDetails) Decode(content string) *GroupDetails {
	_ = json.Unmarshal([]byte(content), g)
	return g
}

type NodeDetails struct {
	Id    string `json:"id"`
	State int32  `json:"state"`
}

type ClientDetails struct {
	Ip    string `json:"ip"`
	Path  string `json:"path"`
	Group string `json:"group"`
}

type ScheduleSnapshotDetails struct {
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

func (s *ScheduleSnapshotDetails) Decode(content string) *ScheduleSnapshotDetails {
	_ = json.Unmarshal([]byte(content), s)
	return s
}

type ExecuteSnapshotDetails struct {
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

func (e *ExecuteSnapshotDetails) Decode(content string) *ExecuteSnapshotDetails {
	_ = json.Unmarshal([]byte(content), e)
	return e
}

type ExecuteSnapshotHisDetails struct {
	Id           int64  `json:"id"`
	SnapshotId   string `json:"snapshot_id"`
	JobId        string `json:"job_id"`
	JobName      string `json:"job_name"`
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
	Remark       string `json:"remark"`
}
