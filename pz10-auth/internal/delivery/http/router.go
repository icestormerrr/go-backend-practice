package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/icestormerrr/pz10-auth/internal/core"
	"github.com/icestormerrr/pz10-auth/internal/delivery/http/handlers"
	"github.com/icestormerrr/pz10-auth/internal/delivery/http/middleware"
)

func Build(authHandler *handlers.AuthHandler, userHandler *handlers.UserHandler, tokenManager core.TokenManager) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Post("/api/v1/login", authHandler.Login)
	r.Post("/api/v1/refresh", authHandler.Refresh)

	r.Group(func(priv chi.Router) {
		priv.Use(middleware.AuthN(tokenManager))
		priv.Use(middleware.AuthZRoles("admin", "user"))
		priv.Get("/api/v1/me", userHandler.Me)
		priv.Get("/api/v1/user/{id}", userHandler.GetByID)
	})

	r.Group(func(admin chi.Router) {
		admin.Use(middleware.AuthN(tokenManager))
		admin.Use(middleware.AuthZRoles("admin"))
		admin.Get("/api/v1/admin/stats", userHandler.GetAdminStats)
	})

	return r
}
