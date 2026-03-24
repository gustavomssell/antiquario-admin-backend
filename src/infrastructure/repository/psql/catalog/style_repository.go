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

var ColumnsStyleMapping = map[string]string{
	"id":            "id",
	"name":          "name",
	"description":   "description",
	"periodId":      "period_id",
	"originCountry": "origin_country",
	"createdAt":     "created_at",
	"updatedAt":     "updated_at",
}

type StyleRepositoryInterface interface {
	GetAll() (*[]domainCatalog.Style, error)
	GetByID(id string) (*domainCatalog.Style, error)
	Create(objDomain *domainCatalog.Style) (*domainCatalog.Style, error)
	Update(id string, objMap map[string]interface{}) (*domainCatalog.Style, error)
	Delete(id string) error
	SearchPaginated(filters domain.DataFilters) (*domainCatalog.SearchResultStyle, error)
}

type StyleRepository struct {
	DB     *gorm.DB
	Logger *logger.Logger
}

func NewStyleRepository(db *gorm.DB, loggerInstance *logger.Logger) StyleRepositoryInterface {
	return &StyleRepository{DB: db, Logger: loggerInstance}
}

func arrayToDomainStyleMapper(rows *[]Style) *[]domainCatalog.Style {
	var domainRows []domainCatalog.Style
	for _, c := range *rows {
		domainRows = append(domainRows, *c.toDomainMapper())
	}
	return &domainRows
}

func (r *StyleRepository) GetAll() (*[]domainCatalog.Style, error) {
	var results []Style
	if err := r.DB.Find(&results).Error; err != nil {
		r.Logger.Error("Error getting all styles", zap.Error(err))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	return arrayToDomainStyleMapper(&results), nil
}

func (r *StyleRepository) GetByID(id string) (*domainCatalog.Style, error) {
	var obj Style
	if err := r.DB.Where("id = ?", id).First(&obj).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainErrors.NewAppErrorWithType(domainErrors.NotFound)
		}
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	return obj.toDomainMapper(), nil
}

func (r *StyleRepository) Create(newObj *domainCatalog.Style) (*domainCatalog.Style, error) {
	obj := fromDomainStyleMapper(newObj)
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

func (r *StyleRepository) Update(id string, objMap map[string]interface{}) (*domainCatalog.Style, error) {
	var obj Style
	if err := r.DB.Where("id = ?", id).First(&obj).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainErrors.NewAppErrorWithType(domainErrors.NotFound)
		}
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}

	updateData := make(map[string]interface{})
	for k, v := range objMap {
		if column, ok := ColumnsStyleMapping[k]; ok {
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

func (r *StyleRepository) Delete(id string) error {
	tx := r.DB.Where("id = ?", id).Delete(&Style{})
	if tx.Error != nil {
		return domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	if tx.RowsAffected == 0 {
		return domainErrors.NewAppErrorWithType(domainErrors.NotFound)
	}
	return nil
}

func (r *StyleRepository) SearchPaginated(filters domain.DataFilters) (*domainCatalog.SearchResultStyle, error) {
	query := r.DB.Model(&Style{})

	for field, values := range filters.LikeFilters {
		if len(values) > 0 {
			for _, value := range values {
				if value != "" && ColumnsStyleMapping[field] != "" {
					query = query.Where(ColumnsStyleMapping[field]+" LIKE ?", "%"+value+"%")
				}
			}
		}
	}
	for field, values := range filters.Matches {
		if len(values) > 0 && ColumnsStyleMapping[field] != "" {
			query = query.Where(ColumnsStyleMapping[field]+" IN ?", values)
		}
	}

	var count int64
	query.Count(&count)

	if len(filters.SortBy) > 0 && ColumnsStyleMapping[filters.SortBy[0]] != "" {
		direction := "ASC"
		if filters.SortDirection == domain.SortDesc {
			direction = "DESC"
		}
		query = query.Order(ColumnsStyleMapping[filters.SortBy[0]] + " " + direction)
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

	var list []Style
	if err := query.Find(&list).Error; err != nil {
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}

	return &domainCatalog.SearchResultStyle{
		Data:       arrayToDomainStyleMapper(&list),
		Total:      count,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: int(math.Ceil(float64(count) / float64(pageSize))),
	}, nil
}
