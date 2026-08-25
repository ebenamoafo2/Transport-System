package handler

import (
	"net/http"

	"github.com/ebenamoafo2/transport/ride/internal/adapters/http/api"
	"github.com/ebenamoafo2/transport/ride/internal/ports"
	"github.com/gin-gonic/gin"
)

// AssignmentHandler is the HTTP adapter. It implements api.ServiceInterface
// and delegates to the core Assignment(a hexagonal port)
type AssignmentHandler struct {
	service ports.AssignmentService
}

func NewAssignmentHandler(service ports.AssignmentService) *AssignmentHandler {
	return &AssignmentHandler{
		service: service,
	}
}

func (h *AssignmentHandler) ListAssignments(c *gin.Context, params api.ListAssignmentsParams) {
	notImplemented(c)
}

func (h *AssignmentHandler) CreateAssignment(c *gin.Context) {
	notImplemented(c)
}

func (h *AssignmentHandler) GetAssignment(c *gin.Context, id string) {
	notImplemented(c)
}

func notImplemented(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"error":   "not implemented",
		"message": "This endpoint is not yet implemented.",
	})
}

// Ensure AssignmentHandler implements the api.ServerInterface
var _ api.ServerInterface = (*AssignmentHandler)(nil)
