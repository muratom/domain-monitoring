package service

type Runnable interface {
	Startable
	Stoppable
}

type Startable interface {
	Start()
}

type Stoppable interface {
	Stop()
}
