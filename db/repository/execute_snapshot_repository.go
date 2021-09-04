package repository

import (
	"github.com/busgo/pink/db"
	"github.com/jmoiron/sqlx"
)

// ExecuteSnapshotRepository
type ExecuteSnapshotRepository struct {
	connection *sqlx.DB
}

// new ExecuteSnapshotRepository
func NewExecuteSnapshotRepository(db *db.DB) *ExecuteSnapshotRepository {
	return &ExecuteSnapshotRepository{connection: db.DS}
}
