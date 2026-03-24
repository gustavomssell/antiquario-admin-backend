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

var ColumnsProductMapping = map[string]string{
	"id":                "id",
	"code":              "code",
	"name":              "name",
	"categoryId":        "category_id",
	"status":            "status",
	"estimatedValue":    "estimated_value",
	"sellingPrice":      "selling_price",
	"createdAt":         "created_at",
	"updatedAt":         "updated_at",
}

type ProductRepositoryInterface interface {
	GetAll() (*[]domainProduct.Product, error)
	GetByID(id string) (*domainProduct.Product, error)
	Create(productDomain *domainProduct.Product) (*domainProduct.Product, error)
	Update(id string, productMap map[string]interface{}) (*domainProduct.Product, error)
	Delete(id string) error
	SearchPaginated(filters domain.DataFilters) (*domainProduct.SearchResultProduct, error)
}

type ProductRepository struct {
	DB     *gorm.DB
	Logger *logger.Logger
}

func NewProductRepository(db *gorm.DB, loggerInstance *logger.Logger) ProductRepositoryInterface {
	return &ProductRepository{DB: db, Logger: loggerInstance}
}

func arrayToDomainProductMapper(products *[]Product) *[]domainProduct.Product {
	var domainProducts []domainProduct.Product
	for _, p := range *products {
		domainProducts = append(domainProducts, *p.toDomainMapper())
	}
	return &domainProducts
}

func (r *ProductRepository) GetAll() (*[]domainProduct.Product, error) {
	var products []Product
	if err := r.DB.Find(&products).Error; err != nil {
		r.Logger.Error("Error getting all products", zap.Error(err))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	r.Logger.Info("Successfully retrieved all products", zap.Int("count", len(products)))
	return arrayToDomainProductMapper(&products), nil
}

func (r *ProductRepository) GetByID(id string) (*domainProduct.Product, error) {
	var productObj Product
	err := r.DB.Where("id = ?", id).First(&productObj).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			r.Logger.Warn("Product not found", zap.String("id", id))
			err = domainErrors.NewAppErrorWithType(domainErrors.NotFound)
		} else {
			r.Logger.Error("Error getting product by ID", zap.Error(err), zap.String("id", id))
			err = domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		}
		return nil, err
	}
	r.Logger.Info("Successfully retrieved product by ID", zap.String("id", id))
	return productObj.toDomainMapper(), nil
}

func (r *ProductRepository) Create(productDomain *domainProduct.Product) (*domainProduct.Product, error) {
	r.Logger.Info("Creating new product", zap.String("name", productDomain.Name))
	
	productObj := fromDomainProductMapper(productDomain)
	txDb := r.DB.Create(productObj)
	err := txDb.Error
	if err != nil {
		r.Logger.Error("Error creating product", zap.Error(err), zap.String("code", productDomain.Code))
		byteErr, _ := json.Marshal(err)
		var newError domainErrors.GormErr
		if errUnmarshal := json.Unmarshal(byteErr, &newError); errUnmarshal == nil {
			if newError.Number == 1062 || newError.Number == 23505 { // Unique constraint violation (23505 in PSQL)
				err = domainErrors.NewAppErrorWithType(domainErrors.ResourceAlreadyExists)
				return nil, err
			}
		}
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	r.Logger.Info("Successfully created product", zap.String("id", productObj.ID))
	return productObj.toDomainMapper(), nil
}

func (r *ProductRepository) Update(id string, productMap map[string]interface{}) (*domainProduct.Product, error) {
	var productObj Product
	err := r.DB.Where("id = ?", id).First(&productObj).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainErrors.NewAppErrorWithType(domainErrors.NotFound)
		}
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}

	updateData := make(map[string]interface{})
	for k, v := range productMap {
		if column, ok := ColumnsProductMapping[k]; ok {
			updateData[column] = v
		} else {
			updateData[k] = v // Fallback if key is exactly the db column or custom passing
		}
	}

	err = r.DB.Model(&productObj).Select("*").Updates(updateData).Error
	if err != nil {
		r.Logger.Error("Error updating product", zap.Error(err), zap.String("id", id))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}

	// Fetch updated object
	r.DB.Where("id = ?", id).First(&productObj)
	r.Logger.Info("Successfully updated product", zap.String("id", id))
	return productObj.toDomainMapper(), nil
}

func (r *ProductRepository) Delete(id string) error {
	tx := r.DB.Where("id = ?", id).Delete(&Product{})
	if tx.Error != nil {
		r.Logger.Error("Error deleting product", zap.Error(tx.Error), zap.String("id", id))
		return domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	if tx.RowsAffected == 0 {
		r.Logger.Warn("Product not found for deletion", zap.String("id", id))
		return domainErrors.NewAppErrorWithType(domainErrors.NotFound)
	}
	r.Logger.Info("Successfully deleted product", zap.String("id", id))
	return nil
}

func (r *ProductRepository) SearchPaginated(filters domain.DataFilters) (*domainProduct.SearchResultProduct, error) {
	query := r.DB.Model(&Product{})

	// Apply like filters
	for field, values := range filters.LikeFilters {
		if len(values) > 0 {
			for _, value := range values {
				if value != "" {
					column := ColumnsProductMapping[field]
					if column != "" {
						query = query.Where(column+" LIKE ?", "%"+value+"%")
					}
				}
			}
		}
	}

	// Apply equality filters
	for field, values := range filters.Matches {
		if len(values) > 0 {
			column := ColumnsProductMapping[field]
			if column != "" {
				query = query.Where(column+" IN ?", values)
			}
		}
	}

	// Pagination parameters
	var count int64
	query.Count(&count)

	if len(filters.SortBy) > 0 {
		sortColumn := ColumnsProductMapping[filters.SortBy[0]]
		if sortColumn != "" {
			direction := "ASC"
			if filters.SortDirection == domain.SortDesc {
				direction = "DESC"
			}
			query = query.Order(sortColumn + " " + direction)
		}
	} else {
		query = query.Order("created_at DESC")
	}

	page := 1
	if filters.Page > 0 {
		page = filters.Page
	}
	pageSize := 10
	if filters.PageSize > 0 {
		pageSize = filters.PageSize
	}

	query = query.Offset((page - 1) * pageSize).Limit(pageSize)

	var products []Product
	if err := query.Find(&products).Error; err != nil {
		r.Logger.Error("Error searching paginated products", zap.Error(err))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}

	totalPages := int(math.Ceil(float64(count) / float64(pageSize)))

	return &domainProduct.SearchResultProduct{
		Data:       arrayToDomainProductMapper(&products),
		Total:      count,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}
