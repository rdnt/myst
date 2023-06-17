package scheduler

import (
	"time"

	"github.com/go-co-op/gocron"
	"github.com/pkg/errors"

	"myst/pkg/logger"
	"myst/src/client/application"
	"myst/src/client/remote"
)

type Scheduler struct {
	app    application.Application
	remote *remote.Remote
	sched  *gocron.Scheduler
	log    *logger.Logger
}

func New(app application.Application, remote *remote.Remote) (*Scheduler, error) {
	s := &Scheduler{
		app:    app,
		remote: remote,
		sched:  gocron.NewScheduler(time.UTC),
		log:    logger.New("scheduler", logger.DefaultColor),
	}

	_, err := s.sched.Every(10).Second().
		StartAt(time.Now()).
		Do(s.sync)
	if err != nil {
		return nil, errors.Wrap(err, "failed to schedule sync")
	}

	_, err = s.sched.Every(10).Minute().
		StartAt(time.Now()).
		Do(s.reauthenticate)
	if err != nil {
		return nil, errors.Wrap(err, "failed to schedule silentAuth")
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

func (s *Scheduler) reauthenticate() {
	err := s.remote.Reauthenticate()
	if err != nil {
		s.log.Error(err)
	}
}
