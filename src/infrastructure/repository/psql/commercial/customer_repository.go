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

var ColumnsCustomerMapping = map[string]string{
	"id":                     "id",
	"userId":                 "user_id",
	"customerType":           "customer_type",
	"preferences":            "preferences",
	"purchaseHistorySummary": "purchase_history_summary",
	"creditLimit":            "credit_limit",
	"notes":                  "notes",
	"createdAt":              "created_at",
	"updatedAt":              "updated_at",
}

type CustomerRepositoryInterface interface {
	GetAll() (*[]domainCommercial.Customer, error)
	GetByID(id string) (*domainCommercial.Customer, error)
	Create(objDomain *domainCommercial.Customer) (*domainCommercial.Customer, error)
	Update(id string, objMap map[string]interface{}) (*domainCommercial.Customer, error)
	Delete(id string) error
	SearchPaginated(filters domain.DataFilters) (*domainCommercial.SearchResultCustomer, error)
}

type CustomerRepository struct {
	DB     *gorm.DB
	Logger *logger.Logger
}

func NewCustomerRepository(db *gorm.DB, loggerInstance *logger.Logger) CustomerRepositoryInterface {
	return &CustomerRepository{DB: db, Logger: loggerInstance}
}

func arrayToDomainCustomerMapper(rows *[]Customer) *[]domainCommercial.Customer {
	var domainRows []domainCommercial.Customer
	for _, c := range *rows {
		domainRows = append(domainRows, *c.toDomainMapper())
	}
	return &domainRows
}

func (r *CustomerRepository) GetAll() (*[]domainCommercial.Customer, error) {
	var results []Customer
	if err := r.DB.Find(&results).Error; err != nil {
		r.Logger.Error("Error getting all customers", zap.Error(err))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	return arrayToDomainCustomerMapper(&results), nil
}

func (r *CustomerRepository) GetByID(id string) (*domainCommercial.Customer, error) {
	var obj Customer
	if err := r.DB.Where("id = ?", id).Preload("User").First(&obj).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainErrors.NewAppErrorWithType(domainErrors.NotFound)
		}
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	return obj.toDomainMapper(), nil
}

func (r *CustomerRepository) Create(newObj *domainCommercial.Customer) (*domainCommercial.Customer, error) {
	obj := fromDomainCustomerMapper(newObj)
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

func (r *CustomerRepository) Update(id string, objMap map[string]interface{}) (*domainCommercial.Customer, error) {
	var obj Customer
	if err := r.DB.Where("id = ?", id).First(&obj).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainErrors.NewAppErrorWithType(domainErrors.NotFound)
		}
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}

	updateData := make(map[string]interface{})
	for k, v := range objMap {
		if column, ok := ColumnsCustomerMapping[k]; ok {
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

func (r *CustomerRepository) Delete(id string) error {
	tx := r.DB.Where("id = ?", id).Delete(&Customer{})
	if tx.Error != nil {
		return domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	if tx.RowsAffected == 0 {
		return domainErrors.NewAppErrorWithType(domainErrors.NotFound)
	}
	return nil
}

func (r *CustomerRepository) SearchPaginated(filters domain.DataFilters) (*domainCommercial.SearchResultCustomer, error) {
	query := r.DB.Model(&Customer{})

	for field, values := range filters.LikeFilters {
		if len(values) > 0 {
			for _, value := range values {
				if value != "" && ColumnsCustomerMapping[field] != "" {
					query = query.Where(ColumnsCustomerMapping[field]+" LIKE ?", "%"+value+"%")
				}
			}
		}
	}
	for field, values := range filters.Matches {
		if len(values) > 0 && ColumnsCustomerMapping[field] != "" {
			query = query.Where(ColumnsCustomerMapping[field]+" IN ?", values)
		}
	}

	var count int64
	query.Count(&count)

	if len(filters.SortBy) > 0 && ColumnsCustomerMapping[filters.SortBy[0]] != "" {
		direction := "ASC"
		if filters.SortDirection == domain.SortDesc {
			direction = "DESC"
		}
		query = query.Order(ColumnsCustomerMapping[filters.SortBy[0]] + " " + direction)
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

	var list []Customer
	if err := query.Find(&list).Error; err != nil {
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}

	return &domainCommercial.SearchResultCustomer{
		Data:       arrayToDomainCustomerMapper(&list),
		Total:      count,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: int(math.Ceil(float64(count) / float64(pageSize))),
	}, nil
}
