package worker

import (
	"github.com/jamesxuhaozhe/croncord/worker/model"
	"math/rand"
	"os/exec"
	"time"
)

type DefaultExecutor struct {

}

type Executor interface {
	Execute(jobInfo *model.JobExecuteInfo)
}

var WorkerExecutor *DefaultExecutor


func (exe *DefaultExecutor) Execute(jobInfo *model.JobExecuteInfo) {
	go func() {

		// the exec result
		result := &model.JobExecuteResult{
			ExecuteInfo: jobInfo,
			Output:      make([]byte, 0),
		}

		// TODO create the joblock

		// record the job start time
		result.StartTime = time.Now()

		// give some random time
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

		// try lock and defer unlock
		if err != nil {
			result.Err = err
			result.EndTime = time.Now()
		} else {
			// if the lock is acquired, the job start time needs to re-assigned
			result.StartTime = time.Now()

			// invoke the shell command
			cmd := exec.CommandContext(jobInfo.CancelCtx, "/bin/bash", "-c", jobInfo.Job.Command)

			// catch the output
			output, err = cmd.CombinedOutput()

			// record the time
			result.EndTime = time.Now()
			result.Output = output
			result.Err = err
		}

		// once the job is finished we need to publish the event to scheduler, scheduler need to remove it from the table
		// TODO
	}()
}

func InitExecutor() error {
	WorkerExecutor = &DefaultExecutor{}
	return nil
}




