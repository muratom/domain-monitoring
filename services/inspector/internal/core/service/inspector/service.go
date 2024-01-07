package inspector

import (
	"context"
	"github.com/gammazero/workerpool"
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/entity/notification"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"time"
)

type Service struct {
	domainChecker   checker
	domainUpdater   updater
	domainRetriever retriever
	tick            time.Duration
	workerNumber    int
	notifiers       []Notifier
	closeChan       chan struct{}
}

func New(
	checker checker,
	updater updater,
	retriever retriever,
	opts ...Option,
) *Service {
	const (
		defaultTick         = 1 * time.Minute
		defaultWorkerNumber = 5
	)

	s := &Service{
		domainChecker:   checker,
		domainUpdater:   updater,
		domainRetriever: retriever,

		tick:         defaultTick,
		workerNumber: defaultWorkerNumber,
		closeChan:    make(chan struct{}),
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *Service) Start(ctx context.Context) {
	ticker := time.NewTicker(s.tick)
	go func() {
		logrus.Infof("starting cron")
		defer func() {
			logrus.Info("ticker is stopping...")

			ticker.Stop()
		}()

		for {
			select {
			case <-ticker.C:
				st := time.Now()
				ctx, span := otel.Tracer("").Start(ctx, "tick")
				logrus.Infof("Tick")
				rottenFQDNs, err := s.domainRetriever.RetrieveRottenDomainsFQDN(ctx)
				if err != nil {
					logrus.Warnf("get rotten domains' FQDNs: %v", err)
					continue
				}
				span.SetAttributes(attribute.StringSlice("rotten_fqdns", rottenFQDNs))

				wp := workerpool.New(s.workerNumber)
				checkResults := make(chan checkResult, len(rottenFQDNs))
				for _, fqdn := range rottenFQDNs {
					fqdn := fqdn
					wp.Submit(makeFQDNProcessor(ctx, s.domainChecker, s.domainUpdater, fqdn, checkResults))
				}

				wp.StopWait()
				close(checkResults)

				for _, notifier := range s.notifiers {
					for result := range checkResults {
						err := notifier.Notify(result.notifications)
						if err != nil {
							logrus.Errorf("notifier %v: %v", notifier.Name(), err)
						}
					}
				}
				span.End()
				logrus.Infof("elapsed time: %v", time.Since(st))
			case <-ctx.Done():
				return
			case <-s.closeChan:
				return
			}
		}

	}()
}

func (s *Service) Stop(context.Context) {
	s.closeChan <- struct{}{}
}

type checkResult struct {
	notifications []notification.Notification
	err           error
}

func makeFQDNProcessor(
	ctx context.Context,
	domainChecker checker,
	domainUpdater updater,
	fqdn string,
	result chan checkResult,
) func() {
	return func() {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		ctx, span := otel.Tracer("").Start(ctx, "worker", trace.WithAttributes(
			attribute.String("FQDN", fqdn),
		))

		var notifications []notification.Notification
		nots, err := domainChecker.CheckDomainNameServers(ctx, fqdn)
		if err != nil {
			logrus.Warnf("check name servers for FQDN (%v): %v", fqdn, err)
		}
		notifications = append(notifications, nots...)

		nots, err = domainChecker.CheckDomainRegistration(ctx, fqdn)
		if err != nil {
			logrus.Warnf("check registration for FQDN (%v): %v", fqdn, err)
		}
		notifications = append(notifications, nots...)

		nots, err = domainChecker.CheckDomainChanges(ctx, fqdn)
		if err != nil {
			logrus.Warnf("check changes for FQDN (%v): %v", fqdn, err)
		}
		notifications = append(notifications, nots...)

		_, err = domainUpdater.UpdateDomain(ctx, fqdn)
		if err != nil {
			logrus.Warnf("updating domain (%v): %v", fqdn, err)
		}

		result <- checkResult{
			notifications: notifications,
			err:           err,
		}
		span.End()
	}
}