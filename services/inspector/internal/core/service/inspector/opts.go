package inspector

import (
	"time"
)

type Option func(*Service)

func WithTick(t time.Duration) Option {
	return func(s *Service) {
		s.tick = t
	}
}

func WithWorkerNumber(n int) Option {
	return func(s *Service) {
		s.workerNumber = n
	}
}

func WithNotifiers(ns []Notifier) Option {
	return func(s *Service) {
		s.notifiers = ns
	}
}
