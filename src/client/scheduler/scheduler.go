package scheduler

import (
	"time"

	"github.com/go-co-op/gocron"

	"myst/pkg/logger"
	"myst/src/client/application"
)

type Scheduler struct {
	app   application.Application
	sched *gocron.Scheduler
	log   *logger.Logger
}

func New(app application.Application) (*Scheduler, error) {
	s := &Scheduler{
		app:   app,
		sched: gocron.NewScheduler(time.UTC),
		log:   logger.New("scheduler", logger.DefaultColor),
	}

	_, err := s.sched.Every(10).Second().
		StartAt(time.Now()).
		Do(s.sync)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Scheduler) Start() error {
	s.sched.StartAsync()

	return nil
}

func (s *Scheduler) Stop() error {
	s.sched.Stop()

	return nil
}

func (s *Scheduler) sync() {
	err := s.app.Sync()
	if err != nil {
		s.log.Error(err)
	}
}
