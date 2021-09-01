package handler

import (
	"context"
	"fmt"
	"github.com/busgo/pink/http/check"
	"github.com/busgo/pink/http/model"
	"github.com/busgo/pink/pkg/etcd"
	"github.com/busgo/pink/pkg/log"
	"github.com/busgo/pink/pkg/protocol"
	"github.com/busgo/pink/pkg/protocol/builder"
	"github.com/busgo/pink/pkg/util"
	"github.com/busgo/pink/schedule"
	"github.com/labstack/echo/v4"
	"math/rand"
	"strings"
	"time"
)

// PinkWebHandler  web api
type PinkWebHandler struct {
	etcdCli   *etcd.Cli
	scheduler *schedule.PinkScheduler
}

func NewPinkWebHandler(etcdCli *etcd.Cli, scheduler *schedule.PinkScheduler) *PinkWebHandler {
	return &PinkWebHandler{etcdCli: etcdCli, scheduler: scheduler}
}

// add job conf
func (h *PinkWebHandler) AddJobConf(c echo.Context) error {
	req := new(model.JobConfAddRequest)
	_ = c.Bind(req)
	err := check.JobConfAddRequest(req)
	if err != nil {
		return WriteParamError(c, err.Error())
	}

	jobConf := builder.NewJobConf(req)
	key := fmt.Sprintf("%s%s/%s", protocol.JobConfPath, req.Group, jobConf.Id)
	err = h.etcdCli.Put(c.Request().Context(), key, jobConf.Encode())
	if err != nil {
		return WriteBusinessError(c, err.Error())
	}
	return WriteOK(c, true)
}

// job conf update
func (h *PinkWebHandler) UpdateJobConf(c echo.Context) error {

	req := new(model.JobConfUpdateRequest)
	_ = c.Bind(req)

	err := check.JobConfUpdateRequest(req)
	if err != nil {
		return WriteParamError(c, err.Error())
	}
	ctx := c.Request().Context()
	key := fmt.Sprintf("%s%s/%s", protocol.JobConfPath, req.Group, req.Id)
	v, err := h.etcdCli.Get(ctx, key)
	if err != nil {
		return WriteBusinessError(c, err.Error())
	}
	if v == "" {
		return WriteParamError(c, fmt.Sprintf("the job conf group %s,id %s not found", req.Group, req.Id))
	}
	jobConf := new(protocol.JobConf).Decode(v)
	jobConf.Name = req.Name
	jobConf.Cron = req.Cron
	jobConf.Target = req.Target
	jobConf.Param = req.Param
	jobConf.Mobile = req.Mobile
	jobConf.Remark = req.Remark
	jobConf.State = req.State
	jobConf.UpdateTime = time.Now().Unix()
	jobConf.Version = jobConf.Version + 1
	content := jobConf.Encode()

	err = h.etcdCli.Put(ctx, key, content)
	if err != nil {
		return WriteBusinessError(c, err.Error())
	}
	return WriteOK(c, true)
}

// delete a job conf
func (h *PinkWebHandler) DeleteJobConf(c echo.Context) error {

	req := new(model.JobConfDeleteRequest)
	_ = c.Bind(req)
	err := check.JobConfDeleteRequest(req)
	if err != nil {
		return WriteParamError(c, err.Error())
	}
	key := fmt.Sprintf("%s%s/%s", protocol.JobConfPath, req.Group, req.Id)
	ctx := c.Request().Context()
	_, err = h.etcdCli.Get(ctx, key)
	if err != nil {
		return WriteBusinessError(c, err.Error())
	}

	err = h.etcdCli.Delete(ctx, key)
	if err != nil {
		return WriteBusinessError(c, err.Error())
	}

	return WriteOK(c, true)
}

// execute job
func (h *PinkWebHandler) JobExecute(c echo.Context) error {

	req := new(model.JobExecuteRequest)
	_ = c.Bind(req)
	if strings.TrimSpace(req.Group) == "" {
		return WriteParamError(c, "group is nil")
	}

	if strings.TrimSpace(req.Id) == "" {
		return WriteParamError(c, "id is nil")
	}

	ctx := c.Request().Context()
	key := fmt.Sprintf("%s%s/%s", protocol.JobConfPath, req.Group, req.Id)
	v, err := h.etcdCli.Get(ctx, key)
	if err != nil {
		return WriteBusinessError(c, err.Error())
	}

	if v == "" {
		return WriteBusinessError(c, "the job conf not found")
	}
	jobConf := new(protocol.JobConf).Decode(v)

	if jobConf.Group != req.Group {
		return WriteParamError(c, "the bad group")
	}
	clientInstancePath := fmt.Sprintf(protocol.ClientInstancePath, req.Group)
	_, values, err := h.etcdCli.GetWithPrefix(ctx, clientInstancePath)

	if err != nil {
		return WriteBusinessError(c, err.Error())
	}
	if len(values) == 0 {
		return WriteBusinessError(c, "has no client to dispatch the task")
	}
	ip := values[rand.Intn(len(values))]

	now := time.Now().Format("2006-01-02 15:04:05")
	executeSnapshot := &protocol.ExecuteSnapshot{
		Id:           util.Generate(),
		JobId:        jobConf.Id,
		Name:         jobConf.Name,
		Group:        jobConf.Group,
		Cron:         jobConf.Cron,
		Target:       jobConf.Target,
		Ip:           ip,
		Param:        jobConf.Param,
		State:        protocol.ExecuteSnapshotInit,
		BeforeTime:   now,
		ScheduleTime: now,
		StartTime:    "",
		EndTime:      "",
		Times:        0,
		Mobile:       jobConf.Mobile,
		Version:      0,
		Remark:       jobConf.Remark,
	}

	path := fmt.Sprintf(protocol.ExecuteSnapshotPath, req.Group, ip, executeSnapshot.Id)
	content := executeSnapshot.Encode()
	ctx, _ = context.WithTimeout(context.Background(), time.Second*3)
	err = h.etcdCli.PutWithNotExist(ctx, path, content)
	if err != nil {
		log.Errorf("create execute plan snapshot %s  for ip %s fail err %+v", executeSnapshot.Encode(), ip, err)
		return WriteBusinessError(c, err.Error())
	}
	log.Infof(" create execute plan snapshot %s  for ip %s success", content, ip)
	return WriteOK(c, fmt.Sprintf("dispatch the task to client:%s", ip))
}

// query job conf list
func (h *PinkWebHandler) JobConfList(c echo.Context) error {

	req := new(model.JobConfListRequest)
	_ = c.Bind(req)

	prefix := protocol.JobConfPath
	if strings.TrimSpace(req.Group) != "" {
		prefix = fmt.Sprintf("%s%s", protocol.JobConfPath, req.Group)
	}
	keys, values, err := h.etcdCli.GetWithPrefix(c.Request().Context(), prefix)
	if err != nil {
		return WriteBusinessError(c, err.Error())
	}

	details := make([]*model.JobConfDetails, 0)
	if len(keys) == 0 || len(values) == 0 {
		return WriteOK(c, details)
	}

	for _, value := range values {
		jobConf := new(protocol.JobConf).Decode(value)
		details = append(details, &model.JobConfDetails{
			Id:         jobConf.Id,
			Name:       jobConf.Name,
			Group:      jobConf.Group,
			Cron:       jobConf.Cron,
			Target:     jobConf.Target,
			Param:      jobConf.Param,
			State:      jobConf.State,
			Mobile:     jobConf.Mobile,
			Remark:     jobConf.Remark,
			CreateTime: time.Unix(jobConf.CreateTime, 0).Format("2006-01-02 15:04:05"),
			UpdateTime: time.Unix(jobConf.UpdateTime, 0).Format("2006-01-02 15:04:05"),
			Version:    0,
		})
	}
	return WriteOK(c, details)
}

func (h *PinkWebHandler) AddGroup(c echo.Context) error {

	req := new(model.AddGroupRequest)

	_ = c.Bind(req)

	err := check.AddGroupRequest(req)
	if err != nil {
		return WriteParamError(c, err.Error())
	}

	ctx := c.Request().Context()
	key := fmt.Sprintf("%s%s", protocol.GroupPath, req.Name)
	content, err := h.etcdCli.Get(ctx, key)
	if err != nil {
		return WriteBusinessError(c, err.Error())
	}
	if content != "" {
		return WriteBusinessError(c, "集群已存在")
	}
	details := &model.GroupDetails{
		Name:   req.Name,
		Remark: req.Remark,
	}
	content = details.Encode()
	err = h.etcdCli.PutWithNotExist(ctx, key, content)
	if err != nil {
		return WriteBusinessError(c, err.Error())
	}
	return WriteOK(c, true)
}

// query group list
func (h *PinkWebHandler) GroupList(c echo.Context) error {

	ctx := c.Request().Context()

	_, values, err := h.etcdCli.GetWithPrefix(ctx, protocol.GroupPath)
	if err != nil {
		return WriteBusinessError(c, err.Error())
	}

	groups := make([]*model.GroupDetails, 0)
	if len(values) == 0 {
		return WriteOK(c, groups)
	}
	for _, content := range values {
		group := new(model.GroupDetails).Decode(content)
		groups = append(groups, group)
	}
	return WriteOK(c, groups)

}

// query group details list
func (h *PinkWebHandler) GroupDetailsList(c echo.Context) error {

	ctx := c.Request().Context()

	_, values, err := h.etcdCli.GetWithPrefix(ctx, protocol.GroupPath)
	if err != nil {
		return WriteBusinessError(c, err.Error())
	}

	groups := make([]*model.GroupDetails, 0)
	if len(values) == 0 {
		return WriteOK(c, groups)
	}
	for _, content := range values {
		group := new(model.GroupDetails).Decode(content)

		clientInstancePath := fmt.Sprintf(protocol.ClientInstancePath, group.Name)
		_, values, err := h.etcdCli.GetWithPrefix(ctx, clientInstancePath)

		if err != nil {
			return WriteBusinessError(c, err.Error())
		}

		if len(values) > 0 {
			group.Clients = values
		}

		groups = append(groups, group)
	}
	return WriteOK(c, groups)

}

// node list
func (h *PinkWebHandler) NodeList(c echo.Context) error {

	ctx := c.Request().Context()
	_, values, err := h.etcdCli.GetWithPrefix(ctx, protocol.NodeInstancePath)
	if err != nil {
		return WriteBusinessError(c, err.Error())
	}

	details := make([]*model.NodeDetails, 0)
	if len(values) == 0 {
		return WriteOK(c, details)
	}

	id, err := h.etcdCli.Leader(ctx, protocol.NodeElectionPath)
	if err != nil {
		return WriteBusinessError(c, err.Error())
	}

	for _, node := range values {

		var state int32 = protocol.Follower
		if id == node {
			state = protocol.Leader
		}

		details = append(details, &model.NodeDetails{
			Id:    node,
			State: state,
		})

	}
	return WriteOK(c, details)
}

// schedule plan list
func (h *PinkWebHandler) SchedulePlanList(c echo.Context) error {

	plans := h.scheduler.GetAllSchedulePlan()

	scheduleDetails := make([]*model.SchedulePlanDetails, 0)
	if len(plans) == 0 {
		return WriteOK(c, scheduleDetails)
	}

	for _, plan := range plans {

		scheduleDetails = append(scheduleDetails, &model.SchedulePlanDetails{
			Id:         plan.Id,
			Name:       plan.Name,
			Group:      plan.Group,
			Cron:       plan.Cron,
			Target:     plan.Target,
			Param:      plan.Param,
			BeforeTime: plan.BeforeTime.Format("2006-01-02 15:04:05"),
			NextTime:   plan.NextTime.Format("2006-01-02 15:04:05"),
			Mobile:     plan.Mobile,
			Remark:     plan.Remark,
		})

	}
	return WriteOK(c, scheduleDetails)
}

// schedule clients
func (h *PinkWebHandler) SchedulePlanClients(c echo.Context) error {

	req := new(model.ClientsRequest)
	_ = c.Bind(req)
	if strings.TrimSpace(req.Group) == "" {
		return WriteParamError(c, "group is nil")
	}
	prefix := fmt.Sprintf(protocol.ClientInstancePath, req.Group)
	keys, values, err := h.etcdCli.GetWithPrefix(c.Request().Context(), prefix)
	if err != nil {
		return WriteBusinessError(c, err.Error())
	}
	details := make([]*model.ClientDetails, 0)
	if len(keys) == 0 {
		return WriteOK(c, details)
	}

	for pos, key := range keys {
		details = append(details, &model.ClientDetails{
			Ip:    values[pos],
			Path:  key,
			Group: req.Group,
		})
	}
	return WriteOK(c, details)
}

// schedule snapshots
func (h *PinkWebHandler) ScheduleSnapshots(c echo.Context) error {

	_, values, err := h.etcdCli.GetWithPrefix(c.Request().Context(), protocol.ScheduleSnapshotPath)
	if err != nil {
		return WriteBusinessError(c, err.Error())
	}

	details := make([]*model.ScheduleSnapshotDetails, 0)
	if len(values) == 0 {
		return WriteOK(c, details)
	}

	for _, value := range values {
		details = append(details, new(model.ScheduleSnapshotDetails).Decode(value))
	}
	return WriteOK(c, details)
}

// execute snapshot
func (h *PinkWebHandler) ExecuteSnapshots(c echo.Context) error {

	req := new(model.ExecuteSnapshotsRequest)
	_ = c.Bind(req)

	if strings.TrimSpace(req.Group) == "" {
		return WriteParamError(c, "group is nil")
	}
	if strings.TrimSpace(req.Id) != "" && strings.TrimSpace(req.Ip) == "" {
		return WriteParamError(c, "client is nil")
	}

	prefix := fmt.Sprintf("%s%s", protocol.ExecuteSnapshotBasePath, req.Group)
	if strings.TrimSpace(req.Ip) != "" && strings.TrimSpace(req.Id) == "" {
		prefix = fmt.Sprintf("%s/%s", prefix, req.Ip)
	} else if strings.TrimSpace(req.Ip) != "" && strings.TrimSpace(req.Id) != "" {
		prefix = fmt.Sprintf(protocol.ExecuteSnapshotPath, req.Group, req.Ip, req.Id)
	}
	_, values, err := h.etcdCli.GetWithPrefix(c.Request().Context(), prefix)
	if err != nil {
		return WriteBusinessError(c, err.Error())
	}
	details := make([]*model.ExecuteSnapshotDetails, 0)
	if len(values) == 0 {
		return WriteOK(c, details)
	}

	for _, value := range values {
		details = append(details, new(model.ExecuteSnapshotDetails).Decode(value))
	}
	return WriteOK(c, details)
}

// delete the schedule snapshot
func (h *PinkWebHandler) ScheduleSnapshotDelete(c echo.Context) error {

	req := new(model.ScheduleSnapshotDeleteRequest)

	_ = c.Bind(req)
	if strings.TrimSpace(req.Id) == "" {
		return WriteParamError(c, "id is nil")
	}

	key := fmt.Sprintf("%s%s", protocol.ScheduleSnapshotPath, req.Id)
	ctx := c.Request().Context()
	v, err := h.etcdCli.Get(ctx, key)
	if err != nil {
		return WriteBusinessError(c, err.Error())
	}

	if strings.TrimSpace(v) == "" {
		return WriteBusinessError(c, "the snapshot not found")
	}

	err = h.etcdCli.Delete(ctx, key)
	if err != nil {
		return WriteBusinessError(c, err.Error())
	}

	return WriteOK(c, true)

}

/**
 * execute history snapshot
 */
func (h *PinkWebHandler) ExecuteHistorySnapshots(c echo.Context) error {
	req := new(model.ExecuteSnapshotsRequest)
	_ = c.Bind(req)
	prefix := protocol.ExecuteSnapshotHistoryPath
	if strings.TrimSpace(req.Id) != "" {
		prefix = fmt.Sprintf("%s%s", protocol.ExecuteSnapshotHistoryPath, req.Id)
	}
	_, values, err := h.etcdCli.GetWithPrefix(c.Request().Context(), prefix)
	if err != nil {
		return WriteBusinessError(c, err.Error())
	}
	details := make([]*model.ExecuteSnapshotDetails, 0)
	if len(values) == 0 {
		return WriteOK(c, details)
	}

	for _, value := range values {
		details = append(details, new(model.ExecuteSnapshotDetails).Decode(value))
	}
	return WriteOK(c, details)
}

// delete execute history snapshot
func (h *PinkWebHandler) ExecuteHistorySnapshotDelete(c echo.Context) error {

	req := new(model.ExecuteHistorySnapshotDeleteRequest)
	_ = c.Bind(req)
	if strings.TrimSpace(req.Id) == "" {
		return WriteParamError(c, "id is nil")
	}

	err := h.etcdCli.Delete(c.Request().Context(), fmt.Sprintf("%s%s", protocol.ExecuteSnapshotHistoryPath, req.Id))
	if err != nil {
		return WriteBusinessError(c, err.Error())
	}

	return WriteOK(c, true)

}
