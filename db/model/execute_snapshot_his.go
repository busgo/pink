package model

import "encoding/json"

// CREATE TABLE `execute_snapshot_his` (
//  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ä¸»é”®',
//  `snapshot_id` varchar(32) NOT NULL,
//  `job_name` varchar(32) NOT NULL,
//  `group` varchar(32) NOT NULL,
//  `cron` varchar(255) DEFAULT NULL,
//  `target` varchar(255) DEFAULT NULL,
//  `ip` varchar(32) DEFAULT NULL,
//  `param` varchar(255) DEFAULT NULL,
//  `state` tinyint(4) DEFAULT NULL,
//  `before_time` datetime NOT NULL,
//  `schedule_time` datetime NOT NULL,
//  `end_time` datetime DEFAULT NULL,
//  `times` bigint(20) NOT NULL,
//  `mobile` varchar(32) DEFAULT NULL,
//  `remark` varchar(32) DEFAULT NULL,
//  PRIMARY KEY (`id`)
//) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
type ExecuteSnapshotHis struct {
	Id           int64  `db:"id"`
	SnapshotId   string `db:"snapshot_id"`
	JobId        string `db:"job_id"`
	JobName      string `db:"job_name"`
	Group        string `db:"group"`
	Cron         string `db:"cron"`
	Target       string `db:"target"`
	Ip           string `db:"ip"`
	Param        string `db:"param"`
	State        int32  `db:"state"`
	BeforeTime   string `db:"before_time"`
	ScheduleTime string `db:"schedule_time"`
	EndTime      string `db:"end_time"`
	Times        int64  `db:"times"`
	Mobile       string `db:"mobile"`
	Remark       string `db:"remark"`
}

//
func (his *ExecuteSnapshotHis) Decode(content string) *ExecuteSnapshotHis {

	if len(content) == 0 {
		return nil
	}
	err := json.Unmarshal([]byte(content), his)
	if err != nil {
		return nil
	}
	return his
}
