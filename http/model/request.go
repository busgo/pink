package model

type Result struct {
	Code    int32       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type JobConfUpdateRequest struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Group  string `json:"group"`
	Cron   string `json:"cron"`
	Target string `json:"target"`
	Param  string `json:"param"`
	State  int32  `json:"state"`
	Mobile string `json:"mobile"`
	Remark string `json:"remark"`
}

type JobConfListRequest struct {
	Group string `json:"group"`
}

type JobConfAddRequest struct {
	Name   string `json:"name"`
	Group  string `json:"group"`
	Cron   string `json:"cron"`
	Target string `json:"target"`
	Param  string `json:"param"`
	State  int32  `json:"state"`
	Mobile string `json:"mobile"`
	Remark string `json:"remark"`
}

type JobExecuteRequest struct {
	Group string `json:"group"`
	Id    string `json:"id"`
}

type JobConfDeleteRequest struct {
	Id    string `json:"id"`
	Group string `json:"group"`
}

type AddGroupRequest struct {
	Name   string `json:"name"`
	Remark string `json:"remark"`
}

type ClientsRequest struct {
	Group string `json:"group"`
}

type ExecuteSnapshotsRequest struct {
	Id        string `json:"id"`
	JobId     string `json:"job_id"`
	Group     string `json:"group"`
	Ip        string `json:"ip"`
	State     string `json:"state"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type ScheduleSnapshotDeleteRequest struct {
	Id string `json:"id"`
}

type ExecuteHistorySnapshotsRequest struct {
	JobId             string `json:"job_id"`
	SnapshotId        string `json:"snapshot_id"`
	JobName           string `json:"job_name"`
	Group             string `json:"group"`
	Ip                string `json:"ip"`
	State             int32  `json:"state"`
	ScheduleStartTime string `json:"schedule_start_time"`
	ScheduleEndTime   string `json:"schedule_end_time"`
	PageNo            int32  `json:"page_no"`
	PageSize          int32  `json:"page_size"`
}

type ExecuteHistorySnapshotDeleteRequest struct {
	Id int64 `json:"id"`
}
