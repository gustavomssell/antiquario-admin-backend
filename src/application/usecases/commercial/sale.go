package commercial

import (
	"github.com/gbrayhan/microservices-go/src/domain"
	domainCommercial "github.com/gbrayhan/microservices-go/src/domain/commercial"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	repositoryCommercial "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/commercial"
	"go.uber.org/zap"
)

type ISaleUseCase interface {
	GetAll() (*[]domainCommercial.Sale, error)
	GetByID(id string) (*domainCommercial.Sale, error)
	Create(newSale *domainCommercial.Sale) (*domainCommercial.Sale, error)
	Update(id string, saleMap map[string]interface{}) (*domainCommercial.Sale, error)
	Delete(id string) error
	SearchPaginated(filters domain.DataFilters) (*domainCommercial.SearchResultSale, error)
}

type SaleUseCase struct {
	saleRepository repositoryCommercial.SaleRepositoryInterface
	Logger         *logger.Logger
}

func NewSaleUseCase(saleRepository repositoryCommercial.SaleRepositoryInterface, logger *logger.Logger) ISaleUseCase {
	return &SaleUseCase{
		saleRepository: saleRepository,
		Logger:         logger,
	}
}

func (s *SaleUseCase) GetAll() (*[]domainCommercial.Sale, error) {
	s.Logger.Info("Getting all sales")
	return s.saleRepository.GetAll()
}

func (s *SaleUseCase) GetByID(id string) (*domainCommercial.Sale, error) {
	s.Logger.Info("Getting sale by ID", zap.String("id", id))
	return s.saleRepository.GetByID(id)
}

func (s *SaleUseCase) Create(newSale *domainCommercial.Sale) (*domainCommercial.Sale, error) {
	s.Logger.Info("Creating new sale", zap.String("customerID", newSale.CustomerID))
	return s.saleRepository.Create(newSale)
}

func (s *SaleUseCase) Update(id string, saleMap map[string]interface{}) (*domainCommercial.Sale, error) {
	s.Logger.Info("Updating sale", zap.String("id", id))
	return s.saleRepository.Update(id, saleMap)
}

func (s *SaleUseCase) Delete(id string) error {
	s.Logger.Info("Deleting sale", zap.String("id", id))
	return s.saleRepository.Delete(id)
}

func (s *SaleUseCase) SearchPaginated(filters domain.DataFilters) (*domainCommercial.SearchResultSale, error) {
	s.Logger.Info("Searching sales with pagination",
		zap.Int("page", filters.Page),
		zap.Int("pageSize", filters.PageSize))
	return s.saleRepository.SearchPaginated(filters)
}
