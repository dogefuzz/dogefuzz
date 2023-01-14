package job

import (
	"context"

	"github.com/dogefuzz/dogefuzz/config"
	"github.com/dogefuzz/dogefuzz/pkg/interfaces"
	cron "github.com/robfig/cron/v3"
)

type scheduler struct {
	scheduler *cron.Cron
	cfg       *config.Config
	env       *env
}

func NewJobScheduler(cfg *config.Config) *scheduler {
	return &scheduler{cfg: cfg, env: NewEnv(cfg)}
}

func (s *scheduler) Start() {
	s.scheduler = cron.New()
	jobs := s.getAvailableJobs()
	for _, id := range s.cfg.JobConfig.EnabledJobs {
		if cronjJob, ok := jobs[id]; ok {
			s.scheduler.AddFunc(cronjJob.CronConfig(), cronjJob.Handler)
			s.env.logger.Sugar().Infof("starting job %s", id)
		} else {
			s.env.logger.Sugar().Warnf("ignore job %s because it's not implemented", id)
		}
	}

	s.env.logger.Info("starting job scheduler")
	go s.scheduler.Run()
}

func (s *scheduler) Shutdown() context.Context {
	s.env.logger.Info("stoping job scheduler")
	return s.scheduler.Stop()
}

func (s *scheduler) getAvailableJobs() map[string]interfaces.CronJob {
	return map[string]interfaces.CronJob{
		s.env.TasksCheckerJob().Id():        s.env.TasksCheckerJob(),
		s.env.TransactionsCheckerJob().Id(): s.env.TransactionsCheckerJob(),
	}
}
