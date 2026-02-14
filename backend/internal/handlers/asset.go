package handlers

import (
	"net/http"
	"strconv"

	"assetsentinel/internal/middleware"
	"assetsentinel/internal/repository"

	"github.com/gin-gonic/gin"
)

type AssetHandler struct {
	assetService interface {
		Create(asset *repository.Asset) error
		Get(id, orgID uint) (*repository.Asset, error)
		List(orgID uint, page, pageSize int, status, category string) ([]repository.Asset, int, error)
		Update(asset *repository.Asset) error
		Delete(id, orgID uint) error
	}
}

func NewAssetHandler(assetService interface {
	Create(asset *repository.Asset) error
	Get(id, orgID uint) (*repository.Asset, error)
	List(orgID uint, page, pageSize int, status, category string) ([]repository.Asset, int, error)
	Update(asset *repository.Asset) error
	Delete(id, orgID uint) error
}) *AssetHandler {
	return &AssetHandler{assetService: assetService}
}

func (h *AssetHandler) Create(c *gin.Context) {
	var asset repository.Asset
	if err := c.ShouldBindJSON(&asset); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	asset.OrganizationID = middleware.GetOrganizationID(c)

	if err := h.assetService.Create(&asset); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, asset)
}

func (h *AssetHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	orgID := middleware.GetOrganizationID(c)

	asset, err := h.assetService.Get(uint(id), orgID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Asset not found"})
		return
	}

	c.JSON(http.StatusOK, asset)
}

func (h *AssetHandler) List(c *gin.Context) {
	orgID := middleware.GetOrganizationID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")
	category := c.Query("category")

	assets, total, err := h.assetService.List(orgID, page, pageSize, status, category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      assets,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (h *AssetHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	orgID := middleware.GetOrganizationID(c)

	var asset repository.Asset
	if err := c.ShouldBindJSON(&asset); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	asset.ID = uint(id)
	asset.OrganizationID = orgID

	if err := h.assetService.Update(&asset); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, asset)
}

func (h *AssetHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	orgID := middleware.GetOrganizationID(c)

	if err := h.assetService.Delete(uint(id), orgID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Asset deleted"})
}
