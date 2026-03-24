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

type NewAuctionRequest struct {
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	StartDate   *time.Time `json:"startDate"`
	EndDate     *time.Time `json:"endDate"`
	Status      string     `json:"status"`
	CreatedBy   int        `json:"createdBy"`
}

type ResponseAuction struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	StartDate   *time.Time `json:"startDate"`
	EndDate     *time.Time `json:"endDate"`
	Status      string     `json:"status"`
	CreatedBy   int        `json:"createdBy"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

type IAuctionController interface {
	NewAuction(ctx *gin.Context)
	GetAllAuctions(ctx *gin.Context)
	GetAuctionByID(ctx *gin.Context)
	UpdateAuction(ctx *gin.Context)
	DeleteAuction(ctx *gin.Context)
	SearchPaginated(ctx *gin.Context)
}

type AuctionController struct {
	auctionUseCase usecaseOperation.IAuctionUseCase
	Logger         *logger.Logger
}

func NewAuctionController(auctionUseCase usecaseOperation.IAuctionUseCase, loggerInstance *logger.Logger) IAuctionController {
	return &AuctionController{auctionUseCase: auctionUseCase, Logger: loggerInstance}
}

func (c *AuctionController) NewAuction(ctx *gin.Context) {
	c.Logger.Info("Creating new auction")
	var request NewAuctionRequest
	if err := controllers.BindJSON(ctx, &request); err != nil {
		c.Logger.Error("Error binding JSON for new auction", zap.Error(err))
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	model, err := c.auctionUseCase.Create(toAuctionUsecaseMapper(&request))
	if err != nil {
		c.Logger.Error("Error creating auction", zap.Error(err), zap.String("title", request.Title))
		_ = ctx.Error(err)
		return
	}
	response := domainToAuctionResponseMapper(model)
	ctx.JSON(http.StatusCreated, response)
}

func (c *AuctionController) GetAllAuctions(ctx *gin.Context) {
	items, err := c.auctionUseCase.GetAll()
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, arrayDomainToAuctionResponseMapper(items))
}

func (c *AuctionController) GetAuctionByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		appError := domainErrors.NewAppError(errors.New("auction id is invalid"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	item, err := c.auctionUseCase.GetByID(id)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, domainToAuctionResponseMapper(item))
}

func (c *AuctionController) UpdateAuction(ctx *gin.Context) {
	id := ctx.Param("id")
	var requestMap map[string]any
	err := controllers.BindJSONMap(ctx, &requestMap)
	if err != nil {
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	updated, err := c.auctionUseCase.Update(id, requestMap)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, domainToAuctionResponseMapper(updated))
}

func (c *AuctionController) DeleteAuction(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.auctionUseCase.Delete(id)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "resource deleted successfully"})
}

func (c *AuctionController) SearchPaginated(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))

	filters := domain.DataFilters{
		Page:     page,
		PageSize: pageSize,
	}

	likeFilters := make(map[string][]string)
	for field := range repositoryOperation.ColumnsAuctionMapping {
		if values := ctx.QueryArray(field + "_like"); len(values) > 0 {
			likeFilters[field] = values
		}
	}
	filters.LikeFilters = likeFilters

	result, err := c.auctionUseCase.SearchPaginated(filters)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	response := gin.H{
		"data":       arrayDomainToAuctionResponseMapper(result.Data),
		"total":      result.Total,
		"page":       result.Page,
		"pageSize":   result.PageSize,
		"totalPages": result.TotalPages,
		"filters":    filters,
	}
	ctx.JSON(http.StatusOK, response)
}

func domainToAuctionResponseMapper(p *domainOperation.Auction) *ResponseAuction {
	return &ResponseAuction{
		ID:          p.ID,
		Title:       p.Title,
		Description: p.Description,
		StartDate:   p.StartDate,
		EndDate:     p.EndDate,
		Status:      p.Status,
		CreatedBy:   p.CreatedBy,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

func arrayDomainToAuctionResponseMapper(items *[]domainOperation.Auction) *[]ResponseAuction {
	res := make([]ResponseAuction, len(*items))
	for i, p := range *items {
		res[i] = *domainToAuctionResponseMapper(&p)
	}
	return &res
}

func toAuctionUsecaseMapper(req *NewAuctionRequest) *domainOperation.Auction {
	return &domainOperation.Auction{
		Title:       req.Title,
		Description: req.Description,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Status:      req.Status,
		CreatedBy:   req.CreatedBy,
	}
}
