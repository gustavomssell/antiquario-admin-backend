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

var ColumnsTagMapping = map[string]string{
	"id":          "id",
	"name":        "name",
	"description": "description",
	"createdAt":   "created_at",
	"updatedAt":   "updated_at",
}

type TagRepositoryInterface interface {
	GetAll() (*[]domainCatalog.Tag, error)
	GetByID(id string) (*domainCatalog.Tag, error)
	Create(objDomain *domainCatalog.Tag) (*domainCatalog.Tag, error)
	Update(id string, objMap map[string]interface{}) (*domainCatalog.Tag, error)
	Delete(id string) error
	SearchPaginated(filters domain.DataFilters) (*domainCatalog.SearchResultTag, error)
}

type TagRepository struct {
	DB     *gorm.DB
	Logger *logger.Logger
}

func NewTagRepository(db *gorm.DB, loggerInstance *logger.Logger) TagRepositoryInterface {
	return &TagRepository{DB: db, Logger: loggerInstance}
}

func arrayToDomainTagMapper(rows *[]Tag) *[]domainCatalog.Tag {
	var domainRows []domainCatalog.Tag
	for _, c := range *rows {
		domainRows = append(domainRows, *c.toDomainMapper())
	}
	return &domainRows
}

func (r *TagRepository) GetAll() (*[]domainCatalog.Tag, error) {
	var results []Tag
	if err := r.DB.Find(&results).Error; err != nil {
		r.Logger.Error("Error getting all tags", zap.Error(err))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	return arrayToDomainTagMapper(&results), nil
}

func (r *TagRepository) GetByID(id string) (*domainCatalog.Tag, error) {
	var obj Tag
	if err := r.DB.Where("id = ?", id).First(&obj).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainErrors.NewAppErrorWithType(domainErrors.NotFound)
		}
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	return obj.toDomainMapper(), nil
}

func (r *TagRepository) Create(newObj *domainCatalog.Tag) (*domainCatalog.Tag, error) {
	obj := fromDomainTagMapper(newObj)
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

func (r *TagRepository) Update(id string, objMap map[string]interface{}) (*domainCatalog.Tag, error) {
	var obj Tag
	if err := r.DB.Where("id = ?", id).First(&obj).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainErrors.NewAppErrorWithType(domainErrors.NotFound)
		}
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}

	updateData := make(map[string]interface{})
	for k, v := range objMap {
		if column, ok := ColumnsTagMapping[k]; ok {
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

func (r *TagRepository) Delete(id string) error {
	tx := r.DB.Where("id = ?", id).Delete(&Tag{})
	if tx.Error != nil {
		return domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	if tx.RowsAffected == 0 {
		return domainErrors.NewAppErrorWithType(domainErrors.NotFound)
	}
	return nil
}

func (r *TagRepository) SearchPaginated(filters domain.DataFilters) (*domainCatalog.SearchResultTag, error) {
	query := r.DB.Model(&Tag{})

	for field, values := range filters.LikeFilters {
		if len(values) > 0 {
			for _, value := range values {
				if value != "" && ColumnsTagMapping[field] != "" {
					query = query.Where(ColumnsTagMapping[field]+" LIKE ?", "%"+value+"%")
				}
			}
		}
	}
	for field, values := range filters.Matches {
		if len(values) > 0 && ColumnsTagMapping[field] != "" {
			query = query.Where(ColumnsTagMapping[field]+" IN ?", values)
		}
	}

	var count int64
	query.Count(&count)

	if len(filters.SortBy) > 0 && ColumnsTagMapping[filters.SortBy[0]] != "" {
		direction := "ASC"
		if filters.SortDirection == domain.SortDesc {
			direction = "DESC"
		}
		query = query.Order(ColumnsTagMapping[filters.SortBy[0]] + " " + direction)
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

	var list []Tag
	if err := query.Find(&list).Error; err != nil {
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}

	return &domainCatalog.SearchResultTag{
		Data:       arrayToDomainTagMapper(&list),
		Total:      count,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: int(math.Ceil(float64(count) / float64(pageSize))),
	}, nil
}
