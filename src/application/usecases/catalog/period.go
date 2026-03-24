package catalog

import (
	"github.com/gbrayhan/microservices-go/src/domain"
	domainCatalog "github.com/gbrayhan/microservices-go/src/domain/catalog"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	repositoryCatalog "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/catalog"
	"go.uber.org/zap"
)

type IPeriodUseCase interface {
	GetAll() (*[]domainCatalog.Period, error)
	GetByID(id string) (*domainCatalog.Period, error)
	Create(newPeriod *domainCatalog.Period) (*domainCatalog.Period, error)
	Update(id string, periodMap map[string]interface{}) (*domainCatalog.Period, error)
	Delete(id string) error
	SearchPaginated(filters domain.DataFilters) (*domainCatalog.SearchResultPeriod, error)
}

type PeriodUseCase struct {
	periodRepository repositoryCatalog.PeriodRepositoryInterface
	Logger           *logger.Logger
}

func NewPeriodUseCase(periodRepository repositoryCatalog.PeriodRepositoryInterface, logger *logger.Logger) IPeriodUseCase {
	return &PeriodUseCase{
		periodRepository: periodRepository,
		Logger:           logger,
	}
}

func (s *PeriodUseCase) GetAll() (*[]domainCatalog.Period, error) {
	s.Logger.Info("Getting all periods")
	return s.periodRepository.GetAll()
}

func (s *PeriodUseCase) GetByID(id string) (*domainCatalog.Period, error) {
	s.Logger.Info("Getting period by ID", zap.String("id", id))
	return s.periodRepository.GetByID(id)
}

func (s *PeriodUseCase) Create(newPeriod *domainCatalog.Period) (*domainCatalog.Period, error) {
	s.Logger.Info("Creating new period", zap.String("name", newPeriod.Name))
	return s.periodRepository.Create(newPeriod)
}

func (s *PeriodUseCase) Update(id string, periodMap map[string]interface{}) (*domainCatalog.Period, error) {
	s.Logger.Info("Updating period", zap.String("id", id))
	return s.periodRepository.Update(id, periodMap)
}

func (s *PeriodUseCase) Delete(id string) error {
	s.Logger.Info("Deleting period", zap.String("id", id))
	return s.periodRepository.Delete(id)
}

func (s *PeriodUseCase) SearchPaginated(filters domain.DataFilters) (*domainCatalog.SearchResultPeriod, error) {
	s.Logger.Info("Searching periods with pagination",
		zap.Int("page", filters.Page),
		zap.Int("pageSize", filters.PageSize))
	return s.periodRepository.SearchPaginated(filters)
}
