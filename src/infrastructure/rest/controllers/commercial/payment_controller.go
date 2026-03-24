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

type NewPaymentRequest struct {
	SaleID        string     `json:"saleId" binding:"required"`
	Amount        float64    `json:"amount"`
	PaymentDate   *time.Time `json:"paymentDate"`
	PaymentMethod string     `json:"paymentMethod"`
	Status        string     `json:"status"`
	Reference     string     `json:"reference"`
}

type ResponsePayment struct {
	ID            string     `json:"id"`
	SaleID        string     `json:"saleId"`
	Amount        float64    `json:"amount"`
	PaymentDate   *time.Time `json:"paymentDate"`
	PaymentMethod string     `json:"paymentMethod"`
	Status        string     `json:"status"`
	Reference     string     `json:"reference"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
}

type IPaymentController interface {
	NewPayment(ctx *gin.Context)
	GetAllPayments(ctx *gin.Context)
	GetPaymentByID(ctx *gin.Context)
	UpdatePayment(ctx *gin.Context)
	DeletePayment(ctx *gin.Context)
	SearchPaginated(ctx *gin.Context)
}

type PaymentController struct {
	paymentUseCase usecaseCommercial.IPaymentUseCase
	Logger         *logger.Logger
}

func NewPaymentController(paymentUseCase usecaseCommercial.IPaymentUseCase, loggerInstance *logger.Logger) IPaymentController {
	return &PaymentController{paymentUseCase: paymentUseCase, Logger: loggerInstance}
}

func (c *PaymentController) NewPayment(ctx *gin.Context) {
	c.Logger.Info("Creating new payment")
	var request NewPaymentRequest
	if err := controllers.BindJSON(ctx, &request); err != nil {
		c.Logger.Error("Error binding JSON for new payment", zap.Error(err))
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	model, err := c.paymentUseCase.Create(toPaymentUsecaseMapper(&request))
	if err != nil {
		c.Logger.Error("Error creating payment", zap.Error(err), zap.String("saleId", request.SaleID))
		_ = ctx.Error(err)
		return
	}
	response := domainToPaymentResponseMapper(model)
	ctx.JSON(http.StatusCreated, response)
}

func (c *PaymentController) GetAllPayments(ctx *gin.Context) {
	items, err := c.paymentUseCase.GetAll()
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, arrayDomainToPaymentResponseMapper(items))
}

func (c *PaymentController) GetPaymentByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		appError := domainErrors.NewAppError(errors.New("payment id is invalid"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	item, err := c.paymentUseCase.GetByID(id)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, domainToPaymentResponseMapper(item))
}

func (c *PaymentController) UpdatePayment(ctx *gin.Context) {
	id := ctx.Param("id")
	var requestMap map[string]any
	err := controllers.BindJSONMap(ctx, &requestMap)
	if err != nil {
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	updated, err := c.paymentUseCase.Update(id, requestMap)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, domainToPaymentResponseMapper(updated))
}

func (c *PaymentController) DeletePayment(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.paymentUseCase.Delete(id)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "resource deleted successfully"})
}

func (c *PaymentController) SearchPaginated(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))

	filters := domain.DataFilters{
		Page:     page,
		PageSize: pageSize,
	}

	likeFilters := make(map[string][]string)
	for field := range repositoryCommercial.ColumnsPaymentMapping {
		if values := ctx.QueryArray(field + "_like"); len(values) > 0 {
			likeFilters[field] = values
		}
	}
	filters.LikeFilters = likeFilters

	result, err := c.paymentUseCase.SearchPaginated(filters)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	response := gin.H{
		"data":       arrayDomainToPaymentResponseMapper(result.Data),
		"total":      result.Total,
		"page":       result.Page,
		"pageSize":   result.PageSize,
		"totalPages": result.TotalPages,
		"filters":    filters,
	}
	ctx.JSON(http.StatusOK, response)
}

func domainToPaymentResponseMapper(p *domainCommercial.Payment) *ResponsePayment {
	return &ResponsePayment{
		ID:            p.ID,
		SaleID:        p.SaleID,
		Amount:        p.Amount,
		PaymentDate:   p.PaymentDate,
		PaymentMethod: p.PaymentMethod,
		Status:        p.Status,
		Reference:     p.Reference,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}
}

func arrayDomainToPaymentResponseMapper(items *[]domainCommercial.Payment) *[]ResponsePayment {
	res := make([]ResponsePayment, len(*items))
	for i, p := range *items {
		res[i] = *domainToPaymentResponseMapper(&p)
	}
	return &res
}

func toPaymentUsecaseMapper(req *NewPaymentRequest) *domainCommercial.Payment {
	return &domainCommercial.Payment{
		SaleID:        req.SaleID,
		Amount:        req.Amount,
		PaymentDate:   req.PaymentDate,
		PaymentMethod: req.PaymentMethod,
		Status:        req.Status,
		Reference:     req.Reference,
	}
}
