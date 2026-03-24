package commercial

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gbrayhan/microservices-go/src/domain"
	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	domainCommercial "github.com/gbrayhan/microservices-go/src/domain/commercial"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	repositoryCommercial "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/commercial"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers"
	usecaseCommercial "github.com/gbrayhan/microservices-go/src/application/usecases/commercial"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type NewAcquisitionRequest struct {
	SupplierID      string     `json:"supplierId" binding:"required"`
	AcquisitionDate *time.Time `json:"acquisitionDate"`
	TotalValue      float64    `json:"totalValue"`
	PaymentMethod   string     `json:"paymentMethod"`
	Notes           string     `json:"notes"`
	CreatedBy       int        `json:"createdBy"`
}

type ResponseAcquisition struct {
	ID              string     `json:"id"`
	SupplierID      string     `json:"supplierId"`
	AcquisitionDate *time.Time `json:"acquisitionDate"`
	TotalValue      float64    `json:"totalValue"`
	PaymentMethod   string     `json:"paymentMethod"`
	Notes           string     `json:"notes"`
	CreatedBy       int        `json:"createdBy"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
}

type IAcquisitionController interface {
	NewAcquisition(ctx *gin.Context)
	GetAllAcquisitions(ctx *gin.Context)
	GetAcquisitionByID(ctx *gin.Context)
	UpdateAcquisition(ctx *gin.Context)
	DeleteAcquisition(ctx *gin.Context)
	SearchPaginated(ctx *gin.Context)
}

type AcquisitionController struct {
	acquisitionUseCase usecaseCommercial.IAcquisitionUseCase
	Logger             *logger.Logger
}

func NewAcquisitionController(acquisitionUseCase usecaseCommercial.IAcquisitionUseCase, loggerInstance *logger.Logger) IAcquisitionController {
	return &AcquisitionController{acquisitionUseCase: acquisitionUseCase, Logger: loggerInstance}
}

func (c *AcquisitionController) NewAcquisition(ctx *gin.Context) {
	c.Logger.Info("Creating new acquisition")
	var request NewAcquisitionRequest
	if err := controllers.BindJSON(ctx, &request); err != nil {
		c.Logger.Error("Error binding JSON for new acquisition", zap.Error(err))
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	model, err := c.acquisitionUseCase.Create(toAcquisitionUsecaseMapper(&request))
	if err != nil {
		c.Logger.Error("Error creating acquisition", zap.Error(err), zap.String("supplierId", request.SupplierID))
		_ = ctx.Error(err)
		return
	}
	response := domainToAcquisitionResponseMapper(model)
	ctx.JSON(http.StatusCreated, response)
}

func (c *AcquisitionController) GetAllAcquisitions(ctx *gin.Context) {
	items, err := c.acquisitionUseCase.GetAll()
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, arrayDomainToAcquisitionResponseMapper(items))
}

func (c *AcquisitionController) GetAcquisitionByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		appError := domainErrors.NewAppError(errors.New("acquisition id is invalid"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	item, err := c.acquisitionUseCase.GetByID(id)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, domainToAcquisitionResponseMapper(item))
}

func (c *AcquisitionController) UpdateAcquisition(ctx *gin.Context) {
	id := ctx.Param("id")
	var requestMap map[string]any
	err := controllers.BindJSONMap(ctx, &requestMap)
	if err != nil {
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	updated, err := c.acquisitionUseCase.Update(id, requestMap)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, domainToAcquisitionResponseMapper(updated))
}

func (c *AcquisitionController) DeleteAcquisition(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.acquisitionUseCase.Delete(id)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "resource deleted successfully"})
}

func (c *AcquisitionController) SearchPaginated(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))

	filters := domain.DataFilters{
		Page:     page,
		PageSize: pageSize,
	}

	likeFilters := make(map[string][]string)
	for field := range repositoryCommercial.ColumnsAcquisitionMapping {
		if values := ctx.QueryArray(field + "_like"); len(values) > 0 {
			likeFilters[field] = values
		}
	}
	filters.LikeFilters = likeFilters

	result, err := c.acquisitionUseCase.SearchPaginated(filters)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	response := gin.H{
		"data":       arrayDomainToAcquisitionResponseMapper(result.Data),
		"total":      result.Total,
		"page":       result.Page,
		"pageSize":   result.PageSize,
		"totalPages": result.TotalPages,
		"filters":    filters,
	}
	ctx.JSON(http.StatusOK, response)
}

func domainToAcquisitionResponseMapper(p *domainCommercial.Acquisition) *ResponseAcquisition {
	return &ResponseAcquisition{
		ID:              p.ID,
		SupplierID:      p.SupplierID,
		AcquisitionDate: p.AcquisitionDate,
		TotalValue:      p.TotalValue,
		PaymentMethod:   p.PaymentMethod,
		Notes:           p.Notes,
		CreatedBy:       p.CreatedBy,
		CreatedAt:       p.CreatedAt,
		UpdatedAt:       p.UpdatedAt,
	}
}

func arrayDomainToAcquisitionResponseMapper(items *[]domainCommercial.Acquisition) *[]ResponseAcquisition {
	res := make([]ResponseAcquisition, len(*items))
	for i, p := range *items {
		res[i] = *domainToAcquisitionResponseMapper(&p)
	}
	return &res
}

func toAcquisitionUsecaseMapper(req *NewAcquisitionRequest) *domainCommercial.Acquisition {
	return &domainCommercial.Acquisition{
		SupplierID:      req.SupplierID,
		AcquisitionDate: req.AcquisitionDate,
		TotalValue:      req.TotalValue,
		PaymentMethod:   req.PaymentMethod,
		Notes:           req.Notes,
		CreatedBy:       req.CreatedBy,
	}
}
