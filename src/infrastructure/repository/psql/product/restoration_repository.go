package product

import (
	"encoding/json"
	"math"

	"github.com/gbrayhan/microservices-go/src/domain"
	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	domainProduct "github.com/gbrayhan/microservices-go/src/domain/product"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var ColumnsRestorationOrderMapping = map[string]string{
	"id":             "id",
	"productId":      "product_id",
	"status":         "status",
	"priority":       "priority",
	"estimatedCost":  "estimated_cost",
	"actualCost":     "actual_cost",
	"assignedTo":     "assigned_to",
	"startDate":      "start_date",
	"completionDate": "completion_date",
	"createdBy":      "created_by",
	"createdAt":      "created_at",
	"updatedAt":      "updated_at",
}

type RestorationOrderRepositoryInterface interface {
	GetAll() (*[]domainProduct.RestorationOrder, error)
	GetByID(id string) (*domainProduct.RestorationOrder, error)
	Create(objDomain *domainProduct.RestorationOrder) (*domainProduct.RestorationOrder, error)
	Update(id string, objMap map[string]interface{}) (*domainProduct.RestorationOrder, error)
	Delete(id string) error
	SearchPaginated(filters domain.DataFilters) (*domainProduct.SearchResultRestorationOrder, error)
}

type RestorationOrderRepository struct {
	DB     *gorm.DB
	Logger *logger.Logger
}

func NewRestorationOrderRepository(db *gorm.DB, loggerInstance *logger.Logger) RestorationOrderRepositoryInterface {
	return &RestorationOrderRepository{DB: db, Logger: loggerInstance}
}

func arrayToDomainRestorationOrderMapper(rows *[]RestorationOrder) *[]domainProduct.RestorationOrder {
	var domainRows []domainProduct.RestorationOrder
	for _, c := range *rows {
		domainRows = append(domainRows, *c.toDomainMapper())
	}
	return &domainRows
}

func (r *RestorationOrderRepository) GetAll() (*[]domainProduct.RestorationOrder, error) {
	var results []RestorationOrder
	if err := r.DB.Find(&results).Error; err != nil {
		r.Logger.Error("Error getting all restoration orders", zap.Error(err))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	return arrayToDomainRestorationOrderMapper(&results), nil
}

func (r *RestorationOrderRepository) GetByID(id string) (*domainProduct.RestorationOrder, error) {
	var obj RestorationOrder
	if err := r.DB.Where("id = ?", id).First(&obj).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainErrors.NewAppErrorWithType(domainErrors.NotFound)
		}
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	return obj.toDomainMapper(), nil
}

func (r *RestorationOrderRepository) Create(newObj *domainProduct.RestorationOrder) (*domainProduct.RestorationOrder, error) {
	obj := fromDomainRestorationOrderMapper(newObj)
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

func (r *RestorationOrderRepository) Update(id string, objMap map[string]interface{}) (*domainProduct.RestorationOrder, error) {
	var obj RestorationOrder
	if err := r.DB.Where("id = ?", id).First(&obj).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainErrors.NewAppErrorWithType(domainErrors.NotFound)
		}
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}

	updateData := make(map[string]interface{})
	for k, v := range objMap {
		if column, ok := ColumnsRestorationOrderMapping[k]; ok {
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

func (r *RestorationOrderRepository) Delete(id string) error {
	tx := r.DB.Where("id = ?", id).Delete(&RestorationOrder{})
	if tx.Error != nil {
		return domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	if tx.RowsAffected == 0 {
		return domainErrors.NewAppErrorWithType(domainErrors.NotFound)
	}
	return nil
}

func (r *RestorationOrderRepository) SearchPaginated(filters domain.DataFilters) (*domainProduct.SearchResultRestorationOrder, error) {
	query := r.DB.Model(&RestorationOrder{})

	for field, values := range filters.LikeFilters {
		if len(values) > 0 {
			for _, value := range values {
				if value != "" && ColumnsRestorationOrderMapping[field] != "" {
					query = query.Where(ColumnsRestorationOrderMapping[field]+" LIKE ?", "%"+value+"%")
				}
			}
		}
	}
	for field, values := range filters.Matches {
		if len(values) > 0 && ColumnsRestorationOrderMapping[field] != "" {
			query = query.Where(ColumnsRestorationOrderMapping[field]+" IN ?", values)
		}
	}

	var count int64
	query.Count(&count)

	if len(filters.SortBy) > 0 && ColumnsRestorationOrderMapping[filters.SortBy[0]] != "" {
		direction := "ASC"
		if filters.SortDirection == domain.SortDesc {
			direction = "DESC"
		}
		query = query.Order(ColumnsRestorationOrderMapping[filters.SortBy[0]] + " " + direction)
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

	var list []RestorationOrder
	if err := query.Find(&list).Error; err != nil {
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}

	return &domainProduct.SearchResultRestorationOrder{
		Data:       arrayToDomainRestorationOrderMapper(&list),
		Total:      count,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: int(math.Ceil(float64(count) / float64(pageSize))),
	}, nil
}
