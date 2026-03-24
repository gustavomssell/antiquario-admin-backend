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

var ColumnsPeriodMapping = map[string]string{
	"id":          "id",
	"name":        "name",
	"description": "description",
	"startYear":   "start_year",
	"endYear":     "end_year",
	"createdAt":   "created_at",
	"updatedAt":   "updated_at",
}

type PeriodRepositoryInterface interface {
	GetAll() (*[]domainCatalog.Period, error)
	GetByID(id string) (*domainCatalog.Period, error)
	Create(objDomain *domainCatalog.Period) (*domainCatalog.Period, error)
	Update(id string, objMap map[string]interface{}) (*domainCatalog.Period, error)
	Delete(id string) error
	SearchPaginated(filters domain.DataFilters) (*domainCatalog.SearchResultPeriod, error)
}

type PeriodRepository struct {
	DB     *gorm.DB
	Logger *logger.Logger
}

func NewPeriodRepository(db *gorm.DB, loggerInstance *logger.Logger) PeriodRepositoryInterface {
	return &PeriodRepository{DB: db, Logger: loggerInstance}
}

func arrayToDomainPeriodMapper(rows *[]Period) *[]domainCatalog.Period {
	var domainRows []domainCatalog.Period
	for _, c := range *rows {
		domainRows = append(domainRows, *c.toDomainMapper())
	}
	return &domainRows
}

func (r *PeriodRepository) GetAll() (*[]domainCatalog.Period, error) {
	var results []Period
	if err := r.DB.Find(&results).Error; err != nil {
		r.Logger.Error("Error getting all periods", zap.Error(err))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	return arrayToDomainPeriodMapper(&results), nil
}

func (r *PeriodRepository) GetByID(id string) (*domainCatalog.Period, error) {
	var obj Period
	if err := r.DB.Where("id = ?", id).First(&obj).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainErrors.NewAppErrorWithType(domainErrors.NotFound)
		}
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	return obj.toDomainMapper(), nil
}

func (r *PeriodRepository) Create(newObj *domainCatalog.Period) (*domainCatalog.Period, error) {
	obj := fromDomainPeriodMapper(newObj)
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

func (r *PeriodRepository) Update(id string, objMap map[string]interface{}) (*domainCatalog.Period, error) {
	var obj Period
	if err := r.DB.Where("id = ?", id).First(&obj).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainErrors.NewAppErrorWithType(domainErrors.NotFound)
		}
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}

	updateData := make(map[string]interface{})
	for k, v := range objMap {
		if column, ok := ColumnsPeriodMapping[k]; ok {
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

func (r *PeriodRepository) Delete(id string) error {
	tx := r.DB.Where("id = ?", id).Delete(&Period{})
	if tx.Error != nil {
		return domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	if tx.RowsAffected == 0 {
		return domainErrors.NewAppErrorWithType(domainErrors.NotFound)
	}
	return nil
}

func (r *PeriodRepository) SearchPaginated(filters domain.DataFilters) (*domainCatalog.SearchResultPeriod, error) {
	query := r.DB.Model(&Period{})

	for field, values := range filters.LikeFilters {
		if len(values) > 0 {
			for _, value := range values {
				if value != "" && ColumnsPeriodMapping[field] != "" {
					query = query.Where(ColumnsPeriodMapping[field]+" LIKE ?", "%"+value+"%")
				}
			}
		}
	}
	for field, values := range filters.Matches {
		if len(values) > 0 && ColumnsPeriodMapping[field] != "" {
			query = query.Where(ColumnsPeriodMapping[field]+" IN ?", values)
		}
	}

	var count int64
	query.Count(&count)

	if len(filters.SortBy) > 0 && ColumnsPeriodMapping[filters.SortBy[0]] != "" {
		direction := "ASC"
		if filters.SortDirection == domain.SortDesc {
			direction = "DESC"
		}
		query = query.Order(ColumnsPeriodMapping[filters.SortBy[0]] + " " + direction)
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

	var list []Period
	if err := query.Find(&list).Error; err != nil {
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}

	return &domainCatalog.SearchResultPeriod{
		Data:       arrayToDomainPeriodMapper(&list),
		Total:      count,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: int(math.Ceil(float64(count) / float64(pageSize))),
	}, nil
}
