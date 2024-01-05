package main

import "github.com/muratom/domain-monitoring/services/inspector/internal/core/service"

type DomainInspector interface {
	service.Runnable
}
