package model

import (
	"context"
	"github.com/gorhill/cronexpr"
	"time"
)

// JobLog is data carrier for the log for job
type JobLog struct {
	JobName      string `json:"jobName" bson:"jobName"`           // job name
	Command      string `json:"command" bson:"command"`           // job shell command
	Err          string `json:"err" bson:"err"`                   // err string
	Output       string `json:"output" bson:"output"`             // job output
	PlanTime     int64  `json:"planTime" bson:"planTime"`         // job planned start time
	ScheduleTime int64  `json:"scheduleTime" bson:"scheduleTime"` // actual scheduled time
	StartTime    int64  `json:"startTime" bson:"startTime"`       // actual start time of the job
	EndTime      int64  `json:"endTime" bson:"endTime"`           // actual job end time
}

// LogBatch
type LogBatch struct {
	Logs []interface{} // multiple logs
}

// Job is a cron job
type Job struct {
	Name     string `json:"name"`     //  job name
	Command  string `json:"command"`  // shell command
	CronExpr string `json:"cronExpr"` // cron expression
}

// JobExecuteInfo wraps around job and has its meta information
type JobExecuteInfo struct {
	Job        *Job               // actual job
	PlanTime   time.Time          // planned exec time in theory
	RealTime   time.Time          // actual exec time
	CancelCtx  context.Context    // context
	CancelFunc context.CancelFunc //  used to cancel the exec
}

// JobExecuteResult
type JobExecuteResult struct {
	ExecuteInfo *JobExecuteInfo
	Output      []byte    // job output
	Err         error     // err
	StartTime   time.Time // start time
	EndTime     time.Time // end time
}

type JobSchedulePlan struct {
	Job *Job
	Expr *cronexpr.Expression
	NextTime time.Time
}

// JobEvent
type JobEvent struct {
	EventType int //  SAVE, DELETE
	Job       *Job
}

const (
	JobEventSave = 1

	JobEventDelete = 2

	JobEventKill = 3

	JobSaveDir = "/croncord/jobs/"

	JobKillerDir = "/croncord/killer/"

	JobLockDir = "/croncord/lock/"

	JobWorkerDir = "/croncord/workers/"
)
