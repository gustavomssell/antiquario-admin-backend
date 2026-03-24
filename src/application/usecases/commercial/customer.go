package commercial

import (
	"github.com/gbrayhan/microservices-go/src/domain"
	domainCommercial "github.com/gbrayhan/microservices-go/src/domain/commercial"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	repositoryCommercial "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/commercial"
	"go.uber.org/zap"
)

type ICustomerUseCase interface {
	GetAll() (*[]domainCommercial.Customer, error)
	GetByID(id string) (*domainCommercial.Customer, error)
	Create(newCustomer *domainCommercial.Customer) (*domainCommercial.Customer, error)
	Update(id string, customerMap map[string]interface{}) (*domainCommercial.Customer, error)
	Delete(id string) error
	SearchPaginated(filters domain.DataFilters) (*domainCommercial.SearchResultCustomer, error)
}

type CustomerUseCase struct {
	customerRepository repositoryCommercial.CustomerRepositoryInterface
	Logger             *logger.Logger
}

func NewCustomerUseCase(customerRepository repositoryCommercial.CustomerRepositoryInterface, logger *logger.Logger) ICustomerUseCase {
	return &CustomerUseCase{
		customerRepository: customerRepository,
		Logger:             logger,
	}
}

func (s *CustomerUseCase) GetAll() (*[]domainCommercial.Customer, error) {
	s.Logger.Info("Getting all customers")
	return s.customerRepository.GetAll()
}

func (s *CustomerUseCase) GetByID(id string) (*domainCommercial.Customer, error) {
	s.Logger.Info("Getting customer by ID", zap.String("id", id))
	return s.customerRepository.GetByID(id)
}

func (s *CustomerUseCase) Create(newCustomer *domainCommercial.Customer) (*domainCommercial.Customer, error) {
	s.Logger.Info("Creating new customer", zap.Int("userID", newCustomer.UserID))
	return s.customerRepository.Create(newCustomer)
}

func (s *CustomerUseCase) Update(id string, customerMap map[string]interface{}) (*domainCommercial.Customer, error) {
	s.Logger.Info("Updating customer", zap.String("id", id))
	return s.customerRepository.Update(id, customerMap)
}

func (s *CustomerUseCase) Delete(id string) error {
	s.Logger.Info("Deleting customer", zap.String("id", id))
	return s.customerRepository.Delete(id)
}

func (s *CustomerUseCase) SearchPaginated(filters domain.DataFilters) (*domainCommercial.SearchResultCustomer, error) {
	s.Logger.Info("Searching customers with pagination",
		zap.Int("page", filters.Page),
		zap.Int("pageSize", filters.PageSize))
	return s.customerRepository.SearchPaginated(filters)
}
