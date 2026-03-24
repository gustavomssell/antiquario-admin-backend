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

type NewCustomerRequest struct {
	UserID                 int     `json:"userId" binding:"required"`
	CustomerType           string  `json:"customerType"`
	Preferences            string  `json:"preferences"`
	PurchaseHistorySummary string  `json:"purchaseHistorySummary"`
	CreditLimit            float64 `json:"creditLimit"`
	Notes                  string  `json:"notes"`
}

type ResponseCustomer struct {
	ID                     string    `json:"id"`
	UserID                 int       `json:"userId"`
	CustomerType           string    `json:"customerType"`
	Preferences            string    `json:"preferences"`
	PurchaseHistorySummary string    `json:"purchaseHistorySummary"`
	CreditLimit            float64   `json:"creditLimit"`
	Notes                  string    `json:"notes"`
	CreatedAt              time.Time `json:"createdAt"`
	UpdatedAt              time.Time `json:"updatedAt"`
}

type ICustomerController interface {
	NewCustomer(ctx *gin.Context)
	GetAllCustomers(ctx *gin.Context)
	GetCustomerByID(ctx *gin.Context)
	UpdateCustomer(ctx *gin.Context)
	DeleteCustomer(ctx *gin.Context)
	SearchPaginated(ctx *gin.Context)
}

type CustomerController struct {
	customerUseCase usecaseCommercial.ICustomerUseCase
	Logger          *logger.Logger
}

func NewCustomerController(customerUseCase usecaseCommercial.ICustomerUseCase, loggerInstance *logger.Logger) ICustomerController {
	return &CustomerController{customerUseCase: customerUseCase, Logger: loggerInstance}
}

func (c *CustomerController) NewCustomer(ctx *gin.Context) {
	c.Logger.Info("Creating new customer")
	var request NewCustomerRequest
	if err := controllers.BindJSON(ctx, &request); err != nil {
		c.Logger.Error("Error binding JSON for new customer", zap.Error(err))
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	model, err := c.customerUseCase.Create(toCustomerUsecaseMapper(&request))
	if err != nil {
		c.Logger.Error("Error creating customer", zap.Error(err), zap.Int("userId", request.UserID))
		_ = ctx.Error(err)
		return
	}
	response := domainToCustomerResponseMapper(model)
	ctx.JSON(http.StatusCreated, response)
}

func (c *CustomerController) GetAllCustomers(ctx *gin.Context) {
	items, err := c.customerUseCase.GetAll()
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, arrayDomainToCustomerResponseMapper(items))
}

func (c *CustomerController) GetCustomerByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		appError := domainErrors.NewAppError(errors.New("customer id is invalid"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	item, err := c.customerUseCase.GetByID(id)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, domainToCustomerResponseMapper(item))
}

func (c *CustomerController) UpdateCustomer(ctx *gin.Context) {
	id := ctx.Param("id")
	var requestMap map[string]any
	err := controllers.BindJSONMap(ctx, &requestMap)
	if err != nil {
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	updated, err := c.customerUseCase.Update(id, requestMap)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, domainToCustomerResponseMapper(updated))
}

func (c *CustomerController) DeleteCustomer(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.customerUseCase.Delete(id)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "resource deleted successfully"})
}

func (c *CustomerController) SearchPaginated(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))

	filters := domain.DataFilters{
		Page:     page,
		PageSize: pageSize,
	}

	likeFilters := make(map[string][]string)
	for field := range repositoryCommercial.ColumnsCustomerMapping {
		if values := ctx.QueryArray(field + "_like"); len(values) > 0 {
			likeFilters[field] = values
		}
	}
	filters.LikeFilters = likeFilters

	result, err := c.customerUseCase.SearchPaginated(filters)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	response := gin.H{
		"data":       arrayDomainToCustomerResponseMapper(result.Data),
		"total":      result.Total,
		"page":       result.Page,
		"pageSize":   result.PageSize,
		"totalPages": result.TotalPages,
		"filters":    filters,
	}
	ctx.JSON(http.StatusOK, response)
}

func domainToCustomerResponseMapper(p *domainCommercial.Customer) *ResponseCustomer {
	return &ResponseCustomer{
		ID:                     p.ID,
		UserID:                 p.UserID,
		CustomerType:           p.CustomerType,
		Preferences:            p.Preferences,
		PurchaseHistorySummary: p.PurchaseHistorySummary,
		CreditLimit:            p.CreditLimit,
		Notes:                  p.Notes,
		CreatedAt:              p.CreatedAt,
		UpdatedAt:              p.UpdatedAt,
	}
}

func arrayDomainToCustomerResponseMapper(items *[]domainCommercial.Customer) *[]ResponseCustomer {
	res := make([]ResponseCustomer, len(*items))
	for i, p := range *items {
		res[i] = *domainToCustomerResponseMapper(&p)
	}
	return &res
}

func toCustomerUsecaseMapper(req *NewCustomerRequest) *domainCommercial.Customer {
	return &domainCommercial.Customer{
		UserID:                 req.UserID,
		CustomerType:           req.CustomerType,
		Preferences:            req.Preferences,
		PurchaseHistorySummary: req.PurchaseHistorySummary,
		CreditLimit:            req.CreditLimit,
		Notes:                  req.Notes,
	}
}
