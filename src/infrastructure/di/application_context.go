package di

import (
	"sync"

	authUseCase "github.com/gbrayhan/microservices-go/src/application/usecases/auth"
	catalogUseCase "github.com/gbrayhan/microservices-go/src/application/usecases/catalog"
	commercialUseCase "github.com/gbrayhan/microservices-go/src/application/usecases/commercial"
	operationUseCase "github.com/gbrayhan/microservices-go/src/application/usecases/operation"
	productUseCase "github.com/gbrayhan/microservices-go/src/application/usecases/product"
	userUseCase "github.com/gbrayhan/microservices-go/src/application/usecases/user"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql"
	repositoryCatalog "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/catalog"
	repositoryCommercial "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/commercial"
	repositoryOperation "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/operation"
	repositoryProduct "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/product"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/user"
	authController "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/auth"
	catalogController "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/catalog"
	commercialController "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/commercial"
	operationController "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/operation"
	productController "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/product"
	userController "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/user"
	wsController "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/websocket"
	"github.com/gbrayhan/microservices-go/src/infrastructure/security"
	ws "github.com/gbrayhan/microservices-go/src/infrastructure/websocket"
	"gorm.io/gorm"
)

// ApplicationContext holds all application dependencies and services
type ApplicationContext struct {
	DB                 *gorm.DB
	Logger             *logger.Logger
	AuthController     authController.IAuthController
	UserController    userController.IUserController
	ProductController productController.IProductController

	CategoryController catalogController.ICategoryController
	PeriodController   catalogController.IPeriodController
	StyleController    catalogController.IStyleController
	TagController      catalogController.ITagController
	MaterialController catalogController.IMaterialController

	SupplierController    commercialController.ISupplierController
	CustomerController    commercialController.ICustomerController
	AcquisitionController commercialController.IAcquisitionController
	SaleController        commercialController.ISaleController
	PaymentController     commercialController.IPaymentController

	RestorationOrderController productController.IRestorationOrderController
	ReservationController      operationController.IReservationController
	AppraisalController        operationController.IAppraisalController
	AuctionController          operationController.IAuctionController

	WebSocketController wsController.IWebSocketController
	WSHub               *ws.Hub

	JWTService        security.IJWTService
	UserRepository    user.UserRepositoryInterface
	ProductRepository repositoryProduct.ProductRepositoryInterface

	CategoryRepository repositoryCatalog.CategoryRepositoryInterface
	PeriodRepository   repositoryCatalog.PeriodRepositoryInterface
	StyleRepository    repositoryCatalog.StyleRepositoryInterface
	TagRepository      repositoryCatalog.TagRepositoryInterface
	MaterialRepository repositoryCatalog.MaterialRepositoryInterface

	SupplierRepository    repositoryCommercial.SupplierRepositoryInterface
	CustomerRepository    repositoryCommercial.CustomerRepositoryInterface
	AcquisitionRepository repositoryCommercial.AcquisitionRepositoryInterface
	SaleRepository        repositoryCommercial.SaleRepositoryInterface
	PaymentRepository     repositoryCommercial.PaymentRepositoryInterface

	RestorationOrderRepository repositoryProduct.RestorationOrderRepositoryInterface
	ReservationRepository      repositoryOperation.ReservationRepositoryInterface
	AppraisalRepository        repositoryOperation.AppraisalRepositoryInterface
	AuctionRepository          repositoryOperation.AuctionRepositoryInterface

	AuthUseCase       authUseCase.IAuthUseCase
	UserUseCase       userUseCase.IUserUseCase
	ProductUseCase    productUseCase.IProductUseCase

	CategoryUseCase catalogUseCase.ICategoryUseCase
	PeriodUseCase   catalogUseCase.IPeriodUseCase
	StyleUseCase    catalogUseCase.IStyleUseCase
	TagUseCase      catalogUseCase.ITagUseCase
	MaterialUseCase catalogUseCase.IMaterialUseCase

	SupplierUseCase    commercialUseCase.ISupplierUseCase
	CustomerUseCase    commercialUseCase.ICustomerUseCase
	AcquisitionUseCase commercialUseCase.IAcquisitionUseCase
	SaleUseCase        commercialUseCase.ISaleUseCase
	PaymentUseCase     commercialUseCase.IPaymentUseCase

	RestorationOrderUseCase productUseCase.IRestorationOrderUseCase
	ReservationUseCase      operationUseCase.IReservationUseCase
	AppraisalUseCase        operationUseCase.IAppraisalUseCase
	AuctionUseCase          operationUseCase.IAuctionUseCase
}

var (
	loggerInstance *logger.Logger
	loggerOnce     sync.Once
)

func GetLogger() *logger.Logger {
	loggerOnce.Do(func() {
		loggerInstance, _ = logger.NewLogger()
	})
	return loggerInstance
}

// SetupDependencies creates a new application context with all dependencies
func SetupDependencies(loggerInstance *logger.Logger) (*ApplicationContext, error) {
	// Initialize database with logger
	db, err := psql.InitPSQLDB(loggerInstance)
	if err != nil {
		return nil, err
	}

	// Initialize WebSocket Hub
	wsHub := ws.NewHub()
	go wsHub.Run()

	// Initialize JWT service (manages its own configuration)
	jwtService := security.NewJWTService()

	// Initialize repositories with logger
	userRepo := user.NewUserRepository(db, loggerInstance)
	productRepo := repositoryProduct.NewProductRepository(db, loggerInstance)
	categoryRepo := repositoryCatalog.NewCategoryRepository(db, loggerInstance)
	periodRepo := repositoryCatalog.NewPeriodRepository(db, loggerInstance)
	styleRepo := repositoryCatalog.NewStyleRepository(db, loggerInstance)
	tagRepo := repositoryCatalog.NewTagRepository(db, loggerInstance)
	materialRepo := repositoryCatalog.NewMaterialRepository(db, loggerInstance)

	supplierRepo := repositoryCommercial.NewSupplierRepository(db, loggerInstance)
	customerRepo := repositoryCommercial.NewCustomerRepository(db, loggerInstance)
	acquisitionRepo := repositoryCommercial.NewAcquisitionRepository(db, loggerInstance)
	saleRepo := repositoryCommercial.NewSaleRepository(db, loggerInstance)
	paymentRepo := repositoryCommercial.NewPaymentRepository(db, loggerInstance)

	restorationOrderRepo := repositoryProduct.NewRestorationOrderRepository(db, loggerInstance)
	reservationRepo := repositoryOperation.NewReservationRepository(db, loggerInstance)
	appraisalRepo := repositoryOperation.NewAppraisalRepository(db, loggerInstance)
	auctionRepo := repositoryOperation.NewAuctionRepository(db, loggerInstance)

	// Initialize use cases with logger
	authUC := authUseCase.NewAuthUseCase(userRepo, jwtService, loggerInstance)
	userUC := userUseCase.NewUserUseCase(userRepo, loggerInstance)
	productUC := productUseCase.NewProductUseCase(productRepo, loggerInstance)

	categoryUC := catalogUseCase.NewCategoryUseCase(categoryRepo, loggerInstance)
	periodUC := catalogUseCase.NewPeriodUseCase(periodRepo, loggerInstance)
	styleUC := catalogUseCase.NewStyleUseCase(styleRepo, loggerInstance)
	tagUC := catalogUseCase.NewTagUseCase(tagRepo, loggerInstance)
	materialUC := catalogUseCase.NewMaterialUseCase(materialRepo, loggerInstance)

	supplierUC := commercialUseCase.NewSupplierUseCase(supplierRepo, loggerInstance)
	customerUC := commercialUseCase.NewCustomerUseCase(customerRepo, loggerInstance)
	acquisitionUC := commercialUseCase.NewAcquisitionUseCase(acquisitionRepo, loggerInstance)
	saleUC := commercialUseCase.NewSaleUseCase(saleRepo, loggerInstance)
	paymentUC := commercialUseCase.NewPaymentUseCase(paymentRepo, loggerInstance)

	restorationOrderUC := productUseCase.NewRestorationOrderUseCase(restorationOrderRepo, loggerInstance)
	reservationUC := operationUseCase.NewReservationUseCase(reservationRepo, loggerInstance)
	appraisalUC := operationUseCase.NewAppraisalUseCase(appraisalRepo, loggerInstance)
	auctionUC := operationUseCase.NewAuctionUseCase(auctionRepo, loggerInstance)

	// Initialize controllers with logger
	authController := authController.NewAuthController(authUC, loggerInstance)
	userController := userController.NewUserController(userUC, loggerInstance)
	mainProductController := productController.NewProductController(productUC, loggerInstance)

	categoryController := catalogController.NewCategoryController(categoryUC, loggerInstance)
	periodController := catalogController.NewPeriodController(periodUC, loggerInstance)
	styleController := catalogController.NewStyleController(styleUC, loggerInstance)
	tagController := catalogController.NewTagController(tagUC, loggerInstance)
	materialController := catalogController.NewMaterialController(materialUC, loggerInstance)

	supplierController := commercialController.NewSupplierController(supplierUC, loggerInstance)
	customerController := commercialController.NewCustomerController(customerUC, loggerInstance)
	acquisitionController := commercialController.NewAcquisitionController(acquisitionUC, loggerInstance)
	saleController := commercialController.NewSaleController(saleUC, loggerInstance)
	paymentController := commercialController.NewPaymentController(paymentUC, loggerInstance)

	restorationOrderController := productController.NewRestorationOrderController(restorationOrderUC, loggerInstance)
	reservationController := operationController.NewReservationController(reservationUC, loggerInstance)
	appraisalController := operationController.NewAppraisalController(appraisalUC, loggerInstance)
	auctionController := operationController.NewAuctionController(auctionUC, loggerInstance)
	webSocketController := wsController.NewWebSocketController(wsHub, loggerInstance)

	return &ApplicationContext{
		DB:                 db,
		Logger:             loggerInstance,
		AuthController:     authController,
		UserController:     userController,
		ProductController:  mainProductController,
		CategoryController: categoryController,
		PeriodController:   periodController,
		StyleController:       styleController,
		TagController:         tagController,
		MaterialController:    materialController,
		SupplierController:    supplierController,
		CustomerController:    customerController,
		AcquisitionController: acquisitionController,
		SaleController:        saleController,
		PaymentController:     paymentController,
		RestorationOrderController: restorationOrderController,
		ReservationController:      reservationController,
		AppraisalController:        appraisalController,
		AuctionController:          auctionController,
		WebSocketController:        webSocketController,
		WSHub:                      wsHub,
		JWTService:         jwtService,
		UserRepository:     userRepo,
		ProductRepository:  productRepo,
		CategoryRepository: categoryRepo,
		PeriodRepository:   periodRepo,
		StyleRepository:       styleRepo,
		TagRepository:         tagRepo,
		MaterialRepository:    materialRepo,
		SupplierRepository:    supplierRepo,
		CustomerRepository:    customerRepo,
		AcquisitionRepository: acquisitionRepo,
		SaleRepository:        saleRepo,
		PaymentRepository:     paymentRepo,
		RestorationOrderRepository: restorationOrderRepo,
		ReservationRepository:      reservationRepo,
		AppraisalRepository:        appraisalRepo,
		AuctionRepository:          auctionRepo,
		AuthUseCase:        authUC,
		UserUseCase:        userUC,
		ProductUseCase:     productUC,
		CategoryUseCase:    categoryUC,
		PeriodUseCase:      periodUC,
		StyleUseCase:       styleUC,
		TagUseCase:         tagUC,
		MaterialUseCase:    materialUC,
		SupplierUseCase:    supplierUC,
		CustomerUseCase:    customerUC,
		AcquisitionUseCase: acquisitionUC,
		SaleUseCase:        saleUC,
		PaymentUseCase:     paymentUC,
		RestorationOrderUseCase: restorationOrderUC,
		ReservationUseCase:      reservationUC,
		AppraisalUseCase:        appraisalUC,
		AuctionUseCase:          auctionUC,
	}, nil
}

// NewTestApplicationContext creates an application context for testing with mocked dependencies
func NewTestApplicationContext(
	mockUserRepo user.UserRepositoryInterface,
	mockJWTService security.IJWTService,
	loggerInstance *logger.Logger,
) *ApplicationContext {
	// Initialize use cases with mocked repositories and logger
	authUC := authUseCase.NewAuthUseCase(mockUserRepo, mockJWTService, loggerInstance)
	userUC := userUseCase.NewUserUseCase(mockUserRepo, loggerInstance)

	// Initialize controllers with logger
	authController := authController.NewAuthController(authUC, loggerInstance)
	userController := userController.NewUserController(userUC, loggerInstance)

	return &ApplicationContext{
		Logger:         loggerInstance,
		AuthController: authController,
		UserController: userController,
		JWTService:     mockJWTService,
		UserRepository: mockUserRepo,
		AuthUseCase:    authUC,
		UserUseCase:    userUC,
	}
}
