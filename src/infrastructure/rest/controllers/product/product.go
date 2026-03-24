package product

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gbrayhan/microservices-go/src/domain"
	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	domainProduct "github.com/gbrayhan/microservices-go/src/domain/product"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	repositoryProduct "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/product"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers"
	usecaseProduct "github.com/gbrayhan/microservices-go/src/application/usecases/product"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type NewProductRequest struct {
	Code            string  `json:"code" binding:"required"`
	Name            string  `json:"name" binding:"required"`
	Description     string  `json:"description"`
	EstimatedValue  float64 `json:"estimatedValue"`
	SellingPrice    float64 `json:"sellingPrice"`
	ConditionRating int     `json:"conditionRating"`
	Status          string  `json:"status"`
	CategoryID      *string `json:"categoryId"`
}

type ResponseProduct struct {
	ID              string  `json:"id"`
	Code            string  `json:"code"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	EstimatedValue  float64 `json:"estimatedValue"`
	SellingPrice    float64 `json:"sellingPrice"`
	ConditionRating int     `json:"conditionRating"`
	Status          string  `json:"status"`
	CategoryID      *string `json:"categoryId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type IProductController interface {
	NewProduct(ctx *gin.Context)
	GetAllProducts(ctx *gin.Context)
	GetProductByID(ctx *gin.Context)
	UpdateProduct(ctx *gin.Context)
	DeleteProduct(ctx *gin.Context)
	SearchPaginated(ctx *gin.Context)
}

type ProductController struct {
	productUseCase usecaseProduct.IProductUseCase
	Logger         *logger.Logger
}

func NewProductController(productUseCase usecaseProduct.IProductUseCase, loggerInstance *logger.Logger) IProductController {
	return &ProductController{productUseCase: productUseCase, Logger: loggerInstance}
}

func (c *ProductController) NewProduct(ctx *gin.Context) {
	c.Logger.Info("Creating new product")
	var request NewProductRequest
	if err := controllers.BindJSON(ctx, &request); err != nil {
		c.Logger.Error("Error binding JSON for new product", zap.Error(err))
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	productModel, err := c.productUseCase.Create(toUsecaseMapper(&request))
	if err != nil {
		c.Logger.Error("Error creating product", zap.Error(err), zap.String("code", request.Code))
		_ = ctx.Error(err)
		return
	}
	productResponse := domainToResponseMapper(productModel)
	ctx.JSON(http.StatusCreated, productResponse)
}

func (c *ProductController) GetAllProducts(ctx *gin.Context) {
	products, err := c.productUseCase.GetAll()
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, arrayDomainToResponseMapper(products))
}

func (c *ProductController) GetProductByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		appError := domainErrors.NewAppError(errors.New("product id is invalid"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	product, err := c.productUseCase.GetByID(id)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, domainToResponseMapper(product))
}

func (c *ProductController) UpdateProduct(ctx *gin.Context) {
	id := ctx.Param("id")
	var requestMap map[string]any
	err := controllers.BindJSONMap(ctx, &requestMap)
	if err != nil {
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	
	productUpdated, err := c.productUseCase.Update(id, requestMap)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, domainToResponseMapper(productUpdated))
}

func (c *ProductController) DeleteProduct(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.productUseCase.Delete(id)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "resource deleted successfully"})
}

func (c *ProductController) SearchPaginated(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))

	filters := domain.DataFilters{
		Page:     page,
		PageSize: pageSize,
	}

	likeFilters := make(map[string][]string)
	for field := range repositoryProduct.ColumnsProductMapping {
		if values := ctx.QueryArray(field + "_like"); len(values) > 0 {
			likeFilters[field] = values
		}
	}
	filters.LikeFilters = likeFilters

	result, err := c.productUseCase.SearchPaginated(filters)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	response := gin.H{
		"data":       arrayDomainToResponseMapper(result.Data),
		"total":      result.Total,
		"page":       result.Page,
		"pageSize":   result.PageSize,
		"totalPages": result.TotalPages,
		"filters":    filters,
	}
	ctx.JSON(http.StatusOK, response)
}

func domainToResponseMapper(p *domainProduct.Product) *ResponseProduct {
	return &ResponseProduct{
		ID:              p.ID,
		Code:            p.Code,
		Name:            p.Name,
		Description:     p.Description,
		EstimatedValue:  p.EstimatedValue,
		SellingPrice:    p.SellingPrice,
		ConditionRating: p.ConditionRating,
		Status:          p.Status,
		CategoryID:      p.CategoryID,
		CreatedAt:       p.CreatedAt,
		UpdatedAt:       p.UpdatedAt,
	}
}

func arrayDomainToResponseMapper(products *[]domainProduct.Product) *[]ResponseProduct {
	res := make([]ResponseProduct, len(*products))
	for i, p := range *products {
		res[i] = *domainToResponseMapper(&p)
	}
	return &res
}

func toUsecaseMapper(req *NewProductRequest) *domainProduct.Product {
	return &domainProduct.Product{
		Code:            req.Code,
		Name:            req.Name,
		Description:     req.Description,
		EstimatedValue:  req.EstimatedValue,
		SellingPrice:    req.SellingPrice,
		ConditionRating: req.ConditionRating,
		Status:          req.Status,
		CategoryID:      req.CategoryID,
	}
}
