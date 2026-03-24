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

type NewSupplierRequest struct {
	Name           string `json:"name" binding:"required"`
	Type           string `json:"type"`
	ContactInfo    string `json:"contactInfo"`
	Address        string `json:"address"`
	DocumentNumber string `json:"documentNumber"`
	Notes          string `json:"notes"`
}

type ResponseSupplier struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Type           string    `json:"type"`
	ContactInfo    string    `json:"contactInfo"`
	Address        string    `json:"address"`
	DocumentNumber string    `json:"documentNumber"`
	Notes          string    `json:"notes"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type ISupplierController interface {
	NewSupplier(ctx *gin.Context)
	GetAllSuppliers(ctx *gin.Context)
	GetSupplierByID(ctx *gin.Context)
	UpdateSupplier(ctx *gin.Context)
	DeleteSupplier(ctx *gin.Context)
	SearchPaginated(ctx *gin.Context)
}

type SupplierController struct {
	supplierUseCase usecaseCommercial.ISupplierUseCase
	Logger          *logger.Logger
}

func NewSupplierController(supplierUseCase usecaseCommercial.ISupplierUseCase, loggerInstance *logger.Logger) ISupplierController {
	return &SupplierController{supplierUseCase: supplierUseCase, Logger: loggerInstance}
}

func (c *SupplierController) NewSupplier(ctx *gin.Context) {
	c.Logger.Info("Creating new supplier")
	var request NewSupplierRequest
	if err := controllers.BindJSON(ctx, &request); err != nil {
		c.Logger.Error("Error binding JSON for new supplier", zap.Error(err))
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	model, err := c.supplierUseCase.Create(toSupplierUsecaseMapper(&request))
	if err != nil {
		c.Logger.Error("Error creating supplier", zap.Error(err), zap.String("name", request.Name))
		_ = ctx.Error(err)
		return
	}
	response := domainToSupplierResponseMapper(model)
	ctx.JSON(http.StatusCreated, response)
}

func (c *SupplierController) GetAllSuppliers(ctx *gin.Context) {
	items, err := c.supplierUseCase.GetAll()
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, arrayDomainToSupplierResponseMapper(items))
}

func (c *SupplierController) GetSupplierByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		appError := domainErrors.NewAppError(errors.New("supplier id is invalid"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	item, err := c.supplierUseCase.GetByID(id)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, domainToSupplierResponseMapper(item))
}

func (c *SupplierController) UpdateSupplier(ctx *gin.Context) {
	id := ctx.Param("id")
	var requestMap map[string]any
	err := controllers.BindJSONMap(ctx, &requestMap)
	if err != nil {
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	updated, err := c.supplierUseCase.Update(id, requestMap)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, domainToSupplierResponseMapper(updated))
}

func (c *SupplierController) DeleteSupplier(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.supplierUseCase.Delete(id)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "resource deleted successfully"})
}

func (c *SupplierController) SearchPaginated(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))

	filters := domain.DataFilters{
		Page:     page,
		PageSize: pageSize,
	}

	likeFilters := make(map[string][]string)
	for field := range repositoryCommercial.ColumnsSupplierMapping {
		if values := ctx.QueryArray(field + "_like"); len(values) > 0 {
			likeFilters[field] = values
		}
	}
	filters.LikeFilters = likeFilters

	result, err := c.supplierUseCase.SearchPaginated(filters)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	response := gin.H{
		"data":       arrayDomainToSupplierResponseMapper(result.Data),
		"total":      result.Total,
		"page":       result.Page,
		"pageSize":   result.PageSize,
		"totalPages": result.TotalPages,
		"filters":    filters,
	}
	ctx.JSON(http.StatusOK, response)
}

func domainToSupplierResponseMapper(p *domainCommercial.Supplier) *ResponseSupplier {
	return &ResponseSupplier{
		ID:             p.ID,
		Name:           p.Name,
		Type:           p.Type,
		ContactInfo:    p.ContactInfo,
		Address:        p.Address,
		DocumentNumber: p.DocumentNumber,
		Notes:          p.Notes,
		CreatedAt:      p.CreatedAt,
		UpdatedAt:      p.UpdatedAt,
	}
}

func arrayDomainToSupplierResponseMapper(items *[]domainCommercial.Supplier) *[]ResponseSupplier {
	res := make([]ResponseSupplier, len(*items))
	for i, p := range *items {
		res[i] = *domainToSupplierResponseMapper(&p)
	}
	return &res
}

func toSupplierUsecaseMapper(req *NewSupplierRequest) *domainCommercial.Supplier {
	return &domainCommercial.Supplier{
		Name:           req.Name,
		Type:           req.Type,
		ContactInfo:    req.ContactInfo,
		Address:        req.Address,
		DocumentNumber: req.DocumentNumber,
		Notes:          req.Notes,
	}
}
