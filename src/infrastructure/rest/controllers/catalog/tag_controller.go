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

type NewTagRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type ResponseTag struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type ITagController interface {
	NewTag(ctx *gin.Context)
	GetAllTags(ctx *gin.Context)
	GetTagByID(ctx *gin.Context)
	UpdateTag(ctx *gin.Context)
	DeleteTag(ctx *gin.Context)
	SearchPaginated(ctx *gin.Context)
}

type TagController struct {
	tagUseCase usecaseCatalog.ITagUseCase
	Logger     *logger.Logger
}

func NewTagController(tagUseCase usecaseCatalog.ITagUseCase, loggerInstance *logger.Logger) ITagController {
	return &TagController{tagUseCase: tagUseCase, Logger: loggerInstance}
}

func (c *TagController) NewTag(ctx *gin.Context) {
	c.Logger.Info("Creating new tag")
	var request NewTagRequest
	if err := controllers.BindJSON(ctx, &request); err != nil {
		c.Logger.Error("Error binding JSON for new tag", zap.Error(err))
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	model, err := c.tagUseCase.Create(toTagUsecaseMapper(&request))
	if err != nil {
		c.Logger.Error("Error creating tag", zap.Error(err), zap.String("name", request.Name))
		_ = ctx.Error(err)
		return
	}
	response := domainToTagResponseMapper(model)
	ctx.JSON(http.StatusCreated, response)
}

func (c *TagController) GetAllTags(ctx *gin.Context) {
	items, err := c.tagUseCase.GetAll()
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, arrayDomainToTagResponseMapper(items))
}

func (c *TagController) GetTagByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		appError := domainErrors.NewAppError(errors.New("tag id is invalid"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	item, err := c.tagUseCase.GetByID(id)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, domainToTagResponseMapper(item))
}

func (c *TagController) UpdateTag(ctx *gin.Context) {
	id := ctx.Param("id")
	var requestMap map[string]any
	err := controllers.BindJSONMap(ctx, &requestMap)
	if err != nil {
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	updated, err := c.tagUseCase.Update(id, requestMap)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, domainToTagResponseMapper(updated))
}

func (c *TagController) DeleteTag(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.tagUseCase.Delete(id)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "resource deleted successfully"})
}

func (c *TagController) SearchPaginated(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))

	filters := domain.DataFilters{
		Page:     page,
		PageSize: pageSize,
	}

	likeFilters := make(map[string][]string)
	for field := range repositoryCatalog.ColumnsTagMapping {
		if values := ctx.QueryArray(field + "_like"); len(values) > 0 {
			likeFilters[field] = values
		}
	}
	filters.LikeFilters = likeFilters

	result, err := c.tagUseCase.SearchPaginated(filters)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	response := gin.H{
		"data":       arrayDomainToTagResponseMapper(result.Data),
		"total":      result.Total,
		"page":       result.Page,
		"pageSize":   result.PageSize,
		"totalPages": result.TotalPages,
		"filters":    filters,
	}
	ctx.JSON(http.StatusOK, response)
}

func domainToTagResponseMapper(p *domainCatalog.Tag) *ResponseTag {
	return &ResponseTag{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

func arrayDomainToTagResponseMapper(items *[]domainCatalog.Tag) *[]ResponseTag {
	res := make([]ResponseTag, len(*items))
	for i, p := range *items {
		res[i] = *domainToTagResponseMapper(&p)
	}
	return &res
}

func toTagUsecaseMapper(req *NewTagRequest) *domainCatalog.Tag {
	return &domainCatalog.Tag{
		Name:        req.Name,
		Description: req.Description,
	}
}
