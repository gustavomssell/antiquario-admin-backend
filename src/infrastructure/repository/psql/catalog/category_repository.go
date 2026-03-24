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

var ColumnsCategoryMapping = map[string]string{
	"id":          "id",
	"name":        "name",
	"description": "description",
	"parentId":    "parent_id",
	"level":       "level",
	"active":      "active",
	"createdAt":   "created_at",
	"updatedAt":   "updated_at",
}

type CategoryRepositoryInterface interface {
	GetAll() (*[]domainCatalog.Category, error)
	GetByID(id string) (*domainCatalog.Category, error)
	Create(objDomain *domainCatalog.Category) (*domainCatalog.Category, error)
	Update(id string, objMap map[string]interface{}) (*domainCatalog.Category, error)
	Delete(id string) error
	SearchPaginated(filters domain.DataFilters) (*domainCatalog.SearchResultCategory, error)
}

type CategoryRepository struct {
	DB     *gorm.DB
	Logger *logger.Logger
}

func NewCategoryRepository(db *gorm.DB, loggerInstance *logger.Logger) CategoryRepositoryInterface {
	return &CategoryRepository{DB: db, Logger: loggerInstance}
}

func arrayToDomainCategoryMapper(categories *[]Category) *[]domainCatalog.Category {
	var domainCategories []domainCatalog.Category
	for _, c := range *categories {
		domainCategories = append(domainCategories, *c.toDomainMapper())
	}
	return &domainCategories
}

func (r *CategoryRepository) GetAll() (*[]domainCatalog.Category, error) {
	var results []Category
	if err := r.DB.Find(&results).Error; err != nil {
		r.Logger.Error("Error getting all categories", zap.Error(err))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	return arrayToDomainCategoryMapper(&results), nil
}

func (r *CategoryRepository) GetByID(id string) (*domainCatalog.Category, error) {
	var obj Category
	if err := r.DB.Where("id = ?", id).First(&obj).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainErrors.NewAppErrorWithType(domainErrors.NotFound)
		}
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	return obj.toDomainMapper(), nil
}

func (r *CategoryRepository) Create(newObj *domainCatalog.Category) (*domainCatalog.Category, error) {
	obj := fromDomainCategoryMapper(newObj)
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

func (r *CategoryRepository) Update(id string, objMap map[string]interface{}) (*domainCatalog.Category, error) {
	var obj Category
	if err := r.DB.Where("id = ?", id).First(&obj).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainErrors.NewAppErrorWithType(domainErrors.NotFound)
		}
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}

	updateData := make(map[string]interface{})
	for k, v := range objMap {
		if column, ok := ColumnsCategoryMapping[k]; ok {
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

func (r *CategoryRepository) Delete(id string) error {
	tx := r.DB.Where("id = ?", id).Delete(&Category{})
	if tx.Error != nil {
		return domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	if tx.RowsAffected == 0 {
		return domainErrors.NewAppErrorWithType(domainErrors.NotFound)
	}
	return nil
}

func (r *CategoryRepository) SearchPaginated(filters domain.DataFilters) (*domainCatalog.SearchResultCategory, error) {
	query := r.DB.Model(&Category{})

	for field, values := range filters.LikeFilters {
		if len(values) > 0 {
			for _, value := range values {
				if value != "" && ColumnsCategoryMapping[field] != "" {
					query = query.Where(ColumnsCategoryMapping[field]+" LIKE ?", "%"+value+"%")
				}
			}
		}
	}
	for field, values := range filters.Matches {
		if len(values) > 0 && ColumnsCategoryMapping[field] != "" {
			query = query.Where(ColumnsCategoryMapping[field]+" IN ?", values)
		}
	}

	var count int64
	query.Count(&count)

	if len(filters.SortBy) > 0 && ColumnsCategoryMapping[filters.SortBy[0]] != "" {
		direction := "ASC"
		if filters.SortDirection == domain.SortDesc {
			direction = "DESC"
		}
		query = query.Order(ColumnsCategoryMapping[filters.SortBy[0]] + " " + direction)
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

	var list []Category
	if err := query.Find(&list).Error; err != nil {
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}

	return &domainCatalog.SearchResultCategory{
		Data:       arrayToDomainCategoryMapper(&list),
		Total:      count,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: int(math.Ceil(float64(count) / float64(pageSize))),
	}, nil
}
