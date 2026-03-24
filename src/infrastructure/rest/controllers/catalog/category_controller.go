package catalog

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gbrayhan/microservices-go/src/domain"
	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	domainCatalog "github.com/gbrayhan/microservices-go/src/domain/catalog"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	repositoryCatalog "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/catalog"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers"
	usecaseCatalog "github.com/gbrayhan/microservices-go/src/application/usecases/catalog"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type NewCategoryRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	ParentID    *string `json:"parentId"`
	Level       int     `json:"level"`
	Active      bool    `json:"active"`
}

type ResponseCategory struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ParentID    *string   `json:"parentId"`
	Level       int       `json:"level"`
	Active      bool      `json:"active"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type ICategoryController interface {
	NewCategory(ctx *gin.Context)
	GetAllCategories(ctx *gin.Context)
	GetCategoryByID(ctx *gin.Context)
	UpdateCategory(ctx *gin.Context)
	DeleteCategory(ctx *gin.Context)
	SearchPaginated(ctx *gin.Context)
}

type CategoryController struct {
	categoryUseCase usecaseCatalog.ICategoryUseCase
	Logger          *logger.Logger
}

func NewCategoryController(categoryUseCase usecaseCatalog.ICategoryUseCase, loggerInstance *logger.Logger) ICategoryController {
	return &CategoryController{categoryUseCase: categoryUseCase, Logger: loggerInstance}
}

func (c *CategoryController) NewCategory(ctx *gin.Context) {
	c.Logger.Info("Creating new category")
	var request NewCategoryRequest
	if err := controllers.BindJSON(ctx, &request); err != nil {
		c.Logger.Error("Error binding JSON for new category", zap.Error(err))
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	model, err := c.categoryUseCase.Create(toUsecaseMapper(&request))
	if err != nil {
		c.Logger.Error("Error creating category", zap.Error(err), zap.String("name", request.Name))
		_ = ctx.Error(err)
		return
	}
	response := domainToResponseMapper(model)
	ctx.JSON(http.StatusCreated, response)
}

func (c *CategoryController) GetAllCategories(ctx *gin.Context) {
	items, err := c.categoryUseCase.GetAll()
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, arrayDomainToResponseMapper(items))
}

func (c *CategoryController) GetCategoryByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		appError := domainErrors.NewAppError(errors.New("category id is invalid"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	item, err := c.categoryUseCase.GetByID(id)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, domainToResponseMapper(item))
}

func (c *CategoryController) UpdateCategory(ctx *gin.Context) {
	id := ctx.Param("id")
	var requestMap map[string]any
	err := controllers.BindJSONMap(ctx, &requestMap)
	if err != nil {
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	updated, err := c.categoryUseCase.Update(id, requestMap)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, domainToResponseMapper(updated))
}

func (c *CategoryController) DeleteCategory(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.categoryUseCase.Delete(id)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "resource deleted successfully"})
}

func (c *CategoryController) SearchPaginated(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))

	filters := domain.DataFilters{
		Page:     page,
		PageSize: pageSize,
	}

	likeFilters := make(map[string][]string)
	for field := range repositoryCatalog.ColumnsCategoryMapping {
		if values := ctx.QueryArray(field + "_like"); len(values) > 0 {
			likeFilters[field] = values
		}
	}
	filters.LikeFilters = likeFilters

	result, err := c.categoryUseCase.SearchPaginated(filters)
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

func domainToResponseMapper(p *domainCatalog.Category) *ResponseCategory {
	return &ResponseCategory{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		ParentID:    p.ParentID,
		Level:       p.Level,
		Active:      p.Active,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

func arrayDomainToResponseMapper(items *[]domainCatalog.Category) *[]ResponseCategory {
	res := make([]ResponseCategory, len(*items))
	for i, p := range *items {
		res[i] = *domainToResponseMapper(&p)
	}
	return &res
}

func toUsecaseMapper(req *NewCategoryRequest) *domainCatalog.Category {
	return &domainCatalog.Category{
		Name:        req.Name,
		Description: req.Description,
		ParentID:    req.ParentID,
		Level:       req.Level,
		Active:      req.Active,
	}
}
