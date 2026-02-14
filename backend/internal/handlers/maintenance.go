package handlers

import (
	"net/http"
	"strconv"

	"assetsentinel/internal/middleware"
	"assetsentinel/internal/repository"

	"github.com/gin-gonic/gin"
)

type MaintenanceHandler struct {
	maintenanceService interface {
		Create(plan *repository.MaintenancePlan) error
		Get(id, orgID uint) (*repository.MaintenancePlan, error)
		List(orgID uint, page, pageSize int) ([]repository.MaintenancePlan, int, error)
		Update(plan *repository.MaintenancePlan) error
		Delete(id, orgID uint) error
	}
}

func NewMaintenanceHandler(maintenanceService interface {
	Create(plan *repository.MaintenancePlan) error
	Get(id, orgID uint) (*repository.MaintenancePlan, error)
	List(orgID uint, page, pageSize int) ([]repository.MaintenancePlan, int, error)
	Update(plan *repository.MaintenancePlan) error
	Delete(id, orgID uint) error
}) *MaintenanceHandler {
	return &MaintenanceHandler{maintenanceService: maintenanceService}
}

func (h *MaintenanceHandler) Create(c *gin.Context) {
	var plan repository.MaintenancePlan
	if err := c.ShouldBindJSON(&plan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	plan.OrganizationID = middleware.GetOrganizationID(c)

	if err := h.maintenanceService.Create(&plan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, plan)
}

func (h *MaintenanceHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	orgID := middleware.GetOrganizationID(c)

	plan, err := h.maintenanceService.Get(uint(id), orgID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Maintenance plan not found"})
		return
	}

	c.JSON(http.StatusOK, plan)
}

func (h *MaintenanceHandler) List(c *gin.Context) {
	orgID := middleware.GetOrganizationID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	plans, total, err := h.maintenanceService.List(orgID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      plans,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (h *MaintenanceHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	orgID := middleware.GetOrganizationID(c)

	var plan repository.MaintenancePlan
	if err := c.ShouldBindJSON(&plan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	plan.ID = uint(id)
	plan.OrganizationID = orgID

	if err := h.maintenanceService.Update(&plan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, plan)
}

func (h *MaintenanceHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	orgID := middleware.GetOrganizationID(c)

	if err := h.maintenanceService.Delete(uint(id), orgID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Maintenance plan deleted"})
}
