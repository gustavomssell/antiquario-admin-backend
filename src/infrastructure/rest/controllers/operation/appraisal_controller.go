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

type NewAppraisalRequest struct {
	ProductID           string     `json:"productId" binding:"required"`
	AppraiserName       string     `json:"appraiserName"`
	AppraisalDate       *time.Time `json:"appraisalDate"`
	EstimatedValue      float64    `json:"estimatedValue"`
	ConditionAssessment string     `json:"conditionAssessment"`
	AuthenticityRating  float64    `json:"authenticityRating"`
	Notes               string     `json:"notes"`
	DocumentURL         string     `json:"documentUrl"`
}

type ResponseAppraisal struct {
	ID                  string     `json:"id"`
	ProductID           string     `json:"productId"`
	AppraiserName       string     `json:"appraiserName"`
	AppraisalDate       *time.Time `json:"appraisalDate"`
	EstimatedValue      float64    `json:"estimatedValue"`
	ConditionAssessment string     `json:"conditionAssessment"`
	AuthenticityRating  float64    `json:"authenticityRating"`
	Notes               string     `json:"notes"`
	DocumentURL         string     `json:"documentUrl"`
	CreatedAt           time.Time  `json:"createdAt"`
	UpdatedAt           time.Time  `json:"updatedAt"`
}

type IAppraisalController interface {
	NewAppraisal(ctx *gin.Context)
	GetAllAppraisals(ctx *gin.Context)
	GetAppraisalByID(ctx *gin.Context)
	UpdateAppraisal(ctx *gin.Context)
	DeleteAppraisal(ctx *gin.Context)
	SearchPaginated(ctx *gin.Context)
}

type AppraisalController struct {
	appraisalUseCase usecaseOperation.IAppraisalUseCase
	Logger           *logger.Logger
}

func NewAppraisalController(appraisalUseCase usecaseOperation.IAppraisalUseCase, loggerInstance *logger.Logger) IAppraisalController {
	return &AppraisalController{appraisalUseCase: appraisalUseCase, Logger: loggerInstance}
}

func (c *AppraisalController) NewAppraisal(ctx *gin.Context) {
	c.Logger.Info("Creating new appraisal")
	var request NewAppraisalRequest
	if err := controllers.BindJSON(ctx, &request); err != nil {
		c.Logger.Error("Error binding JSON for new appraisal", zap.Error(err))
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	model, err := c.appraisalUseCase.Create(toAppraisalUsecaseMapper(&request))
	if err != nil {
		c.Logger.Error("Error creating appraisal", zap.Error(err), zap.String("productId", request.ProductID))
		_ = ctx.Error(err)
		return
	}
	response := domainToAppraisalResponseMapper(model)
	ctx.JSON(http.StatusCreated, response)
}

func (c *AppraisalController) GetAllAppraisals(ctx *gin.Context) {
	items, err := c.appraisalUseCase.GetAll()
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, arrayDomainToAppraisalResponseMapper(items))
}

func (c *AppraisalController) GetAppraisalByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		appError := domainErrors.NewAppError(errors.New("appraisal id is invalid"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	item, err := c.appraisalUseCase.GetByID(id)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, domainToAppraisalResponseMapper(item))
}

func (c *AppraisalController) UpdateAppraisal(ctx *gin.Context) {
	id := ctx.Param("id")
	var requestMap map[string]any
	err := controllers.BindJSONMap(ctx, &requestMap)
	if err != nil {
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	updated, err := c.appraisalUseCase.Update(id, requestMap)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, domainToAppraisalResponseMapper(updated))
}

func (c *AppraisalController) DeleteAppraisal(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.appraisalUseCase.Delete(id)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "resource deleted successfully"})
}

func (c *AppraisalController) SearchPaginated(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))

	filters := domain.DataFilters{
		Page:     page,
		PageSize: pageSize,
	}

	likeFilters := make(map[string][]string)
	for field := range repositoryOperation.ColumnsAppraisalMapping {
		if values := ctx.QueryArray(field + "_like"); len(values) > 0 {
			likeFilters[field] = values
		}
	}
	filters.LikeFilters = likeFilters

	result, err := c.appraisalUseCase.SearchPaginated(filters)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	response := gin.H{
		"data":       arrayDomainToAppraisalResponseMapper(result.Data),
		"total":      result.Total,
		"page":       result.Page,
		"pageSize":   result.PageSize,
		"totalPages": result.TotalPages,
		"filters":    filters,
	}
	ctx.JSON(http.StatusOK, response)
}

func domainToAppraisalResponseMapper(p *domainOperation.Appraisal) *ResponseAppraisal {
	return &ResponseAppraisal{
		ID:                  p.ID,
		ProductID:           p.ProductID,
		AppraiserName:       p.AppraiserName,
		AppraisalDate:       p.AppraisalDate,
		EstimatedValue:      p.EstimatedValue,
		ConditionAssessment: p.ConditionAssessment,
		AuthenticityRating:  p.AuthenticityRating,
		Notes:               p.Notes,
		DocumentURL:         p.DocumentURL,
		CreatedAt:           p.CreatedAt,
		UpdatedAt:           p.UpdatedAt,
	}
}

func arrayDomainToAppraisalResponseMapper(items *[]domainOperation.Appraisal) *[]ResponseAppraisal {
	res := make([]ResponseAppraisal, len(*items))
	for i, p := range *items {
		res[i] = *domainToAppraisalResponseMapper(&p)
	}
	return &res
}

func toAppraisalUsecaseMapper(req *NewAppraisalRequest) *domainOperation.Appraisal {
	return &domainOperation.Appraisal{
		ProductID:           req.ProductID,
		AppraiserName:       req.AppraiserName,
		AppraisalDate:       req.AppraisalDate,
		EstimatedValue:      req.EstimatedValue,
		ConditionAssessment: req.ConditionAssessment,
		AuthenticityRating:  req.AuthenticityRating,
		Notes:               req.Notes,
		DocumentURL:         req.DocumentURL,
	}
}
