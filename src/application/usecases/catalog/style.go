package catalog

import (
	"github.com/gbrayhan/microservices-go/src/domain"
	domainCatalog "github.com/gbrayhan/microservices-go/src/domain/catalog"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	repositoryCatalog "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/catalog"
	"go.uber.org/zap"
)

type IStyleUseCase interface {
	GetAll() (*[]domainCatalog.Style, error)
	GetByID(id string) (*domainCatalog.Style, error)
	Create(newStyle *domainCatalog.Style) (*domainCatalog.Style, error)
	Update(id string, styleMap map[string]interface{}) (*domainCatalog.Style, error)
	Delete(id string) error
	SearchPaginated(filters domain.DataFilters) (*domainCatalog.SearchResultStyle, error)
}

type StyleUseCase struct {
	styleRepository repositoryCatalog.StyleRepositoryInterface
	Logger          *logger.Logger
}

func NewStyleUseCase(styleRepository repositoryCatalog.StyleRepositoryInterface, logger *logger.Logger) IStyleUseCase {
	return &StyleUseCase{
		styleRepository: styleRepository,
		Logger:          logger,
	}
}

func (s *StyleUseCase) GetAll() (*[]domainCatalog.Style, error) {
	s.Logger.Info("Getting all styles")
	return s.styleRepository.GetAll()
}

func (s *StyleUseCase) GetByID(id string) (*domainCatalog.Style, error) {
	s.Logger.Info("Getting style by ID", zap.String("id", id))
	return s.styleRepository.GetByID(id)
}

func (s *StyleUseCase) Create(newStyle *domainCatalog.Style) (*domainCatalog.Style, error) {
	s.Logger.Info("Creating new style", zap.String("name", newStyle.Name))
	return s.styleRepository.Create(newStyle)
}

func (s *StyleUseCase) Update(id string, styleMap map[string]interface{}) (*domainCatalog.Style, error) {
	s.Logger.Info("Updating style", zap.String("id", id))
	return s.styleRepository.Update(id, styleMap)
}

func (s *StyleUseCase) Delete(id string) error {
	s.Logger.Info("Deleting style", zap.String("id", id))
	return s.styleRepository.Delete(id)
}

func (s *StyleUseCase) SearchPaginated(filters domain.DataFilters) (*domainCatalog.SearchResultStyle, error) {
	s.Logger.Info("Searching styles with pagination",
		zap.Int("page", filters.Page),
		zap.Int("pageSize", filters.PageSize))
	return s.styleRepository.SearchPaginated(filters)
}
