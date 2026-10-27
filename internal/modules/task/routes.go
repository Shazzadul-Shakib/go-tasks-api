// Package task is the tasks module: controller, service, repository and routes.
package task

import "github.com/gin-gonic/gin"

// RegisterRoutes wires the module together and mounts its routes on rg.
func RegisterRoutes(rg *gin.RouterGroup, repo Repository) {
	ctrl := NewController(NewService(repo))

	g := rg.Group("/tasks")
	g.GET("", ctrl.List)
}
