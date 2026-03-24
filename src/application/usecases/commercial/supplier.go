package commercial

import (
	"github.com/gbrayhan/microservices-go/src/domain"
	domainCommercial "github.com/gbrayhan/microservices-go/src/domain/commercial"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	repositoryCommercial "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/commercial"
	"go.uber.org/zap"
)

type ISupplierUseCase interface {
	GetAll() (*[]domainCommercial.Supplier, error)
	GetByID(id string) (*domainCommercial.Supplier, error)
	Create(newSupplier *domainCommercial.Supplier) (*domainCommercial.Supplier, error)
	Update(id string, supplierMap map[string]interface{}) (*domainCommercial.Supplier, error)
	Delete(id string) error
	SearchPaginated(filters domain.DataFilters) (*domainCommercial.SearchResultSupplier, error)
}

type SupplierUseCase struct {
	supplierRepository repositoryCommercial.SupplierRepositoryInterface
	Logger             *logger.Logger
}

func NewSupplierUseCase(supplierRepository repositoryCommercial.SupplierRepositoryInterface, logger *logger.Logger) ISupplierUseCase {
	return &SupplierUseCase{
		supplierRepository: supplierRepository,
		Logger:             logger,
	}
}

func (s *SupplierUseCase) GetAll() (*[]domainCommercial.Supplier, error) {
	s.Logger.Info("Getting all suppliers")
	return s.supplierRepository.GetAll()
}

func (s *SupplierUseCase) GetByID(id string) (*domainCommercial.Supplier, error) {
	s.Logger.Info("Getting supplier by ID", zap.String("id", id))
	return s.supplierRepository.GetByID(id)
}

func (s *SupplierUseCase) Create(newSupplier *domainCommercial.Supplier) (*domainCommercial.Supplier, error) {
	s.Logger.Info("Creating new supplier", zap.String("name", newSupplier.Name))
	return s.supplierRepository.Create(newSupplier)
}

func (s *SupplierUseCase) Update(id string, supplierMap map[string]interface{}) (*domainCommercial.Supplier, error) {
	s.Logger.Info("Updating supplier", zap.String("id", id))
	return s.supplierRepository.Update(id, supplierMap)
}

func (s *SupplierUseCase) Delete(id string) error {
	s.Logger.Info("Deleting supplier", zap.String("id", id))
	return s.supplierRepository.Delete(id)
}

func (s *SupplierUseCase) SearchPaginated(filters domain.DataFilters) (*domainCommercial.SearchResultSupplier, error) {
	s.Logger.Info("Searching suppliers with pagination",
		zap.Int("page", filters.Page),
		zap.Int("pageSize", filters.PageSize))
	return s.supplierRepository.SearchPaginated(filters)
}
