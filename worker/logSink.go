package worker

import (
	"context"
	"github.com/jamesxuhaozhe/croncord/worker/model"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongoLogSink struct {
	client         *mongo.Client
	logCollection  *mongo.Collection
	logChan        chan *model.JobLog
	autoCommitChan chan *model.LogBatch
}

var WorkerLogSink LogSink

type LogSink interface {
	saveLogs(batch *model.LogBatch)
	writeLoop()
	append(jobLog *model.JobLog)
}

func InitLogSink() (err error) {
	var client *mongo.Client

	opt := new(options.ClientOptions)
	opt.SetConnectTimeout(time.Duration(WorkerConfig.MongodbConnectTimeout) * time.Millisecond)
	// init mongodb connection
	if client, err = mongo.Connect(context.TODO(),
		options.Client().ApplyURI(WorkerConfig.MongodbUri),
		opt); err != nil {
		return err
	}

	WorkerLogSink = &MongoLogSink{
		client:         client,
		logCollection:  client.Database("croncord").Collection("log"),
		logChan:        make(chan *model.JobLog, 1000),
		autoCommitChan: make(chan *model.LogBatch, 1000),
	}

	go WorkerLogSink.writeLoop()

	return nil
}

func (logSink *MongoLogSink) saveLogs(batch *model.LogBatch) {
	logSink.logCollection.InsertMany(context.TODO(), batch.Logs)
}

func (logSink *MongoLogSink) writeLoop() {
	var (
		log          *model.JobLog
		logBatch     *model.LogBatch // current batch
		commitTimer  *time.Timer
		timeoutBatch *model.LogBatch // time out batch
	)

	for {
		select {
		case log = <-logSink.logChan:
			if logBatch == nil {
				logBatch = &model.LogBatch{}
				// this logBatch should be committed once 1 second is reached anyway
				commitTimer = time.AfterFunc(
					time.Duration(WorkerConfig.JobLogCommitTimeout)*time.Millisecond,
					func(batch *model.LogBatch) func() {
						return func() {
							logSink.autoCommitChan <- batch
						}
					}(logBatch))
			}

			logBatch.Logs = append(logBatch.Logs, log)

			if len(logBatch.Logs) >= WorkerConfig.JobLogBatchSize {
				logSink.saveLogs(logBatch)
				logBatch = nil
				commitTimer.Stop()
			}
		case timeoutBatch = <- logSink.autoCommitChan:
			if timeoutBatch != logBatch {
				continue
			}

			logSink.saveLogs(timeoutBatch)

			logBatch = nil
		}
	}
}

func (logSink *MongoLogSink) append(jobLog *model.JobLog) {
	select {
	case logSink.logChan <- jobLog:
	default:
		// discard the job log
	}
}
