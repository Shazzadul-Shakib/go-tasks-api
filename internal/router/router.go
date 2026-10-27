// Package router builds the Gin engine with all routes and middleware.
package router

import (
	"log/slog"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	"github.com/Shazzadul-Shakib/go-tasks-api/internal/handlers"
	"github.com/Shazzadul-Shakib/go-tasks-api/internal/middleware"
	"github.com/Shazzadul-Shakib/go-tasks-api/internal/modules/task"
	"github.com/Shazzadul-Shakib/go-tasks-api/internal/storage"
)

// New builds the engine. Each module mounts its own routes under /api/v1.
func New(log *slog.Logger) *gin.Engine {
	// Report validation errors with JSON names ("title"), not Go names ("Title").
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(jsonFieldName)
	}

	r := gin.New()
	r.HandleMethodNotAllowed = true // 405 instead of 404 for a wrong method
	r.Use(gin.Recovery(), middleware.RequestID(), middleware.Logger(log))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := r.Group("/api/v1")
	task.RegisterRoutes(api, storage.NewMemoryStore())

	r.NoRoute(func(c *gin.Context) {
		handlers.RespondError(c, http.StatusNotFound, "not_found", "route not found")
	})
	r.NoMethod(func(c *gin.Context) {
		handlers.RespondError(c, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed on this route")
	})
	return r
}

func jsonFieldName(f reflect.StructField) string {
	name, _, _ := strings.Cut(f.Tag.Get("json"), ",")
	if name == "-" {
		return ""
	}
	if name == "" {
		return f.Name
	}
	return name
}
