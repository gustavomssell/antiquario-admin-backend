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

type NewSaleRequest struct {
	CustomerID      string     `json:"customerId" binding:"required"`
	SaleDate        *time.Time `json:"saleDate"`
	SaleType        string     `json:"saleType"`
	TotalAmount     float64    `json:"totalAmount"`
	PaymentMethod   string     `json:"paymentMethod"`
	PaymentStatus   string     `json:"paymentStatus"`
	DeliveryAddress string     `json:"deliveryAddress"`
	DeliveryDate    *time.Time `json:"deliveryDate"`
	Notes           string     `json:"notes"`
	CreatedBy       int        `json:"createdBy"`
}

type ResponseSale struct {
	ID              string     `json:"id"`
	CustomerID      string     `json:"customerId"`
	SaleDate        *time.Time `json:"saleDate"`
	SaleType        string     `json:"saleType"`
	TotalAmount     float64    `json:"totalAmount"`
	PaymentMethod   string     `json:"paymentMethod"`
	PaymentStatus   string     `json:"paymentStatus"`
	DeliveryAddress string     `json:"deliveryAddress"`
	DeliveryDate    *time.Time `json:"deliveryDate"`
	Notes           string     `json:"notes"`
	CreatedBy       int        `json:"createdBy"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
}

type ISaleController interface {
	NewSale(ctx *gin.Context)
	GetAllSales(ctx *gin.Context)
	GetSaleByID(ctx *gin.Context)
	UpdateSale(ctx *gin.Context)
	DeleteSale(ctx *gin.Context)
	SearchPaginated(ctx *gin.Context)
}

type SaleController struct {
	saleUseCase usecaseCommercial.ISaleUseCase
	Logger      *logger.Logger
}

func NewSaleController(saleUseCase usecaseCommercial.ISaleUseCase, loggerInstance *logger.Logger) ISaleController {
	return &SaleController{saleUseCase: saleUseCase, Logger: loggerInstance}
}

func (c *SaleController) NewSale(ctx *gin.Context) {
	c.Logger.Info("Creating new sale")
	var request NewSaleRequest
	if err := controllers.BindJSON(ctx, &request); err != nil {
		c.Logger.Error("Error binding JSON for new sale", zap.Error(err))
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	model, err := c.saleUseCase.Create(toSaleUsecaseMapper(&request))
	if err != nil {
		c.Logger.Error("Error creating sale", zap.Error(err), zap.String("customerId", request.CustomerID))
		_ = ctx.Error(err)
		return
	}
	response := domainToSaleResponseMapper(model)
	ctx.JSON(http.StatusCreated, response)
}

func (c *SaleController) GetAllSales(ctx *gin.Context) {
	items, err := c.saleUseCase.GetAll()
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, arrayDomainToSaleResponseMapper(items))
}

func (c *SaleController) GetSaleByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		appError := domainErrors.NewAppError(errors.New("sale id is invalid"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	item, err := c.saleUseCase.GetByID(id)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, domainToSaleResponseMapper(item))
}

func (c *SaleController) UpdateSale(ctx *gin.Context) {
	id := ctx.Param("id")
	var requestMap map[string]any
	err := controllers.BindJSONMap(ctx, &requestMap)
	if err != nil {
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	updated, err := c.saleUseCase.Update(id, requestMap)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, domainToSaleResponseMapper(updated))
}

func (c *SaleController) DeleteSale(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.saleUseCase.Delete(id)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "resource deleted successfully"})
}

func (c *SaleController) SearchPaginated(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))

	filters := domain.DataFilters{
		Page:     page,
		PageSize: pageSize,
	}

	likeFilters := make(map[string][]string)
	for field := range repositoryCommercial.ColumnsSaleMapping {
		if values := ctx.QueryArray(field + "_like"); len(values) > 0 {
			likeFilters[field] = values
		}
	}
	filters.LikeFilters = likeFilters

	result, err := c.saleUseCase.SearchPaginated(filters)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	response := gin.H{
		"data":       arrayDomainToSaleResponseMapper(result.Data),
		"total":      result.Total,
		"page":       result.Page,
		"pageSize":   result.PageSize,
		"totalPages": result.TotalPages,
		"filters":    filters,
	}
	ctx.JSON(http.StatusOK, response)
}

func domainToSaleResponseMapper(p *domainCommercial.Sale) *ResponseSale {
	return &ResponseSale{
		ID:              p.ID,
		CustomerID:      p.CustomerID,
		SaleDate:        p.SaleDate,
		SaleType:        p.SaleType,
		TotalAmount:     p.TotalAmount,
		PaymentMethod:   p.PaymentMethod,
		PaymentStatus:   p.PaymentStatus,
		DeliveryAddress: p.DeliveryAddress,
		DeliveryDate:    p.DeliveryDate,
		Notes:           p.Notes,
		CreatedBy:       p.CreatedBy,
		CreatedAt:       p.CreatedAt,
		UpdatedAt:       p.UpdatedAt,
	}
}

func arrayDomainToSaleResponseMapper(items *[]domainCommercial.Sale) *[]ResponseSale {
	res := make([]ResponseSale, len(*items))
	for i, p := range *items {
		res[i] = *domainToSaleResponseMapper(&p)
	}
	return &res
}

func toSaleUsecaseMapper(req *NewSaleRequest) *domainCommercial.Sale {
	return &domainCommercial.Sale{
		CustomerID:      req.CustomerID,
		SaleDate:        req.SaleDate,
		SaleType:        req.SaleType,
		TotalAmount:     req.TotalAmount,
		PaymentMethod:   req.PaymentMethod,
		PaymentStatus:   req.PaymentStatus,
		DeliveryAddress: req.DeliveryAddress,
		DeliveryDate:    req.DeliveryDate,
		Notes:           req.Notes,
		CreatedBy:       req.CreatedBy,
	}
}
