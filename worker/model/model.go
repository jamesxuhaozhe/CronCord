package model

// JobLog is data carrier for the log for job
type JobLog struct {
	JobName string `json:"jobName" bson:"jobName"` // job name
	Command string `json:"command" bson:"command"` // job shell command
	Err string `json:"err" bson:"err"` // err string
	Output string `json:"output" bson:"output"`	// job output
	PlanTime int64 `json:"planTime" bson:"planTime"` // job planned start time
	ScheduleTime int64 `json:"scheduleTime" bson:"scheduleTime"` // actual scheduled time
	StartTime int64 `json:"startTime" bson:"startTime"` // actual start time of the job
	EndTime int64 `json:"endTime" bson:"endTime"` // actual job end time
}

// LogBatch
type LogBatch struct {
	Logs []interface{}	// multiple logs
}