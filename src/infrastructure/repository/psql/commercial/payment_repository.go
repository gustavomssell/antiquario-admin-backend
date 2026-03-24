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

var ColumnsPaymentMapping = map[string]string{
	"id":            "id",
	"saleId":        "sale_id",
	"amount":        "amount",
	"paymentDate":   "payment_date",
	"paymentMethod": "payment_method",
	"status":        "status",
	"reference":     "reference",
	"createdAt":     "created_at",
	"updatedAt":     "updated_at",
}

type PaymentRepositoryInterface interface {
	GetAll() (*[]domainCommercial.Payment, error)
	GetByID(id string) (*domainCommercial.Payment, error)
	Create(objDomain *domainCommercial.Payment) (*domainCommercial.Payment, error)
	Update(id string, objMap map[string]interface{}) (*domainCommercial.Payment, error)
	Delete(id string) error
	SearchPaginated(filters domain.DataFilters) (*domainCommercial.SearchResultPayment, error)
}

type PaymentRepository struct {
	DB     *gorm.DB
	Logger *logger.Logger
}

func NewPaymentRepository(db *gorm.DB, loggerInstance *logger.Logger) PaymentRepositoryInterface {
	return &PaymentRepository{DB: db, Logger: loggerInstance}
}

func arrayToDomainPaymentMapper(rows *[]Payment) *[]domainCommercial.Payment {
	var domainRows []domainCommercial.Payment
	for _, c := range *rows {
		domainRows = append(domainRows, *c.toDomainMapper())
	}
	return &domainRows
}

func (r *PaymentRepository) GetAll() (*[]domainCommercial.Payment, error) {
	var results []Payment
	if err := r.DB.Find(&results).Error; err != nil {
		r.Logger.Error("Error getting all payments", zap.Error(err))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	return arrayToDomainPaymentMapper(&results), nil
}

func (r *PaymentRepository) GetByID(id string) (*domainCommercial.Payment, error) {
	var obj Payment
	if err := r.DB.Where("id = ?", id).Preload("Sale").First(&obj).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainErrors.NewAppErrorWithType(domainErrors.NotFound)
		}
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	return obj.toDomainMapper(), nil
}

func (r *PaymentRepository) Create(newObj *domainCommercial.Payment) (*domainCommercial.Payment, error) {
	obj := fromDomainPaymentMapper(newObj)
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

func (r *PaymentRepository) Update(id string, objMap map[string]interface{}) (*domainCommercial.Payment, error) {
	var obj Payment
	if err := r.DB.Where("id = ?", id).First(&obj).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainErrors.NewAppErrorWithType(domainErrors.NotFound)
		}
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}

	updateData := make(map[string]interface{})
	for k, v := range objMap {
		if column, ok := ColumnsPaymentMapping[k]; ok {
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

func (r *PaymentRepository) Delete(id string) error {
	tx := r.DB.Where("id = ?", id).Delete(&Payment{})
	if tx.Error != nil {
		return domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	if tx.RowsAffected == 0 {
		return domainErrors.NewAppErrorWithType(domainErrors.NotFound)
	}
	return nil
}

func (r *PaymentRepository) SearchPaginated(filters domain.DataFilters) (*domainCommercial.SearchResultPayment, error) {
	query := r.DB.Model(&Payment{})

	for field, values := range filters.LikeFilters {
		if len(values) > 0 {
			for _, value := range values {
				if value != "" && ColumnsPaymentMapping[field] != "" {
					query = query.Where(ColumnsPaymentMapping[field]+" LIKE ?", "%"+value+"%")
				}
			}
		}
	}
	for field, values := range filters.Matches {
		if len(values) > 0 && ColumnsPaymentMapping[field] != "" {
			query = query.Where(ColumnsPaymentMapping[field]+" IN ?", values)
		}
	}

	var count int64
	query.Count(&count)

	if len(filters.SortBy) > 0 && ColumnsPaymentMapping[filters.SortBy[0]] != "" {
		direction := "ASC"
		if filters.SortDirection == domain.SortDesc {
			direction = "DESC"
		}
		query = query.Order(ColumnsPaymentMapping[filters.SortBy[0]] + " " + direction)
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

	var list []Payment
	if err := query.Find(&list).Error; err != nil {
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}

	return &domainCommercial.SearchResultPayment{
		Data:       arrayToDomainPaymentMapper(&list),
		Total:      count,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: int(math.Ceil(float64(count) / float64(pageSize))),
	}, nil
}
