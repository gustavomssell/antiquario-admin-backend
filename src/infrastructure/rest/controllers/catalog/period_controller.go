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

type NewPeriodRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	StartYear   *int   `json:"startYear"`
	EndYear     *int   `json:"endYear"`
}

type ResponsePeriod struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	StartYear   *int      `json:"startYear"`
	EndYear     *int      `json:"endYear"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type IPeriodController interface {
	NewPeriod(ctx *gin.Context)
	GetAllPeriods(ctx *gin.Context)
	GetPeriodByID(ctx *gin.Context)
	UpdatePeriod(ctx *gin.Context)
	DeletePeriod(ctx *gin.Context)
	SearchPaginated(ctx *gin.Context)
}

type PeriodController struct {
	periodUseCase usecaseCatalog.IPeriodUseCase
	Logger        *logger.Logger
}

func NewPeriodController(periodUseCase usecaseCatalog.IPeriodUseCase, loggerInstance *logger.Logger) IPeriodController {
	return &PeriodController{periodUseCase: periodUseCase, Logger: loggerInstance}
}

func (c *PeriodController) NewPeriod(ctx *gin.Context) {
	c.Logger.Info("Creating new period")
	var request NewPeriodRequest
	if err := controllers.BindJSON(ctx, &request); err != nil {
		c.Logger.Error("Error binding JSON for new period", zap.Error(err))
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	model, err := c.periodUseCase.Create(toPeriodUsecaseMapper(&request))
	if err != nil {
		c.Logger.Error("Error creating period", zap.Error(err), zap.String("name", request.Name))
		_ = ctx.Error(err)
		return
	}
	response := domainToPeriodResponseMapper(model)
	ctx.JSON(http.StatusCreated, response)
}

func (c *PeriodController) GetAllPeriods(ctx *gin.Context) {
	items, err := c.periodUseCase.GetAll()
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, arrayDomainToPeriodResponseMapper(items))
}

func (c *PeriodController) GetPeriodByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		appError := domainErrors.NewAppError(errors.New("period id is invalid"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	item, err := c.periodUseCase.GetByID(id)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, domainToPeriodResponseMapper(item))
}

func (c *PeriodController) UpdatePeriod(ctx *gin.Context) {
	id := ctx.Param("id")
	var requestMap map[string]any
	err := controllers.BindJSONMap(ctx, &requestMap)
	if err != nil {
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	updated, err := c.periodUseCase.Update(id, requestMap)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, domainToPeriodResponseMapper(updated))
}

func (c *PeriodController) DeletePeriod(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.periodUseCase.Delete(id)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "resource deleted successfully"})
}

func (c *PeriodController) SearchPaginated(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))

	filters := domain.DataFilters{
		Page:     page,
		PageSize: pageSize,
	}

	likeFilters := make(map[string][]string)
	for field := range repositoryCatalog.ColumnsPeriodMapping {
		if values := ctx.QueryArray(field + "_like"); len(values) > 0 {
			likeFilters[field] = values
		}
	}
	filters.LikeFilters = likeFilters

	result, err := c.periodUseCase.SearchPaginated(filters)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	response := gin.H{
		"data":       arrayDomainToPeriodResponseMapper(result.Data),
		"total":      result.Total,
		"page":       result.Page,
		"pageSize":   result.PageSize,
		"totalPages": result.TotalPages,
		"filters":    filters,
	}
	ctx.JSON(http.StatusOK, response)
}

func domainToPeriodResponseMapper(p *domainCatalog.Period) *ResponsePeriod {
	return &ResponsePeriod{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		StartYear:   p.StartYear,
		EndYear:     p.EndYear,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

func arrayDomainToPeriodResponseMapper(items *[]domainCatalog.Period) *[]ResponsePeriod {
	res := make([]ResponsePeriod, len(*items))
	for i, p := range *items {
		res[i] = *domainToPeriodResponseMapper(&p)
	}
	return &res
}

func toPeriodUsecaseMapper(req *NewPeriodRequest) *domainCatalog.Period {
	return &domainCatalog.Period{
		Name:        req.Name,
		Description: req.Description,
		StartYear:   req.StartYear,
		EndYear:     req.EndYear,
	}
}
