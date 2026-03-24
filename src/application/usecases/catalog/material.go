package catalog

import (
	"github.com/gbrayhan/microservices-go/src/domain"
	domainCatalog "github.com/gbrayhan/microservices-go/src/domain/catalog"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	repositoryCatalog "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/catalog"
	"go.uber.org/zap"
)

type IMaterialUseCase interface {
	GetAll() (*[]domainCatalog.Material, error)
	GetByID(id string) (*domainCatalog.Material, error)
	Create(newMaterial *domainCatalog.Material) (*domainCatalog.Material, error)
	Update(id string, materialMap map[string]interface{}) (*domainCatalog.Material, error)
	Delete(id string) error
	SearchPaginated(filters domain.DataFilters) (*domainCatalog.SearchResultMaterial, error)
}

type MaterialUseCase struct {
	materialRepository repositoryCatalog.MaterialRepositoryInterface
	Logger             *logger.Logger
}

func NewMaterialUseCase(materialRepository repositoryCatalog.MaterialRepositoryInterface, logger *logger.Logger) IMaterialUseCase {
	return &MaterialUseCase{
		materialRepository: materialRepository,
		Logger:             logger,
	}
}

func (s *MaterialUseCase) GetAll() (*[]domainCatalog.Material, error) {
	s.Logger.Info("Getting all materials")
	return s.materialRepository.GetAll()
}

func (s *MaterialUseCase) GetByID(id string) (*domainCatalog.Material, error) {
	s.Logger.Info("Getting material by ID", zap.String("id", id))
	return s.materialRepository.GetByID(id)
}

func (s *MaterialUseCase) Create(newMaterial *domainCatalog.Material) (*domainCatalog.Material, error) {
	s.Logger.Info("Creating new material", zap.String("name", newMaterial.Name))
	return s.materialRepository.Create(newMaterial)
}

func (s *MaterialUseCase) Update(id string, materialMap map[string]interface{}) (*domainCatalog.Material, error) {
	s.Logger.Info("Updating material", zap.String("id", id))
	return s.materialRepository.Update(id, materialMap)
}

func (s *MaterialUseCase) Delete(id string) error {
	s.Logger.Info("Deleting material", zap.String("id", id))
	return s.materialRepository.Delete(id)
}

func (s *MaterialUseCase) SearchPaginated(filters domain.DataFilters) (*domainCatalog.SearchResultMaterial, error) {
	s.Logger.Info("Searching materials with pagination",
		zap.Int("page", filters.Page),
		zap.Int("pageSize", filters.PageSize))
	return s.materialRepository.SearchPaginated(filters)
}
