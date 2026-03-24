package operation

import (
	"github.com/gbrayhan/microservices-go/src/domain"
	domainOperation "github.com/gbrayhan/microservices-go/src/domain/operation"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	repositoryOperation "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/operation"
	"go.uber.org/zap"
)

type IAppraisalUseCase interface {
	GetAll() (*[]domainOperation.Appraisal, error)
	GetByID(id string) (*domainOperation.Appraisal, error)
	Create(newAppraisal *domainOperation.Appraisal) (*domainOperation.Appraisal, error)
	Update(id string, appraisalMap map[string]interface{}) (*domainOperation.Appraisal, error)
	Delete(id string) error
	SearchPaginated(filters domain.DataFilters) (*domainOperation.SearchResultAppraisal, error)
}

type AppraisalUseCase struct {
	appraisalRepository repositoryOperation.AppraisalRepositoryInterface
	Logger              *logger.Logger
}

func NewAppraisalUseCase(appraisalRepository repositoryOperation.AppraisalRepositoryInterface, logger *logger.Logger) IAppraisalUseCase {
	return &AppraisalUseCase{
		appraisalRepository: appraisalRepository,
		Logger:              logger,
	}
}

func (s *AppraisalUseCase) GetAll() (*[]domainOperation.Appraisal, error) {
	s.Logger.Info("Getting all appraisals")
	return s.appraisalRepository.GetAll()
}

func (s *AppraisalUseCase) GetByID(id string) (*domainOperation.Appraisal, error) {
	s.Logger.Info("Getting appraisal by ID", zap.String("id", id))
	return s.appraisalRepository.GetByID(id)
}

func (s *AppraisalUseCase) Create(newAppraisal *domainOperation.Appraisal) (*domainOperation.Appraisal, error) {
	s.Logger.Info("Creating new appraisal", zap.String("productID", newAppraisal.ProductID))
	return s.appraisalRepository.Create(newAppraisal)
}

func (s *AppraisalUseCase) Update(id string, appraisalMap map[string]interface{}) (*domainOperation.Appraisal, error) {
	s.Logger.Info("Updating appraisal", zap.String("id", id))
	return s.appraisalRepository.Update(id, appraisalMap)
}

func (s *AppraisalUseCase) Delete(id string) error {
	s.Logger.Info("Deleting appraisal", zap.String("id", id))
	return s.appraisalRepository.Delete(id)
}

func (s *AppraisalUseCase) SearchPaginated(filters domain.DataFilters) (*domainOperation.SearchResultAppraisal, error) {
	s.Logger.Info("Searching appraisals with pagination",
		zap.Int("page", filters.Page),
		zap.Int("pageSize", filters.PageSize))
	return s.appraisalRepository.SearchPaginated(filters)
}
