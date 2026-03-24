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

type NewMaterialRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type ResponseMaterial struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type IMaterialController interface {
	NewMaterial(ctx *gin.Context)
	GetAllMaterials(ctx *gin.Context)
	GetMaterialByID(ctx *gin.Context)
	UpdateMaterial(ctx *gin.Context)
	DeleteMaterial(ctx *gin.Context)
	SearchPaginated(ctx *gin.Context)
}

type MaterialController struct {
	materialUseCase usecaseCatalog.IMaterialUseCase
	Logger          *logger.Logger
}

func NewMaterialController(materialUseCase usecaseCatalog.IMaterialUseCase, loggerInstance *logger.Logger) IMaterialController {
	return &MaterialController{materialUseCase: materialUseCase, Logger: loggerInstance}
}

func (c *MaterialController) NewMaterial(ctx *gin.Context) {
	c.Logger.Info("Creating new material")
	var request NewMaterialRequest
	if err := controllers.BindJSON(ctx, &request); err != nil {
		c.Logger.Error("Error binding JSON for new material", zap.Error(err))
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	model, err := c.materialUseCase.Create(toMaterialUsecaseMapper(&request))
	if err != nil {
		c.Logger.Error("Error creating material", zap.Error(err), zap.String("name", request.Name))
		_ = ctx.Error(err)
		return
	}
	response := domainToMaterialResponseMapper(model)
	ctx.JSON(http.StatusCreated, response)
}

func (c *MaterialController) GetAllMaterials(ctx *gin.Context) {
	items, err := c.materialUseCase.GetAll()
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, arrayDomainToMaterialResponseMapper(items))
}

func (c *MaterialController) GetMaterialByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		appError := domainErrors.NewAppError(errors.New("material id is invalid"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	item, err := c.materialUseCase.GetByID(id)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, domainToMaterialResponseMapper(item))
}

func (c *MaterialController) UpdateMaterial(ctx *gin.Context) {
	id := ctx.Param("id")
	var requestMap map[string]any
	err := controllers.BindJSONMap(ctx, &requestMap)
	if err != nil {
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	updated, err := c.materialUseCase.Update(id, requestMap)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, domainToMaterialResponseMapper(updated))
}

func (c *MaterialController) DeleteMaterial(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.materialUseCase.Delete(id)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "resource deleted successfully"})
}

func (c *MaterialController) SearchPaginated(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))

	filters := domain.DataFilters{
		Page:     page,
		PageSize: pageSize,
	}

	likeFilters := make(map[string][]string)
	for field := range repositoryCatalog.ColumnsMaterialMapping {
		if values := ctx.QueryArray(field + "_like"); len(values) > 0 {
			likeFilters[field] = values
		}
	}
	filters.LikeFilters = likeFilters

	result, err := c.materialUseCase.SearchPaginated(filters)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	response := gin.H{
		"data":       arrayDomainToMaterialResponseMapper(result.Data),
		"total":      result.Total,
		"page":       result.Page,
		"pageSize":   result.PageSize,
		"totalPages": result.TotalPages,
		"filters":    filters,
	}
	ctx.JSON(http.StatusOK, response)
}

func domainToMaterialResponseMapper(p *domainCatalog.Material) *ResponseMaterial {
	return &ResponseMaterial{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

func arrayDomainToMaterialResponseMapper(items *[]domainCatalog.Material) *[]ResponseMaterial {
	res := make([]ResponseMaterial, len(*items))
	for i, p := range *items {
		res[i] = *domainToMaterialResponseMapper(&p)
	}
	return &res
}

func toMaterialUsecaseMapper(req *NewMaterialRequest) *domainCatalog.Material {
	return &domainCatalog.Material{
		Name:        req.Name,
		Description: req.Description,
	}
}
