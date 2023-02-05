package job

import (
	"context"

	"github.com/dogefuzz/dogefuzz/pkg/interfaces"
	cron "github.com/robfig/cron/v3"
)

type scheduler struct {
	scheduler *cron.Cron
	env       Env
}

func NewJobScheduler(env Env) *scheduler {
	return &scheduler{env: env}
}

func (s *scheduler) Start() {
	s.scheduler = cron.New(cron.WithSeconds())
	jobs := s.getAvailableJobs()
	for _, id := range s.env.Config().JobConfig.EnabledJobs {
		if cronjJob, ok := jobs[id]; ok {
			s.scheduler.AddFunc(cronjJob.CronConfig(), cronjJob.Handler)
			s.env.Logger().Sugar().Infof("starting job %s", id)
		} else {
			s.env.Logger().Sugar().Warnf("ignore job %s because it's not implemented", id)
		}
	}

	s.env.Logger().Info("starting job scheduler")
	go s.scheduler.Run()
}

func (s *scheduler) Shutdown() context.Context {
	s.env.Logger().Info("stoping job scheduler")
	return s.scheduler.Stop()
}

func (s *scheduler) getAvailableJobs() map[string]interfaces.CronJob {
	return map[string]interfaces.CronJob{
		s.env.TasksCheckerJob().Id():        s.env.TasksCheckerJob(),
		s.env.TransactionsCheckerJob().Id(): s.env.TransactionsCheckerJob(),
	}
}
