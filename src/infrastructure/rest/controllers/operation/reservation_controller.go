package operation

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gbrayhan/microservices-go/src/domain"
	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	domainOperation "github.com/gbrayhan/microservices-go/src/domain/operation"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	repositoryOperation "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/operation"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers"
	usecaseOperation "github.com/gbrayhan/microservices-go/src/application/usecases/operation"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type NewReservationRequest struct {
	ProductID             string     `json:"productId" binding:"required"`
	CustomerID            string     `json:"customerId" binding:"required"`
	ReservationDate       *time.Time `json:"reservationDate"`
	ExpiryDate            *time.Time `json:"expiryDate"`
	CustomReservationDays *int       `json:"customReservationDays"`
	RequiresDeposit       bool       `json:"requiresDeposit"`
	DepositPercentage     float64    `json:"depositPercentage"`
	DepositAmount         float64    `json:"depositAmount"`
	DepositPaid           bool       `json:"depositPaid"`
	DepositDate           *time.Time `json:"depositDate"`
	Status                string     `json:"status"`
	CancellationReason    string     `json:"cancellationReason"`
	Notes                 string     `json:"notes"`
	CreatedBy             int        `json:"createdBy"`
}

type ResponseReservation struct {
	ID                    string     `json:"id"`
	ProductID             string     `json:"productId"`
	CustomerID            string     `json:"customerId"`
	ReservationDate       *time.Time `json:"reservationDate"`
	ExpiryDate            *time.Time `json:"expiryDate"`
	CustomReservationDays *int       `json:"customReservationDays"`
	RequiresDeposit       bool       `json:"requiresDeposit"`
	DepositPercentage     float64    `json:"depositPercentage"`
	DepositAmount         float64    `json:"depositAmount"`
	DepositPaid           bool       `json:"depositPaid"`
	DepositDate           *time.Time `json:"depositDate"`
	Status                string     `json:"status"`
	CancellationReason    string     `json:"cancellationReason"`
	Notes                 string     `json:"notes"`
	CreatedBy             int        `json:"createdBy"`
	CreatedAt             time.Time  `json:"createdAt"`
	UpdatedAt             time.Time  `json:"updatedAt"`
}

type IReservationController interface {
	NewReservation(ctx *gin.Context)
	GetAllReservations(ctx *gin.Context)
	GetReservationByID(ctx *gin.Context)
	UpdateReservation(ctx *gin.Context)
	DeleteReservation(ctx *gin.Context)
	SearchPaginated(ctx *gin.Context)
}

type ReservationController struct {
	reservationUseCase usecaseOperation.IReservationUseCase
	Logger             *logger.Logger
}

func NewReservationController(reservationUseCase usecaseOperation.IReservationUseCase, loggerInstance *logger.Logger) IReservationController {
	return &ReservationController{reservationUseCase: reservationUseCase, Logger: loggerInstance}
}

func (c *ReservationController) NewReservation(ctx *gin.Context) {
	c.Logger.Info("Creating new reservation")
	var request NewReservationRequest
	if err := controllers.BindJSON(ctx, &request); err != nil {
		c.Logger.Error("Error binding JSON for new reservation", zap.Error(err))
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	model, err := c.reservationUseCase.Create(toReservationUsecaseMapper(&request))
	if err != nil {
		c.Logger.Error("Error creating reservation", zap.Error(err), zap.String("productId", request.ProductID))
		_ = ctx.Error(err)
		return
	}
	response := domainToReservationResponseMapper(model)
	ctx.JSON(http.StatusCreated, response)
}

func (c *ReservationController) GetAllReservations(ctx *gin.Context) {
	items, err := c.reservationUseCase.GetAll()
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, arrayDomainToReservationResponseMapper(items))
}

func (c *ReservationController) GetReservationByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		appError := domainErrors.NewAppError(errors.New("reservation id is invalid"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	item, err := c.reservationUseCase.GetByID(id)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, domainToReservationResponseMapper(item))
}

func (c *ReservationController) UpdateReservation(ctx *gin.Context) {
	id := ctx.Param("id")
	var requestMap map[string]any
	err := controllers.BindJSONMap(ctx, &requestMap)
	if err != nil {
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	updated, err := c.reservationUseCase.Update(id, requestMap)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, domainToReservationResponseMapper(updated))
}

func (c *ReservationController) DeleteReservation(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.reservationUseCase.Delete(id)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "resource deleted successfully"})
}

func (c *ReservationController) SearchPaginated(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))

	filters := domain.DataFilters{
		Page:     page,
		PageSize: pageSize,
	}

	likeFilters := make(map[string][]string)
	for field := range repositoryOperation.ColumnsReservationMapping {
		if values := ctx.QueryArray(field + "_like"); len(values) > 0 {
			likeFilters[field] = values
		}
	}
	filters.LikeFilters = likeFilters

	result, err := c.reservationUseCase.SearchPaginated(filters)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	response := gin.H{
		"data":       arrayDomainToReservationResponseMapper(result.Data),
		"total":      result.Total,
		"page":       result.Page,
		"pageSize":   result.PageSize,
		"totalPages": result.TotalPages,
		"filters":    filters,
	}
	ctx.JSON(http.StatusOK, response)
}

func domainToReservationResponseMapper(p *domainOperation.Reservation) *ResponseReservation {
	return &ResponseReservation{
		ID:                    p.ID,
		ProductID:             p.ProductID,
		CustomerID:            p.CustomerID,
		ReservationDate:       p.ReservationDate,
		ExpiryDate:            p.ExpiryDate,
		CustomReservationDays: p.CustomReservationDays,
		RequiresDeposit:       p.RequiresDeposit,
		DepositPercentage:     p.DepositPercentage,
		DepositAmount:         p.DepositAmount,
		DepositPaid:           p.DepositPaid,
		DepositDate:           p.DepositDate,
		Status:                p.Status,
		CancellationReason:    p.CancellationReason,
		Notes:                 p.Notes,
		CreatedBy:             p.CreatedBy,
		CreatedAt:             p.CreatedAt,
		UpdatedAt:             p.UpdatedAt,
	}
}

func arrayDomainToReservationResponseMapper(items *[]domainOperation.Reservation) *[]ResponseReservation {
	res := make([]ResponseReservation, len(*items))
	for i, p := range *items {
		res[i] = *domainToReservationResponseMapper(&p)
	}
	return &res
}

func toReservationUsecaseMapper(req *NewReservationRequest) *domainOperation.Reservation {
	return &domainOperation.Reservation{
		ProductID:             req.ProductID,
		CustomerID:            req.CustomerID,
		ReservationDate:       req.ReservationDate,
		ExpiryDate:            req.ExpiryDate,
		CustomReservationDays: req.CustomReservationDays,
		RequiresDeposit:       req.RequiresDeposit,
		DepositPercentage:     req.DepositPercentage,
		DepositAmount:         req.DepositAmount,
		DepositPaid:           req.DepositPaid,
		DepositDate:           req.DepositDate,
		Status:                req.Status,
		CancellationReason:    req.CancellationReason,
		Notes:                 req.Notes,
		CreatedBy:             req.CreatedBy,
	}
}
