package operation

import (
	"github.com/gbrayhan/microservices-go/src/domain"
	domainOperation "github.com/gbrayhan/microservices-go/src/domain/operation"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	repositoryOperation "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/operation"
	"go.uber.org/zap"
)

type IReservationUseCase interface {
	GetAll() (*[]domainOperation.Reservation, error)
	GetByID(id string) (*domainOperation.Reservation, error)
	Create(newReservation *domainOperation.Reservation) (*domainOperation.Reservation, error)
	Update(id string, reservationMap map[string]interface{}) (*domainOperation.Reservation, error)
	Delete(id string) error
	SearchPaginated(filters domain.DataFilters) (*domainOperation.SearchResultReservation, error)
}

type ReservationUseCase struct {
	reservationRepository repositoryOperation.ReservationRepositoryInterface
	Logger                *logger.Logger
}

func NewReservationUseCase(reservationRepository repositoryOperation.ReservationRepositoryInterface, logger *logger.Logger) IReservationUseCase {
	return &ReservationUseCase{
		reservationRepository: reservationRepository,
		Logger:                logger,
	}
}

func (s *ReservationUseCase) GetAll() (*[]domainOperation.Reservation, error) {
	s.Logger.Info("Getting all reservations")
	return s.reservationRepository.GetAll()
}

func (s *ReservationUseCase) GetByID(id string) (*domainOperation.Reservation, error) {
	s.Logger.Info("Getting reservation by ID", zap.String("id", id))
	return s.reservationRepository.GetByID(id)
}

func (s *ReservationUseCase) Create(newReservation *domainOperation.Reservation) (*domainOperation.Reservation, error) {
	s.Logger.Info("Creating new reservation", zap.String("productID", newReservation.ProductID))
	return s.reservationRepository.Create(newReservation)
}

func (s *ReservationUseCase) Update(id string, reservationMap map[string]interface{}) (*domainOperation.Reservation, error) {
	s.Logger.Info("Updating reservation", zap.String("id", id))
	return s.reservationRepository.Update(id, reservationMap)
}

func (s *ReservationUseCase) Delete(id string) error {
	s.Logger.Info("Deleting reservation", zap.String("id", id))
	return s.reservationRepository.Delete(id)
}

func (s *ReservationUseCase) SearchPaginated(filters domain.DataFilters) (*domainOperation.SearchResultReservation, error) {
	s.Logger.Info("Searching reservations with pagination",
		zap.Int("page", filters.Page),
		zap.Int("pageSize", filters.PageSize))
	return s.reservationRepository.SearchPaginated(filters)
}
