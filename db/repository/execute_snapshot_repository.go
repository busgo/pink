package repository

import (
	"github.com/busgo/pink/db"
	"github.com/busgo/pink/db/model"
	"github.com/jmoiron/sqlx"
)

const (
	InsertExecuteSnapshotSql = "INSERT INTO execute_snapshot (`snapshot_id`,`job_name`,`group`,`cron`,`target`,`ip`,`param`,`state`,`before_time`,`schedule_time`,`end_time`,`times`,`mobile`,`remark`) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
)

// ExecuteSnapshotRepository
type ExecuteSnapshotRepository struct {
	connection *sqlx.DB
}

// new ExecuteSnapshotRepository
func NewExecuteSnapshotRepository(db *db.DB) *ExecuteSnapshotRepository {
	return &ExecuteSnapshotRepository{connection: db.DS}
}

func (repo *ExecuteSnapshotRepository) Insert(snapshot *model.ExecuteSnapshot) {

	repo.connection.Exec(InsertExecuteSnapshotSql, snapshot)
}
