package inspector

import (
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/service"
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

func WithNotifiers(ns []service.Notifier) Option {
	return func(s *Service) {
		s.notifiers = ns
	}
}
