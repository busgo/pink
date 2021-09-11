package repository

import (
	"github.com/busgo/pink/db"
	"github.com/busgo/pink/db/model"
	"github.com/jmoiron/sqlx"
)

const (
	InsertExecuteSnapshotSql = "INSERT INTO execute_snapshot_his(`snapshot_id`,`job_id`,`job_name`, `group`,`cron`,`target`,`ip`,`param`,`state`,`before_time`,`schedule_time`,`end_time`,`times`,`mobile`,`remark`) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
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
		snapshot.EndTime,
		snapshot.Times,
		snapshot.Mobile,
		snapshot.Remark,
	)

	return err
}

func (repo *ExecuteSnapshotHisRepository) SearchExecuteSnapshotHisByPage() {

	//repo.connection.NamedQuery()
}
