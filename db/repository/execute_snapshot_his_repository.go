package repository

import (
	"fmt"
	"github.com/busgo/pink/db"
	"github.com/busgo/pink/db/model"
	model2 "github.com/busgo/pink/http/model"
	"github.com/jmoiron/sqlx"
	"strings"
)

const (
	InsertExecuteSnapshotSql = "INSERT INTO execute_snapshot_his(`snapshot_id`,`job_id`,`job_name`, `group`,`cron`,`target`,`ip`,`param`,`state`,`before_time`,`schedule_time`,`start_time`,`end_time`,`times`,`mobile`,`remark`) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
)

// ExecuteSnapshotHisRepository
type ExecuteSnapshotHisRepository struct {
	connection *sqlx.DB
}

// new ExecuteSnapshotHisRepository
func NewExecuteSnapshotRepository(db *db.DB) *ExecuteSnapshotHisRepository {
	return &ExecuteSnapshotHisRepository{connection: db.DS}
}

func (repo *ExecuteSnapshotHisRepository) Insert(snapshot *model.ExecuteSnapshotHis) error {

	_, err := repo.connection.Exec(InsertExecuteSnapshotSql,
		snapshot.SnapshotId,
		snapshot.JobId,
		snapshot.JobName,
		snapshot.Group,
		snapshot.Cron,
		snapshot.Target,
		snapshot.Ip,
		snapshot.Param,
		snapshot.State,
		snapshot.BeforeTime,
		snapshot.ScheduleTime,
		snapshot.StartTime,
		snapshot.EndTime,
		snapshot.Times,
		snapshot.Mobile,
		snapshot.Remark,
	)

	return err
}

//
func (repo *ExecuteSnapshotHisRepository) SearchExecuteSnapshotHisByPage(request *model2.ExecuteHistorySnapshotsRequest) ([]*model.ExecuteSnapshotHis, error) {

	sql := "SELECT * FROM execute_snapshot_his WHERE 1=1 "

	parameters := make(map[string]interface{})
	if strings.TrimSpace(request.JobId) != "" {
		sql = sql + " AND job_id=:jobId "
		parameters["jobId"] = request.JobId
	}

	if strings.TrimSpace(request.SnapshotId) != "" {
		sql = sql + " AND snapshot_id=:snapshotId "
		parameters["snapshotId"] = request.SnapshotId
	}
	if strings.TrimSpace(request.Group) != "" {
		sql = sql + " AND `group`=:group "
		parameters["group"] = request.Group
	}
	if strings.TrimSpace(request.Ip) != "" {
		sql = sql + " AND ip=:ip "
		parameters["ip"] = request.Ip
	}

	if request.State != 0 {
		sql = sql + " AND `state`=:state "
		parameters["state"] = request.State
	}

	if strings.TrimSpace(request.JobName) != "" {
		sql = sql + " AND job_name like concat('%',:jobName,'%') "
		parameters["jobName"] = request.JobName
	}

	if strings.TrimSpace(request.ScheduleStartTime) != "" {
		sql = sql + " AND schedule_time >=:scheduleStartTime "
		parameters["scheduleStartTime"] = request.ScheduleStartTime
	}
	if strings.TrimSpace(request.ScheduleEndTime) != "" {
		sql = sql + " AND schedule_time <=:scheduleEndTime "
		parameters["scheduleEndTime"] = request.ScheduleEndTime
	}

	if request.PageNo <= 0 {
		request.PageNo = 1
	}

	if request.PageSize <= 0 || request.PageSize > 100 {
		request.PageSize = 10
	}
	sql = fmt.Sprintf("%s LIMIT %d,%d", sql, request.PageNo-1, request.PageSize)

	rows, err := repo.connection.NamedQuery(sql, parameters)
	if err != nil {
		return make([]*model.ExecuteSnapshotHis, 0), err
	}

	records := make([]*model.ExecuteSnapshotHis, 0)

	for rows.Next() {
		record := new(model.ExecuteSnapshotHis)
		err = rows.StructScan(record)
		if err != nil {
			return make([]*model.ExecuteSnapshotHis, 0), err
		}
		records = append(records, record)
	}

	return records, nil
}
