package worker

import (
	"github.com/jamesxuhaozhe/croncord/worker/model"
	"time"
)

type DefaultScheduler struct {
	jobEventChan chan *model.JobEvent
	jobPlanTable map[string]*model.JobSchedulePlan
	jobExecutingTable map[string]*model.JobExecuteInfo
	jobResultChan chan *model.JobExecuteResult
}

type Scheduler interface {
	handleJobEvent(jobEvent *model.JobEvent)
	tryStartJob(jobPlan *model.JobSchedulePlan)
	trySchedule() (scheduleAfter time.Duration)
	handleJobResult(result *model.JobExecuteResult)
	scheduleLoop()
	PushJobEvent(jobEvent *model.JobEvent)
	PushJobResult(jobResult *model.JobExecuteResult)
}

var jobScheduler Scheduler

func (sch *DefaultScheduler) handleJobEvent(jobEvent *model.JobEvent) {

}

func (sch *DefaultScheduler) tryStartJob(jobPlan *model.JobSchedulePlan) {
	panic("implement me")
}

func (sch *DefaultScheduler) trySchedule() (scheduleAfter time.Duration) {
	panic("implement me")
}

func (sch *DefaultScheduler) handleJobResult(result *model.JobExecuteResult) {
	panic("implement me")
}

func (sch *DefaultScheduler) scheduleLoop() {
	panic("implement me")
}

func (sch *DefaultScheduler) PushJobEvent(jobEvent *model.JobEvent) {
	sch.jobEventChan <- jobEvent
}

func (sch *DefaultScheduler) PushJobResult(jobResult *model.JobExecuteResult) {
	sch.jobResultChan <- jobResult
}


