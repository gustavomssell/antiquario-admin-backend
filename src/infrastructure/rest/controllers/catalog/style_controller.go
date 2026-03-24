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

type NewStyleRequest struct {
	Name          string  `json:"name" binding:"required"`
	Description   string  `json:"description"`
	PeriodID      *string `json:"periodId"`
	OriginCountry string  `json:"originCountry"`
}

type ResponseStyle struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	PeriodID      *string   `json:"periodId"`
	OriginCountry string    `json:"originCountry"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type IStyleController interface {
	NewStyle(ctx *gin.Context)
	GetAllStyles(ctx *gin.Context)
	GetStyleByID(ctx *gin.Context)
	UpdateStyle(ctx *gin.Context)
	DeleteStyle(ctx *gin.Context)
	SearchPaginated(ctx *gin.Context)
}

type StyleController struct {
	styleUseCase usecaseCatalog.IStyleUseCase
	Logger       *logger.Logger
}

func NewStyleController(styleUseCase usecaseCatalog.IStyleUseCase, loggerInstance *logger.Logger) IStyleController {
	return &StyleController{styleUseCase: styleUseCase, Logger: loggerInstance}
}

func (c *StyleController) NewStyle(ctx *gin.Context) {
	c.Logger.Info("Creating new style")
	var request NewStyleRequest
	if err := controllers.BindJSON(ctx, &request); err != nil {
		c.Logger.Error("Error binding JSON for new style", zap.Error(err))
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	model, err := c.styleUseCase.Create(toStyleUsecaseMapper(&request))
	if err != nil {
		c.Logger.Error("Error creating style", zap.Error(err), zap.String("name", request.Name))
		_ = ctx.Error(err)
		return
	}
	response := domainToStyleResponseMapper(model)
	ctx.JSON(http.StatusCreated, response)
}

func (c *StyleController) GetAllStyles(ctx *gin.Context) {
	items, err := c.styleUseCase.GetAll()
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, arrayDomainToStyleResponseMapper(items))
}

func (c *StyleController) GetStyleByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		appError := domainErrors.NewAppError(errors.New("style id is invalid"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	item, err := c.styleUseCase.GetByID(id)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, domainToStyleResponseMapper(item))
}

func (c *StyleController) UpdateStyle(ctx *gin.Context) {
	id := ctx.Param("id")
	var requestMap map[string]any
	err := controllers.BindJSONMap(ctx, &requestMap)
	if err != nil {
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	updated, err := c.styleUseCase.Update(id, requestMap)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, domainToStyleResponseMapper(updated))
}

func (c *StyleController) DeleteStyle(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.styleUseCase.Delete(id)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "resource deleted successfully"})
}

func (c *StyleController) SearchPaginated(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))

	filters := domain.DataFilters{
		Page:     page,
		PageSize: pageSize,
	}

	likeFilters := make(map[string][]string)
	for field := range repositoryCatalog.ColumnsStyleMapping {
		if values := ctx.QueryArray(field + "_like"); len(values) > 0 {
			likeFilters[field] = values
		}
	}
	filters.LikeFilters = likeFilters

	result, err := c.styleUseCase.SearchPaginated(filters)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	response := gin.H{
		"data":       arrayDomainToStyleResponseMapper(result.Data),
		"total":      result.Total,
		"page":       result.Page,
		"pageSize":   result.PageSize,
		"totalPages": result.TotalPages,
		"filters":    filters,
	}
	ctx.JSON(http.StatusOK, response)
}

func domainToStyleResponseMapper(p *domainCatalog.Style) *ResponseStyle {
	return &ResponseStyle{
		ID:            p.ID,
		Name:          p.Name,
		Description:   p.Description,
		PeriodID:      p.PeriodID,
		OriginCountry: p.OriginCountry,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}
}

func arrayDomainToStyleResponseMapper(items *[]domainCatalog.Style) *[]ResponseStyle {
	res := make([]ResponseStyle, len(*items))
	for i, p := range *items {
		res[i] = *domainToStyleResponseMapper(&p)
	}
	return &res
}

func toStyleUsecaseMapper(req *NewStyleRequest) *domainCatalog.Style {
	return &domainCatalog.Style{
		Name:          req.Name,
		Description:   req.Description,
		PeriodID:      req.PeriodID,
		OriginCountry: req.OriginCountry,
	}
}
