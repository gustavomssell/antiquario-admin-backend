package product

import (
	"github.com/gbrayhan/microservices-go/src/domain"
	domainProduct "github.com/gbrayhan/microservices-go/src/domain/product"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	repositoryProduct "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/product"
	"go.uber.org/zap"
)

type IRestorationOrderUseCase interface {
	GetAll() (*[]domainProduct.RestorationOrder, error)
	GetByID(id string) (*domainProduct.RestorationOrder, error)
	Create(newOrder *domainProduct.RestorationOrder) (*domainProduct.RestorationOrder, error)
	Update(id string, orderMap map[string]interface{}) (*domainProduct.RestorationOrder, error)
	Delete(id string) error
	SearchPaginated(filters domain.DataFilters) (*domainProduct.SearchResultRestorationOrder, error)
}

type RestorationOrderUseCase struct {
	restorationOrderRepository repositoryProduct.RestorationOrderRepositoryInterface
	Logger                     *logger.Logger
}

func NewRestorationOrderUseCase(repo repositoryProduct.RestorationOrderRepositoryInterface, logger *logger.Logger) IRestorationOrderUseCase {
	return &RestorationOrderUseCase{
		restorationOrderRepository: repo,
		Logger:                     logger,
	}
}

func (s *RestorationOrderUseCase) GetAll() (*[]domainProduct.RestorationOrder, error) {
	s.Logger.Info("Getting all restoration orders")
	return s.restorationOrderRepository.GetAll()
}

func (s *RestorationOrderUseCase) GetByID(id string) (*domainProduct.RestorationOrder, error) {
	s.Logger.Info("Getting restoration order by ID", zap.String("id", id))
	return s.restorationOrderRepository.GetByID(id)
}

func (s *RestorationOrderUseCase) Create(newOrder *domainProduct.RestorationOrder) (*domainProduct.RestorationOrder, error) {
	s.Logger.Info("Creating new restoration order", zap.String("productID", newOrder.ProductID))
	return s.restorationOrderRepository.Create(newOrder)
}

func (s *RestorationOrderUseCase) Update(id string, orderMap map[string]interface{}) (*domainProduct.RestorationOrder, error) {
	s.Logger.Info("Updating restoration order", zap.String("id", id))
	return s.restorationOrderRepository.Update(id, orderMap)
}

func (s *RestorationOrderUseCase) Delete(id string) error {
	s.Logger.Info("Deleting restoration order", zap.String("id", id))
	return s.restorationOrderRepository.Delete(id)
}

func (s *RestorationOrderUseCase) SearchPaginated(filters domain.DataFilters) (*domainProduct.SearchResultRestorationOrder, error) {
	s.Logger.Info("Searching restoration orders with pagination",
		zap.Int("page", filters.Page),
		zap.Int("pageSize", filters.PageSize))
	return s.restorationOrderRepository.SearchPaginated(filters)
}
