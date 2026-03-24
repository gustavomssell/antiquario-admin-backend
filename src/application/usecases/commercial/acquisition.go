package commercial

import (
	"github.com/gbrayhan/microservices-go/src/domain"
	domainCommercial "github.com/gbrayhan/microservices-go/src/domain/commercial"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	repositoryCommercial "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/commercial"
	"go.uber.org/zap"
)

type IAcquisitionUseCase interface {
	GetAll() (*[]domainCommercial.Acquisition, error)
	GetByID(id string) (*domainCommercial.Acquisition, error)
	Create(newAcquisition *domainCommercial.Acquisition) (*domainCommercial.Acquisition, error)
	Update(id string, acquisitionMap map[string]interface{}) (*domainCommercial.Acquisition, error)
	Delete(id string) error
	SearchPaginated(filters domain.DataFilters) (*domainCommercial.SearchResultAcquisition, error)
}

type AcquisitionUseCase struct {
	acquisitionRepository repositoryCommercial.AcquisitionRepositoryInterface
	Logger                *logger.Logger
}

func NewAcquisitionUseCase(acquisitionRepository repositoryCommercial.AcquisitionRepositoryInterface, logger *logger.Logger) IAcquisitionUseCase {
	return &AcquisitionUseCase{
		acquisitionRepository: acquisitionRepository,
		Logger:                logger,
	}
}

func (s *AcquisitionUseCase) GetAll() (*[]domainCommercial.Acquisition, error) {
	s.Logger.Info("Getting all acquisitions")
	return s.acquisitionRepository.GetAll()
}

func (s *AcquisitionUseCase) GetByID(id string) (*domainCommercial.Acquisition, error) {
	s.Logger.Info("Getting acquisition by ID", zap.String("id", id))
	return s.acquisitionRepository.GetByID(id)
}

func (s *AcquisitionUseCase) Create(newAcquisition *domainCommercial.Acquisition) (*domainCommercial.Acquisition, error) {
	s.Logger.Info("Creating new acquisition", zap.String("supplierID", newAcquisition.SupplierID))
	return s.acquisitionRepository.Create(newAcquisition)
}

func (s *AcquisitionUseCase) Update(id string, acquisitionMap map[string]interface{}) (*domainCommercial.Acquisition, error) {
	s.Logger.Info("Updating acquisition", zap.String("id", id))
	return s.acquisitionRepository.Update(id, acquisitionMap)
}

func (s *AcquisitionUseCase) Delete(id string) error {
	s.Logger.Info("Deleting acquisition", zap.String("id", id))
	return s.acquisitionRepository.Delete(id)
}

func (s *AcquisitionUseCase) SearchPaginated(filters domain.DataFilters) (*domainCommercial.SearchResultAcquisition, error) {
	s.Logger.Info("Searching acquisitions with pagination",
		zap.Int("page", filters.Page),
		zap.Int("pageSize", filters.PageSize))
	return s.acquisitionRepository.SearchPaginated(filters)
}
