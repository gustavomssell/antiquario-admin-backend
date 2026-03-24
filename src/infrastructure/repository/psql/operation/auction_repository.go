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

var ColumnsAuctionMapping = map[string]string{
	"id":          "id",
	"title":       "title",
	"description": "description",
	"startDate":   "start_date",
	"endDate":     "end_date",
	"status":      "status",
	"createdBy":   "created_by",
	"createdAt":   "created_at",
	"updatedAt":   "updated_at",
}

type AuctionRepositoryInterface interface {
	GetAll() (*[]domainOperation.Auction, error)
	GetByID(id string) (*domainOperation.Auction, error)
	Create(objDomain *domainOperation.Auction) (*domainOperation.Auction, error)
	Update(id string, objMap map[string]interface{}) (*domainOperation.Auction, error)
	Delete(id string) error
	SearchPaginated(filters domain.DataFilters) (*domainOperation.SearchResultAuction, error)
}

type AuctionRepository struct {
	DB     *gorm.DB
	Logger *logger.Logger
}

func NewAuctionRepository(db *gorm.DB, loggerInstance *logger.Logger) AuctionRepositoryInterface {
	return &AuctionRepository{DB: db, Logger: loggerInstance}
}

func arrayToDomainAuctionMapper(rows *[]Auction) *[]domainOperation.Auction {
	var domainRows []domainOperation.Auction
	for _, c := range *rows {
		domainRows = append(domainRows, *c.toDomainMapper())
	}
	return &domainRows
}

func (r *AuctionRepository) GetAll() (*[]domainOperation.Auction, error) {
	var results []Auction
	if err := r.DB.Find(&results).Error; err != nil {
		r.Logger.Error("Error getting all auctions", zap.Error(err))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	return arrayToDomainAuctionMapper(&results), nil
}

func (r *AuctionRepository) GetByID(id string) (*domainOperation.Auction, error) {
	var obj Auction
	if err := r.DB.Where("id = ?", id).First(&obj).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainErrors.NewAppErrorWithType(domainErrors.NotFound)
		}
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	return obj.toDomainMapper(), nil
}

func (r *AuctionRepository) Create(newObj *domainOperation.Auction) (*domainOperation.Auction, error) {
	obj := fromDomainAuctionMapper(newObj)
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

func (r *AuctionRepository) Update(id string, objMap map[string]interface{}) (*domainOperation.Auction, error) {
	var obj Auction
	if err := r.DB.Where("id = ?", id).First(&obj).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainErrors.NewAppErrorWithType(domainErrors.NotFound)
		}
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}

	updateData := make(map[string]interface{})
	for k, v := range objMap {
		if column, ok := ColumnsAuctionMapping[k]; ok {
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

func (r *AuctionRepository) Delete(id string) error {
	tx := r.DB.Where("id = ?", id).Delete(&Auction{})
	if tx.Error != nil {
		return domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	if tx.RowsAffected == 0 {
		return domainErrors.NewAppErrorWithType(domainErrors.NotFound)
	}
	return nil
}

func (r *AuctionRepository) SearchPaginated(filters domain.DataFilters) (*domainOperation.SearchResultAuction, error) {
	query := r.DB.Model(&Auction{})

	for field, values := range filters.LikeFilters {
		if len(values) > 0 {
			for _, value := range values {
				if value != "" && ColumnsAuctionMapping[field] != "" {
					query = query.Where(ColumnsAuctionMapping[field]+" LIKE ?", "%"+value+"%")
				}
			}
		}
	}
	for field, values := range filters.Matches {
		if len(values) > 0 && ColumnsAuctionMapping[field] != "" {
			query = query.Where(ColumnsAuctionMapping[field]+" IN ?", values)
		}
	}

	var count int64
	query.Count(&count)

	if len(filters.SortBy) > 0 && ColumnsAuctionMapping[filters.SortBy[0]] != "" {
		direction := "ASC"
		if filters.SortDirection == domain.SortDesc {
			direction = "DESC"
		}
		query = query.Order(ColumnsAuctionMapping[filters.SortBy[0]] + " " + direction)
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

	var list []Auction
	if err := query.Find(&list).Error; err != nil {
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}

	return &domainOperation.SearchResultAuction{
		Data:       arrayToDomainAuctionMapper(&list),
		Total:      count,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: int(math.Ceil(float64(count) / float64(pageSize))),
	}, nil
}
