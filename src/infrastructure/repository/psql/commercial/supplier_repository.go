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

var ColumnsSupplierMapping = map[string]string{
	"id":             "id",
	"name":           "name",
	"type":           "type",
	"contactInfo":    "contact_info",
	"address":        "address",
	"documentNumber": "document_number",
	"notes":          "notes",
	"createdAt":      "created_at",
	"updatedAt":      "updated_at",
}

type SupplierRepositoryInterface interface {
	GetAll() (*[]domainCommercial.Supplier, error)
	GetByID(id string) (*domainCommercial.Supplier, error)
	Create(objDomain *domainCommercial.Supplier) (*domainCommercial.Supplier, error)
	Update(id string, objMap map[string]interface{}) (*domainCommercial.Supplier, error)
	Delete(id string) error
	SearchPaginated(filters domain.DataFilters) (*domainCommercial.SearchResultSupplier, error)
}

type SupplierRepository struct {
	DB     *gorm.DB
	Logger *logger.Logger
}

func NewSupplierRepository(db *gorm.DB, loggerInstance *logger.Logger) SupplierRepositoryInterface {
	return &SupplierRepository{DB: db, Logger: loggerInstance}
}

func arrayToDomainSupplierMapper(rows *[]Supplier) *[]domainCommercial.Supplier {
	var domainRows []domainCommercial.Supplier
	for _, c := range *rows {
		domainRows = append(domainRows, *c.toDomainMapper())
	}
	return &domainRows
}

func (r *SupplierRepository) GetAll() (*[]domainCommercial.Supplier, error) {
	var results []Supplier
	if err := r.DB.Find(&results).Error; err != nil {
		r.Logger.Error("Error getting all suppliers", zap.Error(err))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	return arrayToDomainSupplierMapper(&results), nil
}

func (r *SupplierRepository) GetByID(id string) (*domainCommercial.Supplier, error) {
	var obj Supplier
	if err := r.DB.Where("id = ?", id).First(&obj).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainErrors.NewAppErrorWithType(domainErrors.NotFound)
		}
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	return obj.toDomainMapper(), nil
}

func (r *SupplierRepository) Create(newObj *domainCommercial.Supplier) (*domainCommercial.Supplier, error) {
	obj := fromDomainSupplierMapper(newObj)
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

func (r *SupplierRepository) Update(id string, objMap map[string]interface{}) (*domainCommercial.Supplier, error) {
	var obj Supplier
	if err := r.DB.Where("id = ?", id).First(&obj).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainErrors.NewAppErrorWithType(domainErrors.NotFound)
		}
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}

	updateData := make(map[string]interface{})
	for k, v := range objMap {
		if column, ok := ColumnsSupplierMapping[k]; ok {
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

func (r *SupplierRepository) Delete(id string) error {
	tx := r.DB.Where("id = ?", id).Delete(&Supplier{})
	if tx.Error != nil {
		return domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	if tx.RowsAffected == 0 {
		return domainErrors.NewAppErrorWithType(domainErrors.NotFound)
	}
	return nil
}

func (r *SupplierRepository) SearchPaginated(filters domain.DataFilters) (*domainCommercial.SearchResultSupplier, error) {
	query := r.DB.Model(&Supplier{})

	for field, values := range filters.LikeFilters {
		if len(values) > 0 {
			for _, value := range values {
				if value != "" && ColumnsSupplierMapping[field] != "" {
					query = query.Where(ColumnsSupplierMapping[field]+" LIKE ?", "%"+value+"%")
				}
			}
		}
	}
	for field, values := range filters.Matches {
		if len(values) > 0 && ColumnsSupplierMapping[field] != "" {
			query = query.Where(ColumnsSupplierMapping[field]+" IN ?", values)
		}
	}

	var count int64
	query.Count(&count)

	if len(filters.SortBy) > 0 && ColumnsSupplierMapping[filters.SortBy[0]] != "" {
		direction := "ASC"
		if filters.SortDirection == domain.SortDesc {
			direction = "DESC"
		}
		query = query.Order(ColumnsSupplierMapping[filters.SortBy[0]] + " " + direction)
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

	var list []Supplier
	if err := query.Find(&list).Error; err != nil {
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}

	return &domainCommercial.SearchResultSupplier{
		Data:       arrayToDomainSupplierMapper(&list),
		Total:      count,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: int(math.Ceil(float64(count) / float64(pageSize))),
	}, nil
}
