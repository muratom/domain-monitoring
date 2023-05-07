package http

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/muratom/domain-monitoring/api/rpc/v1/inspector/models"
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/entity"
	"github.com/sirupsen/logrus"
)

type InspectorServer struct {
	domainService DomainService
}

func NewInspectorServer(domainService DomainService) *InspectorServer {
	return &InspectorServer{
		domainService: domainService,
	}
}

func (s *InspectorServer) AddDomain(ctx echo.Context, params models.AddDomainParams) error {
	domain, err := s.domainService.AddDomain(ctx.Request().Context(), params.Fqdn)
	if err != nil {
		resp := models.Error{
			Message: "failed to add domain",
		}
		logrus.Error(err)
		return ctx.JSON(http.StatusInternalServerError, resp)
	}

	return ctx.JSON(http.StatusOK, domain)
}

func (s *InspectorServer) DeleteDomain(ctx echo.Context, params models.DeleteDomainParams) error {
	err := s.domainService.DeleteDomain(ctx.Request().Context(), params.Fqdn)
	if err != nil {
		if err == entity.ErrDomainNotFound {
			return ctx.JSON(http.StatusNotFound, models.Error{Message: fmt.Sprintf("domain %v not found", params.Fqdn)})
		}
		logrus.Error(err)
		return ctx.JSON(http.StatusInternalServerError, models.Error{Message: "failed to delete domain"})
	}

	return ctx.NoContent(http.StatusOK)
}

func (s *InspectorServer) GetDomain(ctx echo.Context, params models.GetDomainParams) error {
	domain, err := s.domainService.GetDomain(ctx.Request().Context(), params.Fqdn)
	if err != nil {
		if err == entity.ErrDomainNotFound {
			return ctx.JSON(http.StatusNotFound, models.Error{Message: fmt.Sprintf("domain %v not found", params.Fqdn)})
		}
		logrus.Error(err)
		return ctx.JSON(http.StatusInternalServerError, models.Error{Message: "failed to get domain"})
	}

	return ctx.JSON(http.StatusOK, domain)
}

func (s *InspectorServer) GetAllDomains(ctx echo.Context) error {
	domains, err := s.domainService.GetAllDomainsFQDN(ctx.Request().Context())
	if err != nil {
		logrus.Error(err)
		return ctx.JSON(http.StatusInternalServerError, models.Error{Message: "failed to get all domains"})
	}

	return ctx.JSON(http.StatusOK, domains)
}

func (s *InspectorServer) UpdateDomain(ctx echo.Context, params models.UpdateDomainParams) error {
	domain, err := s.domainService.UpdateDomain(ctx.Request().Context(), params.Fqdn)
	if err != nil {
		if err == entity.ErrDomainNotFound {
			return ctx.JSON(http.StatusNotFound, models.Error{Message: fmt.Sprintf("domain %v not found", params.Fqdn)})
		}
		logrus.Error(err)
		return ctx.JSON(http.StatusInternalServerError, models.Error{Message: "failed to update domain"})
	}

	return ctx.JSON(http.StatusOK, domain)
}

func (s *InspectorServer) GetChangelogs(ctx echo.Context, params models.GetChangelogsParams) error {
	changelogs, err := s.domainService.GetChangelogs(ctx.Request().Context(), params.Fqdn)
	if err != nil {
		logrus.Error(err)
		return ctx.JSON(http.StatusInternalServerError, models.Error{Message: "failed to get domain's changelogs"})
	}

	return ctx.JSON(http.StatusOK, changelogs)
}

func (s *InspectorServer) Ping(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "pong")
}
