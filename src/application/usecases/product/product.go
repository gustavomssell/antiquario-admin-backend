package product

import (
	"github.com/gbrayhan/microservices-go/src/domain"
	domainProduct "github.com/gbrayhan/microservices-go/src/domain/product"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	repositoryProduct "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/product"
	"go.uber.org/zap"
)

type IProductUseCase interface {
	GetAll() (*[]domainProduct.Product, error)
	GetByID(id string) (*domainProduct.Product, error)
	Create(newProduct *domainProduct.Product) (*domainProduct.Product, error)
	Update(id string, productMap map[string]interface{}) (*domainProduct.Product, error)
	Delete(id string) error
	SearchPaginated(filters domain.DataFilters) (*domainProduct.SearchResultProduct, error)
}

type ProductUseCase struct {
	productRepository repositoryProduct.ProductRepositoryInterface
	Logger            *logger.Logger
}

func NewProductUseCase(productRepository repositoryProduct.ProductRepositoryInterface, logger *logger.Logger) IProductUseCase {
	return &ProductUseCase{
		productRepository: productRepository,
		Logger:            logger,
	}
}

func (s *ProductUseCase) GetAll() (*[]domainProduct.Product, error) {
	s.Logger.Info("Getting all products")
	return s.productRepository.GetAll()
}

func (s *ProductUseCase) GetByID(id string) (*domainProduct.Product, error) {
	s.Logger.Info("Getting product by ID", zap.String("id", id))
	return s.productRepository.GetByID(id)
}

func (s *ProductUseCase) Create(newProduct *domainProduct.Product) (*domainProduct.Product, error) {
	s.Logger.Info("Creating new product", zap.String("code", newProduct.Code))
	
	if newProduct.Status == "" {
		newProduct.Status = "Disponível"
	}
	
	return s.productRepository.Create(newProduct)
}

func (s *ProductUseCase) Update(id string, productMap map[string]interface{}) (*domainProduct.Product, error) {
	s.Logger.Info("Updating product", zap.String("id", id))
	return s.productRepository.Update(id, productMap)
}

func (s *ProductUseCase) Delete(id string) error {
	s.Logger.Info("Deleting product", zap.String("id", id))
	return s.productRepository.Delete(id)
}

func (s *ProductUseCase) SearchPaginated(filters domain.DataFilters) (*domainProduct.SearchResultProduct, error) {
	s.Logger.Info("Searching products with pagination",
		zap.Int("page", filters.Page),
		zap.Int("pageSize", filters.PageSize))
	return s.productRepository.SearchPaginated(filters)
}
