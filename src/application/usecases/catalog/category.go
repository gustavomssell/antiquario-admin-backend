package catalog

import (
	"github.com/gbrayhan/microservices-go/src/domain"
	domainCatalog "github.com/gbrayhan/microservices-go/src/domain/catalog"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	repositoryCatalog "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/catalog"
	"go.uber.org/zap"
)

type ICategoryUseCase interface {
	GetAll() (*[]domainCatalog.Category, error)
	GetByID(id string) (*domainCatalog.Category, error)
	Create(newCategory *domainCatalog.Category) (*domainCatalog.Category, error)
	Update(id string, categoryMap map[string]interface{}) (*domainCatalog.Category, error)
	Delete(id string) error
	SearchPaginated(filters domain.DataFilters) (*domainCatalog.SearchResultCategory, error)
}

type CategoryUseCase struct {
	categoryRepository repositoryCatalog.CategoryRepositoryInterface
	Logger             *logger.Logger
}

func NewCategoryUseCase(categoryRepository repositoryCatalog.CategoryRepositoryInterface, logger *logger.Logger) ICategoryUseCase {
	return &CategoryUseCase{
		categoryRepository: categoryRepository,
		Logger:             logger,
	}
}

func (s *CategoryUseCase) GetAll() (*[]domainCatalog.Category, error) {
	s.Logger.Info("Getting all categories")
	return s.categoryRepository.GetAll()
}

func (s *CategoryUseCase) GetByID(id string) (*domainCatalog.Category, error) {
	s.Logger.Info("Getting category by ID", zap.String("id", id))
	return s.categoryRepository.GetByID(id)
}

func (s *CategoryUseCase) Create(newCategory *domainCatalog.Category) (*domainCatalog.Category, error) {
	s.Logger.Info("Creating new category", zap.String("name", newCategory.Name))
	return s.categoryRepository.Create(newCategory)
}

func (s *CategoryUseCase) Update(id string, categoryMap map[string]interface{}) (*domainCatalog.Category, error) {
	s.Logger.Info("Updating category", zap.String("id", id))
	return s.categoryRepository.Update(id, categoryMap)
}

func (s *CategoryUseCase) Delete(id string) error {
	s.Logger.Info("Deleting category", zap.String("id", id))
	return s.categoryRepository.Delete(id)
}

func (s *CategoryUseCase) SearchPaginated(filters domain.DataFilters) (*domainCatalog.SearchResultCategory, error) {
	s.Logger.Info("Searching categories with pagination",
		zap.Int("page", filters.Page),
		zap.Int("pageSize", filters.PageSize))
	return s.categoryRepository.SearchPaginated(filters)
}
