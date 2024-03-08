package schedule

import (
	"sync"
	"time"
	"wechatbot/internal/pkg/code"
	appCfg "wechatbot/internal/wechatbot/config"
	"wechatbot/pkg/log"

	"github.com/go-co-op/gocron/v2"
	cronUUID "github.com/google/uuid"
)

var (
	once     sync.Once
	instance gocron.Scheduler
)

type ScheduleLogger struct{}

func (s *ScheduleLogger) Debug(msg string, args ...any) {
	log.Debugf(msg, args...)
}

func (s *ScheduleLogger) Error(msg string, args ...any) {
	log.Errorf(msg, args...)
}

func (s *ScheduleLogger) Info(msg string, args ...any) {
	log.Infof(msg, args...)
}

func (s *ScheduleLogger) Warn(msg string, args ...any) {
	log.Warnf(msg, args...)
}

func GetScheduler() (gocron.Scheduler, error) {
	var err error
	once.Do(func() {
		location, _ := time.LoadLocation("Asia/Shanghai")
		// storeIns, _ := postgres.GetPostgreStore(nil)
		s, errCron := gocron.NewScheduler(
			gocron.WithLocation(location),
			gocron.WithLogger(
				&ScheduleLogger{},
			),
		)
		if errCron != nil {
			err = errCron
			return
		}

		s.NewJob(
			gocron.DurationJob(30*time.Minute),
			gocron.NewTask(
				func() {
					// log.Info("Scheduler every 30 minute triggered")
				},
			),
			gocron.WithName("test-task"),
			gocron.WithEventListeners(
				gocron.AfterJobRuns(
					func(jobID cronUUID.UUID, jobName string) {
						// do something after the job completes
						log.Infof("AfterJobRuns %s::%v triggered", jobName, jobID)
					},
				),
			),
		)
		if appCfg.GCfg.GenericServerRunOptions.Mode != code.ServerModelRelease {
			s.RemoveByTags(code.ServerModelRelease)
		}
		instance = s
	})
	return instance, err
}
