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

type NewRestorationOrderRequest struct {
	ProductID      string     `json:"productId" binding:"required"`
	Status         string     `json:"status"`
	Priority       string     `json:"priority"`
	EstimatedCost  float64    `json:"estimatedCost"`
	ActualCost     *float64   `json:"actualCost"`
	AssignedTo     *int       `json:"assignedTo"`
	StartDate      *time.Time `json:"startDate"`
	CompletionDate *time.Time `json:"completionDate"`
	CreatedBy      int        `json:"createdBy"`
}

type ResponseRestorationOrder struct {
	ID             string     `json:"id"`
	ProductID      string     `json:"productId"`
	Status         string     `json:"status"`
	Priority       string     `json:"priority"`
	EstimatedCost  float64    `json:"estimatedCost"`
	ActualCost     *float64   `json:"actualCost"`
	AssignedTo     *int       `json:"assignedTo"`
	StartDate      *time.Time `json:"startDate"`
	CompletionDate *time.Time `json:"completionDate"`
	CreatedBy      int        `json:"createdBy"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
}

type IRestorationOrderController interface {
	NewRestorationOrder(ctx *gin.Context)
	GetAllRestorationOrders(ctx *gin.Context)
	GetRestorationOrderByID(ctx *gin.Context)
	UpdateRestorationOrder(ctx *gin.Context)
	DeleteRestorationOrder(ctx *gin.Context)
	SearchPaginated(ctx *gin.Context)
}

type RestorationOrderController struct {
	restorationOrderUseCase usecaseProduct.IRestorationOrderUseCase
	Logger                  *logger.Logger
}

func NewRestorationOrderController(useCase usecaseProduct.IRestorationOrderUseCase, loggerInstance *logger.Logger) IRestorationOrderController {
	return &RestorationOrderController{restorationOrderUseCase: useCase, Logger: loggerInstance}
}

func (c *RestorationOrderController) NewRestorationOrder(ctx *gin.Context) {
	c.Logger.Info("Creating new restoration order")
	var request NewRestorationOrderRequest
	if err := controllers.BindJSON(ctx, &request); err != nil {
		c.Logger.Error("Error binding JSON for new restoration order", zap.Error(err))
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	model, err := c.restorationOrderUseCase.Create(toRestorationOrderUsecaseMapper(&request))
	if err != nil {
		c.Logger.Error("Error creating restoration order", zap.Error(err), zap.String("productId", request.ProductID))
		_ = ctx.Error(err)
		return
	}
	response := domainToRestorationOrderResponseMapper(model)
	ctx.JSON(http.StatusCreated, response)
}

func (c *RestorationOrderController) GetAllRestorationOrders(ctx *gin.Context) {
	items, err := c.restorationOrderUseCase.GetAll()
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, arrayDomainToRestorationOrderResponseMapper(items))
}

func (c *RestorationOrderController) GetRestorationOrderByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		appError := domainErrors.NewAppError(errors.New("restoration order id is invalid"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	item, err := c.restorationOrderUseCase.GetByID(id)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, domainToRestorationOrderResponseMapper(item))
}

func (c *RestorationOrderController) UpdateRestorationOrder(ctx *gin.Context) {
	id := ctx.Param("id")
	var requestMap map[string]any
	err := controllers.BindJSONMap(ctx, &requestMap)
	if err != nil {
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	updated, err := c.restorationOrderUseCase.Update(id, requestMap)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, domainToRestorationOrderResponseMapper(updated))
}

func (c *RestorationOrderController) DeleteRestorationOrder(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.restorationOrderUseCase.Delete(id)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "resource deleted successfully"})
}

func (c *RestorationOrderController) SearchPaginated(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))

	filters := domain.DataFilters{
		Page:     page,
		PageSize: pageSize,
	}

	likeFilters := make(map[string][]string)
	for field := range repositoryProduct.ColumnsRestorationOrderMapping {
		if values := ctx.QueryArray(field + "_like"); len(values) > 0 {
			likeFilters[field] = values
		}
	}
	filters.LikeFilters = likeFilters

	result, err := c.restorationOrderUseCase.SearchPaginated(filters)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	response := gin.H{
		"data":       arrayDomainToRestorationOrderResponseMapper(result.Data),
		"total":      result.Total,
		"page":       result.Page,
		"pageSize":   result.PageSize,
		"totalPages": result.TotalPages,
		"filters":    filters,
	}
	ctx.JSON(http.StatusOK, response)
}

func domainToRestorationOrderResponseMapper(p *domainProduct.RestorationOrder) *ResponseRestorationOrder {
	return &ResponseRestorationOrder{
		ID:             p.ID,
		ProductID:      p.ProductID,
		Status:         p.Status,
		Priority:       p.Priority,
		EstimatedCost:  p.EstimatedCost,
		ActualCost:     p.ActualCost,
		AssignedTo:     p.AssignedTo,
		StartDate:      p.StartDate,
		CompletionDate: p.CompletionDate,
		CreatedBy:      p.CreatedBy,
		CreatedAt:      p.CreatedAt,
		UpdatedAt:      p.UpdatedAt,
	}
}

func arrayDomainToRestorationOrderResponseMapper(items *[]domainProduct.RestorationOrder) *[]ResponseRestorationOrder {
	res := make([]ResponseRestorationOrder, len(*items))
	for i, p := range *items {
		res[i] = *domainToRestorationOrderResponseMapper(&p)
	}
	return &res
}

func toRestorationOrderUsecaseMapper(req *NewRestorationOrderRequest) *domainProduct.RestorationOrder {
	return &domainProduct.RestorationOrder{
		ProductID:      req.ProductID,
		Status:         req.Status,
		Priority:       req.Priority,
		EstimatedCost:  req.EstimatedCost,
		ActualCost:     req.ActualCost,
		AssignedTo:     req.AssignedTo,
		StartDate:      req.StartDate,
		CompletionDate: req.CompletionDate,
		CreatedBy:      req.CreatedBy,
	}
}
