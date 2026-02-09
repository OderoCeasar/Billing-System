package repositories

import (
	"github.com/OderoCeasar/system/db/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)


type PackageRepository struct {
	db *gorm.DB
}

func NewPackageRepository(db *gorm.DB) *PackageRepository {
	return &PackageRepository{db : db}
}

func (r *PackageRepository) Create(pkg *models.Package) error {
	return r.db.Create(pkg).Error
} 

func (r *PackageRepository) FindByID(id uuid.UUID) (*models.Package, error) {
	var pkg models.Package
	err := r.db.Where("id = ?", id).First(&pkg).Error
	return &pkg, err
}

func (r *PackageRepository) Update(pkg *models.Package) error {
	return r.db.Save(pkg).Error
}

func (r *PackageRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Package{}, id).Error
}

func (r *PackageRepository) ListActive() ([]models.Package, error) {
	var packages []models.Package
	err := r.db.Where("is_active = ?", true).Find(&packages).Error
	return packages, err
}

func (r *PackageRepository) List() ([]models.Package, error) {
	var packages []models.Package
	err := r.db.Find(&packages).Error
	return packages, err
}