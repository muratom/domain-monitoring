package http

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/muratom/domain-monitoring/api/rpc/v1/inspector/models"
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/entity"
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/service"
)

type InspectorServer struct {
	domainService service.DomainService
}

func (s *InspectorServer) AddDomain(ctx echo.Context, params models.AddDomainParams) error {
	domain, err := s.domainService.AddDomain(ctx.Request().Context(), params.Fqdn)
	if err != nil {
		resp := models.Error{
			Message: "failed to add domain",
		}
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
		return ctx.JSON(http.StatusInternalServerError, models.Error{Message: "failed to get domain"})
	}

	return ctx.JSON(http.StatusOK, domain)
}

func (s *InspectorServer) UpdateDomain(ctx echo.Context, params models.UpdateDomainParams) error {
	domain, err := s.domainService.UpdateDomain(ctx.Request().Context(), params.Fqdn)
	if err != nil {
		if err == entity.ErrDomainNotFound {
			return ctx.JSON(http.StatusNotFound, models.Error{Message: fmt.Sprintf("domain %v not found", params.Fqdn)})
		}
		return ctx.JSON(http.StatusInternalServerError, models.Error{Message: "failed to update domain"})
	}

	return ctx.JSON(http.StatusOK, domain)
}

func (s *InspectorServer) Ping(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "pong")
}
