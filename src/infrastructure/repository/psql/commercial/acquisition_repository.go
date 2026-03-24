package commercial

import (
	"encoding/json"
	"math"

	"github.com/gbrayhan/microservices-go/src/domain"
	domainCommercial "github.com/gbrayhan/microservices-go/src/domain/commercial"
	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var ColumnsAcquisitionMapping = map[string]string{
	"id":              "id",
	"supplierId":      "supplier_id",
	"acquisitionDate": "acquisition_date",
	"totalValue":      "total_value",
	"paymentMethod":   "payment_method",
	"notes":           "notes",
	"createdBy":       "created_by",
	"createdAt":       "created_at",
	"updatedAt":       "updated_at",
}

type AcquisitionRepositoryInterface interface {
	GetAll() (*[]domainCommercial.Acquisition, error)
	GetByID(id string) (*domainCommercial.Acquisition, error)
	Create(objDomain *domainCommercial.Acquisition) (*domainCommercial.Acquisition, error)
	Update(id string, objMap map[string]interface{}) (*domainCommercial.Acquisition, error)
	Delete(id string) error
	SearchPaginated(filters domain.DataFilters) (*domainCommercial.SearchResultAcquisition, error)
}

type AcquisitionRepository struct {
	DB     *gorm.DB
	Logger *logger.Logger
}

func NewAcquisitionRepository(db *gorm.DB, loggerInstance *logger.Logger) AcquisitionRepositoryInterface {
	return &AcquisitionRepository{DB: db, Logger: loggerInstance}
}

func arrayToDomainAcquisitionMapper(rows *[]Acquisition) *[]domainCommercial.Acquisition {
	var domainRows []domainCommercial.Acquisition
	for _, c := range *rows {
		domainRows = append(domainRows, *c.toDomainMapper())
	}
	return &domainRows
}

func (r *AcquisitionRepository) GetAll() (*[]domainCommercial.Acquisition, error) {
	var results []Acquisition
	if err := r.DB.Find(&results).Error; err != nil {
		r.Logger.Error("Error getting all acquisitions", zap.Error(err))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	return arrayToDomainAcquisitionMapper(&results), nil
}

func (r *AcquisitionRepository) GetByID(id string) (*domainCommercial.Acquisition, error) {
	var obj Acquisition
	if err := r.DB.Where("id = ?", id).Preload("Supplier").First(&obj).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainErrors.NewAppErrorWithType(domainErrors.NotFound)
		}
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	return obj.toDomainMapper(), nil
}

func (r *AcquisitionRepository) Create(newObj *domainCommercial.Acquisition) (*domainCommercial.Acquisition, error) {
	obj := fromDomainAcquisitionMapper(newObj)
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

func (r *AcquisitionRepository) Update(id string, objMap map[string]interface{}) (*domainCommercial.Acquisition, error) {
	var obj Acquisition
	if err := r.DB.Where("id = ?", id).First(&obj).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainErrors.NewAppErrorWithType(domainErrors.NotFound)
		}
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}

	updateData := make(map[string]interface{})
	for k, v := range objMap {
		if column, ok := ColumnsAcquisitionMapping[k]; ok {
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

func (r *AcquisitionRepository) Delete(id string) error {
	tx := r.DB.Where("id = ?", id).Delete(&Acquisition{})
	if tx.Error != nil {
		return domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	if tx.RowsAffected == 0 {
		return domainErrors.NewAppErrorWithType(domainErrors.NotFound)
	}
	return nil
}

func (r *AcquisitionRepository) SearchPaginated(filters domain.DataFilters) (*domainCommercial.SearchResultAcquisition, error) {
	query := r.DB.Model(&Acquisition{})

	for field, values := range filters.LikeFilters {
		if len(values) > 0 {
			for _, value := range values {
				if value != "" && ColumnsAcquisitionMapping[field] != "" {
					query = query.Where(ColumnsAcquisitionMapping[field]+" LIKE ?", "%"+value+"%")
				}
			}
		}
	}
	for field, values := range filters.Matches {
		if len(values) > 0 && ColumnsAcquisitionMapping[field] != "" {
			query = query.Where(ColumnsAcquisitionMapping[field]+" IN ?", values)
		}
	}

	var count int64
	query.Count(&count)

	if len(filters.SortBy) > 0 && ColumnsAcquisitionMapping[filters.SortBy[0]] != "" {
		direction := "ASC"
		if filters.SortDirection == domain.SortDesc {
			direction = "DESC"
		}
		query = query.Order(ColumnsAcquisitionMapping[filters.SortBy[0]] + " " + direction)
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

	var list []Acquisition
	if err := query.Find(&list).Error; err != nil {
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}

	return &domainCommercial.SearchResultAcquisition{
		Data:       arrayToDomainAcquisitionMapper(&list),
		Total:      count,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: int(math.Ceil(float64(count) / float64(pageSize))),
	}, nil
}
