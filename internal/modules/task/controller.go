package task

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Shazzadul-Shakib/go-tasks-api/internal/handlers"
	"github.com/Shazzadul-Shakib/go-tasks-api/internal/models"
)

type Controller struct {
	svc *Service
}

func NewController(svc *Service) *Controller {
	return &Controller{svc: svc}
}

// List handles GET /tasks, with an optional ?done=true|false filter.
func (h *Controller) List(c *gin.Context) {
	var f models.ListFilter
	if raw := c.Query("done"); raw != "" {
		done, err := strconv.ParseBool(raw)
		if err != nil {
			handlers.RespondError(c, http.StatusBadRequest, "bad_request", "query param \"done\" must be true or false")
			return
		}
		f.Done = &done
	}

	tasks, err := h.svc.List(c.Request.Context(), f)
	if err != nil {
		handlers.RespondError(c, http.StatusInternalServerError, "internal_error", "could not list tasks")
		return
	}
	c.JSON(http.StatusOK, tasks)
}
