package handlers

import (
	"net/http"

	"github.com/OderoCeasar/system/db/models"
	"github.com/OderoCeasar/system/db/repositories"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)


type PackageHandler struct {
	packageRepo *repositories.PackageRepository
}

func NewPackageHandler(packageRepo *repositories.PackageRepository) *PackageHandler {
	return &PackageHandler{packageRepo: packageRepo}
}


type CreatePackageRequest struct {
	Name 		string		`json:"name" binding:"required"`
	Description	string		`json:"description"`
	PackageType string 		`json:"package_type" binding:"required,oneOf=time data"`
	Price 		float64		`json:"price" binding:"required,min=1"`
	DurationMinutes	int		`json:"duration_minutes"`
	DataLimitMD		int64	`json:"data_limit_mb"`
	SpeedLimitUp 	int 	`json:"speed_limit_up"`
	SpeedLimitDown  int		`json:"speed_limit_down"`
	ValidityDays  	int		`json:"validity_days" binding:"min=1"`
}


func (h * PackageHandler) ListPackage(c *gin.Context) {
	packages, err := h.packageRepo.ListActive()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":"failed to fetch the packages"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"packages": packages,	
	})
}


func (h *PackageHandler) GetPackage(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error":"invalid package ID"})
		return
	}

	pkg, err := h.packageRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error":"package not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"package": pkg,
	})
}


func (h *PackageHandler) CreatePackage(c *gin.Context) {
	var req CreatePackageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pkg := &models.Package{
		Name:		req.Name,
		Description: req.Description,
		PackageType: models.PackageType(req.PackageType),
		Price: req.Price,
		DurationMinutes: req.DurationMinutes,
		DataLimitMB: req.DataLimitMD,
		SpeedLimitUp: req.SpeedLimitUp,
		SpeedLimitDown: req.SpeedLimitDown,
		ValidityDays: req.ValidityDays,
		IsActive: true,
	}

	if err := h.packageRepo.Create(pkg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":"failed to create package"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":	"package created successfully",
		"package": pkg,
	})
}

func (h *PackageHandler) UpdatePackage(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
		return
	}

	pkg, err := h.packageRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error":"package not found"})
		return
	}

	var req CreatePackageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pkg.Name = req.Name
	pkg.Description = req.Description
	pkg.PackageType = models.PackageType(req.PackageType)
	pkg.Price = req.Price
	pkg.DurationMinutes = req.DurationMinutes
	pkg.DataLimitMB = req.DataLimitMD
	pkg.SpeedLimitUp = req.SpeedLimitUp
	pkg.SpeedLimitDown = req.SpeedLimitDown
	pkg.ValidityDays = req.ValidityDays
	

	if err := h.packageRepo.Update(pkg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":"failed to update the package"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "package updated successfully",
		"package": pkg,
	})
}	


func (h *PackageHandler) DeletePackage(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error":"invalid package ID"})
		return
	}

	if err := h.packageRepo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":"failed to delete the package"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "package deleted successfully",
		
	})

	
}