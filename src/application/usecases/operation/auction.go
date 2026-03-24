package operation

import (
	"github.com/gbrayhan/microservices-go/src/domain"
	domainOperation "github.com/gbrayhan/microservices-go/src/domain/operation"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	repositoryOperation "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/operation"
	"go.uber.org/zap"
)

type IAuctionUseCase interface {
	GetAll() (*[]domainOperation.Auction, error)
	GetByID(id string) (*domainOperation.Auction, error)
	Create(newAuction *domainOperation.Auction) (*domainOperation.Auction, error)
	Update(id string, auctionMap map[string]interface{}) (*domainOperation.Auction, error)
	Delete(id string) error
	SearchPaginated(filters domain.DataFilters) (*domainOperation.SearchResultAuction, error)
}

type AuctionUseCase struct {
	auctionRepository repositoryOperation.AuctionRepositoryInterface
	Logger            *logger.Logger
}

func NewAuctionUseCase(auctionRepository repositoryOperation.AuctionRepositoryInterface, logger *logger.Logger) IAuctionUseCase {
	return &AuctionUseCase{
		auctionRepository: auctionRepository,
		Logger:            logger,
	}
}

func (s *AuctionUseCase) GetAll() (*[]domainOperation.Auction, error) {
	s.Logger.Info("Getting all auctions")
	return s.auctionRepository.GetAll()
}

func (s *AuctionUseCase) GetByID(id string) (*domainOperation.Auction, error) {
	s.Logger.Info("Getting auction by ID", zap.String("id", id))
	return s.auctionRepository.GetByID(id)
}

func (s *AuctionUseCase) Create(newAuction *domainOperation.Auction) (*domainOperation.Auction, error) {
	s.Logger.Info("Creating new auction", zap.String("title", newAuction.Title))
	return s.auctionRepository.Create(newAuction)
}

func (s *AuctionUseCase) Update(id string, auctionMap map[string]interface{}) (*domainOperation.Auction, error) {
	s.Logger.Info("Updating auction", zap.String("id", id))
	return s.auctionRepository.Update(id, auctionMap)
}

func (s *AuctionUseCase) Delete(id string) error {
	s.Logger.Info("Deleting auction", zap.String("id", id))
	return s.auctionRepository.Delete(id)
}

func (s *AuctionUseCase) SearchPaginated(filters domain.DataFilters) (*domainOperation.SearchResultAuction, error) {
	s.Logger.Info("Searching auctions with pagination",
		zap.Int("page", filters.Page),
		zap.Int("pageSize", filters.PageSize))
	return s.auctionRepository.SearchPaginated(filters)
}
