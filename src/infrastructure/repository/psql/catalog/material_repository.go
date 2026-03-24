package catalog

import (
	"encoding/json"
	"math"

	"github.com/gbrayhan/microservices-go/src/domain"
	domainCatalog "github.com/gbrayhan/microservices-go/src/domain/catalog"
	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var ColumnsMaterialMapping = map[string]string{
	"id":          "id",
	"name":        "name",
	"description": "description",
	"createdAt":   "created_at",
	"updatedAt":   "updated_at",
}

type MaterialRepositoryInterface interface {
	GetAll() (*[]domainCatalog.Material, error)
	GetByID(id string) (*domainCatalog.Material, error)
	Create(objDomain *domainCatalog.Material) (*domainCatalog.Material, error)
	Update(id string, objMap map[string]interface{}) (*domainCatalog.Material, error)
	Delete(id string) error
	SearchPaginated(filters domain.DataFilters) (*domainCatalog.SearchResultMaterial, error)
}

type MaterialRepository struct {
	DB     *gorm.DB
	Logger *logger.Logger
}

func NewMaterialRepository(db *gorm.DB, loggerInstance *logger.Logger) MaterialRepositoryInterface {
	return &MaterialRepository{DB: db, Logger: loggerInstance}
}

func arrayToDomainMaterialMapper(rows *[]Material) *[]domainCatalog.Material {
	var domainRows []domainCatalog.Material
	for _, c := range *rows {
		domainRows = append(domainRows, *c.toDomainMapper())
	}
	return &domainRows
}

func (r *MaterialRepository) GetAll() (*[]domainCatalog.Material, error) {
	var results []Material
	if err := r.DB.Find(&results).Error; err != nil {
		r.Logger.Error("Error getting all materials", zap.Error(err))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	return arrayToDomainMaterialMapper(&results), nil
}

func (r *MaterialRepository) GetByID(id string) (*domainCatalog.Material, error) {
	var obj Material
	if err := r.DB.Where("id = ?", id).First(&obj).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainErrors.NewAppErrorWithType(domainErrors.NotFound)
		}
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	return obj.toDomainMapper(), nil
}

func (r *MaterialRepository) Create(newObj *domainCatalog.Material) (*domainCatalog.Material, error) {
	obj := fromDomainMaterialMapper(newObj)
	if err := r.DB.Create(obj).Error; err != nil {
		byteErr, _ := json.Marshal(err)
		var newError domainErrors.GormErr
		if errUnmarshal := json.Unmarshal(byteErr, &newError); errUnmarshal == nil {
			if newError.Number == 1062 || newError.Number == 23505 {
				return nil, domainErrors.NewAppErrorWithType(domainErrors.ResourceAlreadyExists)
			}
		}
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	return obj.toDomainMapper(), nil
}

func (r *MaterialRepository) Update(id string, objMap map[string]interface{}) (*domainCatalog.Material, error) {
	var obj Material
	if err := r.DB.Where("id = ?", id).First(&obj).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainErrors.NewAppErrorWithType(domainErrors.NotFound)
		}
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}

	updateData := make(map[string]interface{})
	for k, v := range objMap {
		if column, ok := ColumnsMaterialMapping[k]; ok {
			updateData[column] = v
		} else {
			updateData[k] = v
		}
	}

	if err := r.DB.Model(&obj).Select("*").Updates(updateData).Error; err != nil {
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}

	r.DB.Where("id = ?", id).First(&obj)
	return obj.toDomainMapper(), nil
}

func (r *MaterialRepository) Delete(id string) error {
	tx := r.DB.Where("id = ?", id).Delete(&Material{})
	if tx.Error != nil {
		return domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	if tx.RowsAffected == 0 {
		return domainErrors.NewAppErrorWithType(domainErrors.NotFound)
	}
	return nil
}

func (r *MaterialRepository) SearchPaginated(filters domain.DataFilters) (*domainCatalog.SearchResultMaterial, error) {
	query := r.DB.Model(&Material{})

	for field, values := range filters.LikeFilters {
		if len(values) > 0 {
			for _, value := range values {
				if value != "" && ColumnsMaterialMapping[field] != "" {
					query = query.Where(ColumnsMaterialMapping[field]+" LIKE ?", "%"+value+"%")
				}
			}
		}
	}
	for field, values := range filters.Matches {
		if len(values) > 0 && ColumnsMaterialMapping[field] != "" {
			query = query.Where(ColumnsMaterialMapping[field]+" IN ?", values)
		}
	}

	var count int64
	query.Count(&count)

	if len(filters.SortBy) > 0 && ColumnsMaterialMapping[filters.SortBy[0]] != "" {
		direction := "ASC"
		if filters.SortDirection == domain.SortDesc {
			direction = "DESC"
		}
		query = query.Order(ColumnsMaterialMapping[filters.SortBy[0]] + " " + direction)
	} else {
		query = query.Order("created_at DESC")
	}

	page, pageSize := 1, 10
	if filters.Page > 0 {
		page = filters.Page
	}
	if filters.PageSize > 0 {
		pageSize = filters.PageSize
	}
	query = query.Offset((page - 1) * pageSize).Limit(pageSize)

	var list []Material
	if err := query.Find(&list).Error; err != nil {
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}

	return &domainCatalog.SearchResultMaterial{
		Data:       arrayToDomainMaterialMapper(&list),
		Total:      count,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: int(math.Ceil(float64(count) / float64(pageSize))),
	}, nil
}
