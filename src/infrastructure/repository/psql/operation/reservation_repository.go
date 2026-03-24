package operation

import (
	"encoding/json"
	"math"

	"github.com/gbrayhan/microservices-go/src/domain"
	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	domainOperation "github.com/gbrayhan/microservices-go/src/domain/operation"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var ColumnsReservationMapping = map[string]string{
	"id":                    "id",
	"productId":             "product_id",
	"customerId":            "customer_id",
	"reservationDate":       "reservation_date",
	"expiryDate":            "expiry_date",
	"customReservationDays": "custom_reservation_days",
	"requiresDeposit":       "requires_deposit",
	"depositPercentage":     "deposit_percentage",
	"depositAmount":         "deposit_amount",
	"depositPaid":           "deposit_paid",
	"depositDate":           "deposit_date",
	"status":                "status",
	"cancellationReason":    "cancellation_reason",
	"notes":                 "notes",
	"createdBy":             "created_by",
	"createdAt":             "created_at",
	"updatedAt":             "updated_at",
}

type ReservationRepositoryInterface interface {
	GetAll() (*[]domainOperation.Reservation, error)
	GetByID(id string) (*domainOperation.Reservation, error)
	Create(objDomain *domainOperation.Reservation) (*domainOperation.Reservation, error)
	Update(id string, objMap map[string]interface{}) (*domainOperation.Reservation, error)
	Delete(id string) error
	SearchPaginated(filters domain.DataFilters) (*domainOperation.SearchResultReservation, error)
}

type ReservationRepository struct {
	DB     *gorm.DB
	Logger *logger.Logger
}

func NewReservationRepository(db *gorm.DB, loggerInstance *logger.Logger) ReservationRepositoryInterface {
	return &ReservationRepository{DB: db, Logger: loggerInstance}
}

func arrayToDomainReservationMapper(rows *[]Reservation) *[]domainOperation.Reservation {
	var domainRows []domainOperation.Reservation
	for _, c := range *rows {
		domainRows = append(domainRows, *c.toDomainMapper())
	}
	return &domainRows
}

func (r *ReservationRepository) GetAll() (*[]domainOperation.Reservation, error) {
	var results []Reservation
	if err := r.DB.Find(&results).Error; err != nil {
		r.Logger.Error("Error getting all reservations", zap.Error(err))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	return arrayToDomainReservationMapper(&results), nil
}

func (r *ReservationRepository) GetByID(id string) (*domainOperation.Reservation, error) {
	var obj Reservation
	if err := r.DB.Where("id = ?", id).First(&obj).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainErrors.NewAppErrorWithType(domainErrors.NotFound)
		}
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	return obj.toDomainMapper(), nil
}

func (r *ReservationRepository) Create(newObj *domainOperation.Reservation) (*domainOperation.Reservation, error) {
	obj := fromDomainReservationMapper(newObj)
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

func (r *ReservationRepository) Update(id string, objMap map[string]interface{}) (*domainOperation.Reservation, error) {
	var obj Reservation
	if err := r.DB.Where("id = ?", id).First(&obj).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainErrors.NewAppErrorWithType(domainErrors.NotFound)
		}
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}

	updateData := make(map[string]interface{})
	for k, v := range objMap {
		if column, ok := ColumnsReservationMapping[k]; ok {
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

func (r *ReservationRepository) Delete(id string) error {
	tx := r.DB.Where("id = ?", id).Delete(&Reservation{})
	if tx.Error != nil {
		return domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	if tx.RowsAffected == 0 {
		return domainErrors.NewAppErrorWithType(domainErrors.NotFound)
	}
	return nil
}

func (r *ReservationRepository) SearchPaginated(filters domain.DataFilters) (*domainOperation.SearchResultReservation, error) {
	query := r.DB.Model(&Reservation{})

	for field, values := range filters.LikeFilters {
		if len(values) > 0 {
			for _, value := range values {
				if value != "" && ColumnsReservationMapping[field] != "" {
					query = query.Where(ColumnsReservationMapping[field]+" LIKE ?", "%"+value+"%")
				}
			}
		}
	}
	for field, values := range filters.Matches {
		if len(values) > 0 && ColumnsReservationMapping[field] != "" {
			query = query.Where(ColumnsReservationMapping[field]+" IN ?", values)
		}
	}

	var count int64
	query.Count(&count)

	if len(filters.SortBy) > 0 && ColumnsReservationMapping[filters.SortBy[0]] != "" {
		direction := "ASC"
		if filters.SortDirection == domain.SortDesc {
			direction = "DESC"
		}
		query = query.Order(ColumnsReservationMapping[filters.SortBy[0]] + " " + direction)
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

	var list []Reservation
	if err := query.Find(&list).Error; err != nil {
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}

	return &domainOperation.SearchResultReservation{
		Data:       arrayToDomainReservationMapper(&list),
		Total:      count,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: int(math.Ceil(float64(count) / float64(pageSize))),
	}, nil
}
