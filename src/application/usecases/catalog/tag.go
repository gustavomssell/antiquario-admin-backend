package catalog

import (
	"github.com/gbrayhan/microservices-go/src/domain"
	domainCatalog "github.com/gbrayhan/microservices-go/src/domain/catalog"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	repositoryCatalog "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/catalog"
	"go.uber.org/zap"
)

type ITagUseCase interface {
	GetAll() (*[]domainCatalog.Tag, error)
	GetByID(id string) (*domainCatalog.Tag, error)
	Create(newTag *domainCatalog.Tag) (*domainCatalog.Tag, error)
	Update(id string, tagMap map[string]interface{}) (*domainCatalog.Tag, error)
	Delete(id string) error
	SearchPaginated(filters domain.DataFilters) (*domainCatalog.SearchResultTag, error)
}

type TagUseCase struct {
	tagRepository repositoryCatalog.TagRepositoryInterface
	Logger        *logger.Logger
}

func NewTagUseCase(tagRepository repositoryCatalog.TagRepositoryInterface, logger *logger.Logger) ITagUseCase {
	return &TagUseCase{
		tagRepository: tagRepository,
		Logger:        logger,
	}
}

func (s *TagUseCase) GetAll() (*[]domainCatalog.Tag, error) {
	s.Logger.Info("Getting all tags")
	return s.tagRepository.GetAll()
}

func (s *TagUseCase) GetByID(id string) (*domainCatalog.Tag, error) {
	s.Logger.Info("Getting tag by ID", zap.String("id", id))
	return s.tagRepository.GetByID(id)
}

func (s *TagUseCase) Create(newTag *domainCatalog.Tag) (*domainCatalog.Tag, error) {
	s.Logger.Info("Creating new tag", zap.String("name", newTag.Name))
	return s.tagRepository.Create(newTag)
}

func (s *TagUseCase) Update(id string, tagMap map[string]interface{}) (*domainCatalog.Tag, error) {
	s.Logger.Info("Updating tag", zap.String("id", id))
	return s.tagRepository.Update(id, tagMap)
}

func (s *TagUseCase) Delete(id string) error {
	s.Logger.Info("Deleting tag", zap.String("id", id))
	return s.tagRepository.Delete(id)
}

func (s *TagUseCase) SearchPaginated(filters domain.DataFilters) (*domainCatalog.SearchResultTag, error) {
	s.Logger.Info("Searching tags with pagination",
		zap.Int("page", filters.Page),
		zap.Int("pageSize", filters.PageSize))
	return s.tagRepository.SearchPaginated(filters)
}
