package routes

import (
	"net/http"

	"github.com/gbrayhan/microservices-go/src/infrastructure/di"
	"github.com/gin-gonic/gin"
)

func ApplicationRouter(router *gin.Engine, appContext *di.ApplicationContext) {
	v1 := router.Group("/v1")

	v1.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Service is running",
		})
	})

	v1.GET("/ws", appContext.WebSocketController.ServeWS)

	AuthRoutes(v1, appContext.AuthController)
	UserRoutes(v1, appContext.UserController)

	productGroup := v1.Group("/products")
	{
		productGroup.GET("/", appContext.ProductController.GetAllProducts)
		productGroup.GET("/:id", appContext.ProductController.GetProductByID)
		productGroup.POST("/", appContext.ProductController.NewProduct)
		productGroup.PUT("/:id", appContext.ProductController.UpdateProduct)
		productGroup.DELETE("/:id", appContext.ProductController.DeleteProduct)
		productGroup.GET("/search", appContext.ProductController.SearchPaginated)
	}

	categoryGroup := v1.Group("/categories")
	{
		categoryGroup.GET("/", appContext.CategoryController.GetAllCategories)
		categoryGroup.GET("/:id", appContext.CategoryController.GetCategoryByID)
		categoryGroup.POST("/", appContext.CategoryController.NewCategory)
		categoryGroup.PUT("/:id", appContext.CategoryController.UpdateCategory)
		categoryGroup.DELETE("/:id", appContext.CategoryController.DeleteCategory)
		categoryGroup.GET("/search", appContext.CategoryController.SearchPaginated)
	}

	periodGroup := v1.Group("/periods")
	{
		periodGroup.GET("/", appContext.PeriodController.GetAllPeriods)
		periodGroup.GET("/:id", appContext.PeriodController.GetPeriodByID)
		periodGroup.POST("/", appContext.PeriodController.NewPeriod)
		periodGroup.PUT("/:id", appContext.PeriodController.UpdatePeriod)
		periodGroup.DELETE("/:id", appContext.PeriodController.DeletePeriod)
		periodGroup.GET("/search", appContext.PeriodController.SearchPaginated)
	}

	styleGroup := v1.Group("/styles")
	{
		styleGroup.GET("/", appContext.StyleController.GetAllStyles)
		styleGroup.GET("/:id", appContext.StyleController.GetStyleByID)
		styleGroup.POST("/", appContext.StyleController.NewStyle)
		styleGroup.PUT("/:id", appContext.StyleController.UpdateStyle)
		styleGroup.DELETE("/:id", appContext.StyleController.DeleteStyle)
		styleGroup.GET("/search", appContext.StyleController.SearchPaginated)
	}

	tagGroup := v1.Group("/tags")
	{
		tagGroup.GET("/", appContext.TagController.GetAllTags)
		tagGroup.GET("/:id", appContext.TagController.GetTagByID)
		tagGroup.POST("/", appContext.TagController.NewTag)
		tagGroup.PUT("/:id", appContext.TagController.UpdateTag)
		tagGroup.DELETE("/:id", appContext.TagController.DeleteTag)
		tagGroup.GET("/search", appContext.TagController.SearchPaginated)
	}

	materialGroup := v1.Group("/materials")
	{
		materialGroup.GET("/", appContext.MaterialController.GetAllMaterials)
		materialGroup.GET("/:id", appContext.MaterialController.GetMaterialByID)
		materialGroup.POST("/", appContext.MaterialController.NewMaterial)
		materialGroup.PUT("/:id", appContext.MaterialController.UpdateMaterial)
		materialGroup.DELETE("/:id", appContext.MaterialController.DeleteMaterial)
		materialGroup.GET("/search", appContext.MaterialController.SearchPaginated)
	}

	supplierGroup := v1.Group("/suppliers")
	{
		supplierGroup.GET("/", appContext.SupplierController.GetAllSuppliers)
		supplierGroup.GET("/:id", appContext.SupplierController.GetSupplierByID)
		supplierGroup.POST("/", appContext.SupplierController.NewSupplier)
		supplierGroup.PUT("/:id", appContext.SupplierController.UpdateSupplier)
		supplierGroup.DELETE("/:id", appContext.SupplierController.DeleteSupplier)
		supplierGroup.GET("/search", appContext.SupplierController.SearchPaginated)
	}

	customerGroup := v1.Group("/customers")
	{
		customerGroup.GET("/", appContext.CustomerController.GetAllCustomers)
		customerGroup.GET("/:id", appContext.CustomerController.GetCustomerByID)
		customerGroup.POST("/", appContext.CustomerController.NewCustomer)
		customerGroup.PUT("/:id", appContext.CustomerController.UpdateCustomer)
		customerGroup.DELETE("/:id", appContext.CustomerController.DeleteCustomer)
		customerGroup.GET("/search", appContext.CustomerController.SearchPaginated)
	}

	acquisitionGroup := v1.Group("/acquisitions")
	{
		acquisitionGroup.GET("/", appContext.AcquisitionController.GetAllAcquisitions)
		acquisitionGroup.GET("/:id", appContext.AcquisitionController.GetAcquisitionByID)
		acquisitionGroup.POST("/", appContext.AcquisitionController.NewAcquisition)
		acquisitionGroup.PUT("/:id", appContext.AcquisitionController.UpdateAcquisition)
		acquisitionGroup.DELETE("/:id", appContext.AcquisitionController.DeleteAcquisition)
		acquisitionGroup.GET("/search", appContext.AcquisitionController.SearchPaginated)
	}

	saleGroup := v1.Group("/sales")
	{
		saleGroup.GET("/", appContext.SaleController.GetAllSales)
		saleGroup.GET("/:id", appContext.SaleController.GetSaleByID)
		saleGroup.POST("/", appContext.SaleController.NewSale)
		saleGroup.PUT("/:id", appContext.SaleController.UpdateSale)
		saleGroup.DELETE("/:id", appContext.SaleController.DeleteSale)
		saleGroup.GET("/search", appContext.SaleController.SearchPaginated)
	}

	paymentGroup := v1.Group("/payments")
	{
		paymentGroup.GET("/", appContext.PaymentController.GetAllPayments)
		paymentGroup.GET("/:id", appContext.PaymentController.GetPaymentByID)
		paymentGroup.POST("/", appContext.PaymentController.NewPayment)
		paymentGroup.PUT("/:id", appContext.PaymentController.UpdatePayment)
		paymentGroup.DELETE("/:id", appContext.PaymentController.DeletePayment)
		paymentGroup.GET("/search", appContext.PaymentController.SearchPaginated)
	}

	restorationGroup := v1.Group("/restorations")
	{
		restorationGroup.GET("/", appContext.RestorationOrderController.GetAllRestorationOrders)
		restorationGroup.GET("/:id", appContext.RestorationOrderController.GetRestorationOrderByID)
		restorationGroup.POST("/", appContext.RestorationOrderController.NewRestorationOrder)
		restorationGroup.PUT("/:id", appContext.RestorationOrderController.UpdateRestorationOrder)
		restorationGroup.DELETE("/:id", appContext.RestorationOrderController.DeleteRestorationOrder)
		restorationGroup.GET("/search", appContext.RestorationOrderController.SearchPaginated)
	}

	reservationGroup := v1.Group("/reservations")
	{
		reservationGroup.GET("/", appContext.ReservationController.GetAllReservations)
		reservationGroup.GET("/:id", appContext.ReservationController.GetReservationByID)
		reservationGroup.POST("/", appContext.ReservationController.NewReservation)
		reservationGroup.PUT("/:id", appContext.ReservationController.UpdateReservation)
		reservationGroup.DELETE("/:id", appContext.ReservationController.DeleteReservation)
		reservationGroup.GET("/search", appContext.ReservationController.SearchPaginated)
	}

	appraisalGroup := v1.Group("/appraisals")
	{
		appraisalGroup.GET("/", appContext.AppraisalController.GetAllAppraisals)
		appraisalGroup.GET("/:id", appContext.AppraisalController.GetAppraisalByID)
		appraisalGroup.POST("/", appContext.AppraisalController.NewAppraisal)
		appraisalGroup.PUT("/:id", appContext.AppraisalController.UpdateAppraisal)
		appraisalGroup.DELETE("/:id", appContext.AppraisalController.DeleteAppraisal)
		appraisalGroup.GET("/search", appContext.AppraisalController.SearchPaginated)
	}

	auctionGroup := v1.Group("/auctions")
	{
		auctionGroup.GET("/", appContext.AuctionController.GetAllAuctions)
		auctionGroup.GET("/:id", appContext.AuctionController.GetAuctionByID)
		auctionGroup.POST("/", appContext.AuctionController.NewAuction)
		auctionGroup.PUT("/:id", appContext.AuctionController.UpdateAuction)
		auctionGroup.DELETE("/:id", appContext.AuctionController.DeleteAuction)
		auctionGroup.GET("/search", appContext.AuctionController.SearchPaginated)
	}
}
