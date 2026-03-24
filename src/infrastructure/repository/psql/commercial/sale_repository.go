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

var ColumnsSaleMapping = map[string]string{
	"id":              "id",
	"customerId":      "customer_id",
	"saleDate":        "sale_date",
	"saleType":        "sale_type",
	"totalAmount":     "total_amount",
	"paymentMethod":   "payment_method",
	"paymentStatus":   "payment_status",
	"deliveryAddress": "delivery_address",
	"deliveryDate":    "delivery_date",
	"notes":           "notes",
	"createdBy":       "created_by",
	"createdAt":       "created_at",
	"updatedAt":       "updated_at",
}

type SaleRepositoryInterface interface {
	GetAll() (*[]domainCommercial.Sale, error)
	GetByID(id string) (*domainCommercial.Sale, error)
	Create(objDomain *domainCommercial.Sale) (*domainCommercial.Sale, error)
	Update(id string, objMap map[string]interface{}) (*domainCommercial.Sale, error)
	Delete(id string) error
	SearchPaginated(filters domain.DataFilters) (*domainCommercial.SearchResultSale, error)
}

type SaleRepository struct {
	DB     *gorm.DB
	Logger *logger.Logger
}

func NewSaleRepository(db *gorm.DB, loggerInstance *logger.Logger) SaleRepositoryInterface {
	return &SaleRepository{DB: db, Logger: loggerInstance}
}

func arrayToDomainSaleMapper(rows *[]Sale) *[]domainCommercial.Sale {
	var domainRows []domainCommercial.Sale
	for _, c := range *rows {
		domainRows = append(domainRows, *c.toDomainMapper())
	}
	return &domainRows
}

func (r *SaleRepository) GetAll() (*[]domainCommercial.Sale, error) {
	var results []Sale
	if err := r.DB.Find(&results).Error; err != nil {
		r.Logger.Error("Error getting all sales", zap.Error(err))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	return arrayToDomainSaleMapper(&results), nil
}

func (r *SaleRepository) GetByID(id string) (*domainCommercial.Sale, error) {
	var obj Sale
	if err := r.DB.Where("id = ?", id).Preload("Customer").First(&obj).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainErrors.NewAppErrorWithType(domainErrors.NotFound)
		}
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	return obj.toDomainMapper(), nil
}

func (r *SaleRepository) Create(newObj *domainCommercial.Sale) (*domainCommercial.Sale, error) {
	obj := fromDomainSaleMapper(newObj)
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

func (r *SaleRepository) Update(id string, objMap map[string]interface{}) (*domainCommercial.Sale, error) {
	var obj Sale
	if err := r.DB.Where("id = ?", id).First(&obj).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainErrors.NewAppErrorWithType(domainErrors.NotFound)
		}
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}

	updateData := make(map[string]interface{})
	for k, v := range objMap {
		if column, ok := ColumnsSaleMapping[k]; ok {
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

func (r *SaleRepository) Delete(id string) error {
	tx := r.DB.Where("id = ?", id).Delete(&Sale{})
	if tx.Error != nil {
		return domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	if tx.RowsAffected == 0 {
		return domainErrors.NewAppErrorWithType(domainErrors.NotFound)
	}
	return nil
}

func (r *SaleRepository) SearchPaginated(filters domain.DataFilters) (*domainCommercial.SearchResultSale, error) {
	query := r.DB.Model(&Sale{})

	for field, values := range filters.LikeFilters {
		if len(values) > 0 {
			for _, value := range values {
				if value != "" && ColumnsSaleMapping[field] != "" {
					query = query.Where(ColumnsSaleMapping[field]+" LIKE ?", "%"+value+"%")
				}
			}
		}
	}
	for field, values := range filters.Matches {
		if len(values) > 0 && ColumnsSaleMapping[field] != "" {
			query = query.Where(ColumnsSaleMapping[field]+" IN ?", values)
		}
	}

	var count int64
	query.Count(&count)

	if len(filters.SortBy) > 0 && ColumnsSaleMapping[filters.SortBy[0]] != "" {
		direction := "ASC"
		if filters.SortDirection == domain.SortDesc {
			direction = "DESC"
		}
		query = query.Order(ColumnsSaleMapping[filters.SortBy[0]] + " " + direction)
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

	var list []Sale
	if err := query.Find(&list).Error; err != nil {
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}

	return &domainCommercial.SearchResultSale{
		Data:       arrayToDomainSaleMapper(&list),
		Total:      count,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: int(math.Ceil(float64(count) / float64(pageSize))),
	}, nil
}
