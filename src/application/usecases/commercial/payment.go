package commercial

import (
	"github.com/gbrayhan/microservices-go/src/domain"
	domainCommercial "github.com/gbrayhan/microservices-go/src/domain/commercial"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	repositoryCommercial "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/commercial"
	"go.uber.org/zap"
)

type IPaymentUseCase interface {
	GetAll() (*[]domainCommercial.Payment, error)
	GetByID(id string) (*domainCommercial.Payment, error)
	Create(newPayment *domainCommercial.Payment) (*domainCommercial.Payment, error)
	Update(id string, paymentMap map[string]interface{}) (*domainCommercial.Payment, error)
	Delete(id string) error
	SearchPaginated(filters domain.DataFilters) (*domainCommercial.SearchResultPayment, error)
}

type PaymentUseCase struct {
	paymentRepository repositoryCommercial.PaymentRepositoryInterface
	Logger            *logger.Logger
}

func NewPaymentUseCase(paymentRepository repositoryCommercial.PaymentRepositoryInterface, logger *logger.Logger) IPaymentUseCase {
	return &PaymentUseCase{
		paymentRepository: paymentRepository,
		Logger:            logger,
	}
}

func (s *PaymentUseCase) GetAll() (*[]domainCommercial.Payment, error) {
	s.Logger.Info("Getting all payments")
	return s.paymentRepository.GetAll()
}

func (s *PaymentUseCase) GetByID(id string) (*domainCommercial.Payment, error) {
	s.Logger.Info("Getting payment by ID", zap.String("id", id))
	return s.paymentRepository.GetByID(id)
}

func (s *PaymentUseCase) Create(newPayment *domainCommercial.Payment) (*domainCommercial.Payment, error) {
	s.Logger.Info("Creating new payment", zap.String("saleID", newPayment.SaleID))
	return s.paymentRepository.Create(newPayment)
}

func (s *PaymentUseCase) Update(id string, paymentMap map[string]interface{}) (*domainCommercial.Payment, error) {
	s.Logger.Info("Updating payment", zap.String("id", id))
	return s.paymentRepository.Update(id, paymentMap)
}

func (s *PaymentUseCase) Delete(id string) error {
	s.Logger.Info("Deleting payment", zap.String("id", id))
	return s.paymentRepository.Delete(id)
}

func (s *PaymentUseCase) SearchPaginated(filters domain.DataFilters) (*domainCommercial.SearchResultPayment, error) {
	s.Logger.Info("Searching payments with pagination",
		zap.Int("page", filters.Page),
		zap.Int("pageSize", filters.PageSize))
	return s.paymentRepository.SearchPaginated(filters)
}
