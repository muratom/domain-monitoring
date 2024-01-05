package service

import "context"

type Runnable interface {
	Startable
	Stoppable
}

type Startable interface {
	Start(ctx context.Context)
}

type Stoppable interface {
	Stop(ctx context.Context)
}
