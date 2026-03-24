package commercial

import (
	"time"
	domainCommercial "github.com/gbrayhan/microservices-go/src/domain/commercial"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/product"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/user"
)

type Supplier struct {
	ID             string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name           string    `gorm:"column:name"`
	Type           string    `gorm:"column:type"`
	ContactInfo    string    `gorm:"column:contact_info;type:jsonb"`
	Address        string    `gorm:"column:address;type:jsonb"`
	DocumentNumber string    `gorm:"column:document_number"`
	Notes          string    `gorm:"column:notes"`
	CreatedAt      time.Time `gorm:"autoCreateTime:mili"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime:mili"`
}

func (Supplier) TableName() string { return "suppliers" }

func (s *Supplier) toDomainMapper() *domainCommercial.Supplier {
	return &domainCommercial.Supplier{
		ID:             s.ID,
		Name:           s.Name,
		Type:           s.Type,
		ContactInfo:    s.ContactInfo,
		Address:        s.Address,
		DocumentNumber: s.DocumentNumber,
		Notes:          s.Notes,
		CreatedAt:      s.CreatedAt,
		UpdatedAt:      s.UpdatedAt,
	}
}

func fromDomainSupplierMapper(s *domainCommercial.Supplier) *Supplier {
	return &Supplier{
		ID:             s.ID,
		Name:           s.Name,
		Type:           s.Type,
		ContactInfo:    s.ContactInfo,
		Address:        s.Address,
		DocumentNumber: s.DocumentNumber,
		Notes:          s.Notes,
		CreatedAt:      s.CreatedAt,
		UpdatedAt:      s.UpdatedAt,
	}
}

type Customer struct {
	ID                     string     `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID                 int        `gorm:"column:user_id;uniqueIndex"`
	User                   *user.User `gorm:"foreignKey:UserID"`
	CustomerType           string     `gorm:"column:customer_type"`
	Preferences            string     `gorm:"column:preferences;type:jsonb"`
	PurchaseHistorySummary string     `gorm:"column:purchase_history_summary;type:jsonb"`
	CreditLimit            float64    `gorm:"column:credit_limit"`
	Notes                  string     `gorm:"column:notes"`
	CreatedAt              time.Time  `gorm:"autoCreateTime:mili"`
	UpdatedAt              time.Time  `gorm:"autoUpdateTime:mili"`
}

func (Customer) TableName() string { return "customers" }

func (c *Customer) toDomainMapper() *domainCommercial.Customer {
	return &domainCommercial.Customer{
		ID:                     c.ID,
		UserID:                 c.UserID,
		CustomerType:           c.CustomerType,
		Preferences:            c.Preferences,
		PurchaseHistorySummary: c.PurchaseHistorySummary,
		CreditLimit:            c.CreditLimit,
		Notes:                  c.Notes,
		CreatedAt:              c.CreatedAt,
		UpdatedAt:              c.UpdatedAt,
	}
}

func fromDomainCustomerMapper(c *domainCommercial.Customer) *Customer {
	return &Customer{
		ID:                     c.ID,
		UserID:                 c.UserID,
		CustomerType:           c.CustomerType,
		Preferences:            c.Preferences,
		PurchaseHistorySummary: c.PurchaseHistorySummary,
		CreditLimit:            c.CreditLimit,
		Notes:                  c.Notes,
		CreatedAt:              c.CreatedAt,
		UpdatedAt:              c.UpdatedAt,
	}
}

type Acquisition struct {
	ID              string      `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	SupplierID      string      `gorm:"column:supplier_id;type:uuid"`
	Supplier        *Supplier   `gorm:"foreignKey:SupplierID"`
	AcquisitionDate *time.Time  `gorm:"column:acquisition_date"`
	TotalValue      float64     `gorm:"column:total_value"`
	PaymentMethod   string      `gorm:"column:payment_method"`
	Notes           string      `gorm:"column:notes"`
	CreatedBy       int         `gorm:"column:created_by"`
	CreatedByUser   *user.User  `gorm:"foreignKey:CreatedBy"`
	CreatedAt       time.Time   `gorm:"autoCreateTime:mili"`
	UpdatedAt       time.Time   `gorm:"autoUpdateTime:mili"`
}

func (Acquisition) TableName() string { return "acquisitions" }

func (a *Acquisition) toDomainMapper() *domainCommercial.Acquisition {
	return &domainCommercial.Acquisition{
		ID:              a.ID,
		SupplierID:      a.SupplierID,
		AcquisitionDate: a.AcquisitionDate,
		TotalValue:      a.TotalValue,
		PaymentMethod:   a.PaymentMethod,
		Notes:           a.Notes,
		CreatedBy:       a.CreatedBy,
		CreatedAt:       a.CreatedAt,
		UpdatedAt:       a.UpdatedAt,
	}
}

func fromDomainAcquisitionMapper(a *domainCommercial.Acquisition) *Acquisition {
	return &Acquisition{
		ID:              a.ID,
		SupplierID:      a.SupplierID,
		AcquisitionDate: a.AcquisitionDate,
		TotalValue:      a.TotalValue,
		PaymentMethod:   a.PaymentMethod,
		Notes:           a.Notes,
		CreatedBy:       a.CreatedBy,
		CreatedAt:       a.CreatedAt,
		UpdatedAt:       a.UpdatedAt,
	}
}

type AcquisitionItem struct {
	ID            string           `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	AcquisitionID string           `gorm:"column:acquisition_id;type:uuid"`
	Acquisition   *Acquisition     `gorm:"foreignKey:AcquisitionID"`
	ProductID     string           `gorm:"column:product_id;type:uuid"`
	Product       *product.Product `gorm:"foreignKey:ProductID"`
	UnitPrice     float64          `gorm:"column:unit_price"`
	Quantity      int              `gorm:"column:quantity"`
	TotalPrice    float64          `gorm:"column:total_price"`
}

func (AcquisitionItem) TableName() string { return "acquisition_items" }

func (a *AcquisitionItem) toDomainMapper() *domainCommercial.AcquisitionItem {
	return &domainCommercial.AcquisitionItem{
		ID:            a.ID,
		AcquisitionID: a.AcquisitionID,
		ProductID:     a.ProductID,
		UnitPrice:     a.UnitPrice,
		Quantity:      a.Quantity,
		TotalPrice:    a.TotalPrice,
	}
}

func fromDomainAcquisitionItemMapper(a *domainCommercial.AcquisitionItem) *AcquisitionItem {
	return &AcquisitionItem{
		ID:            a.ID,
		AcquisitionID: a.AcquisitionID,
		ProductID:     a.ProductID,
		UnitPrice:     a.UnitPrice,
		Quantity:      a.Quantity,
		TotalPrice:    a.TotalPrice,
	}
}

type Sale struct {
	ID              string      `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CustomerID      string      `gorm:"column:customer_id;type:uuid"`
	Customer        *Customer   `gorm:"foreignKey:CustomerID"`
	SaleDate        *time.Time  `gorm:"column:sale_date"`
	SaleType        string      `gorm:"column:sale_type"`
	TotalAmount     float64     `gorm:"column:total_amount"`
	PaymentMethod   string      `gorm:"column:payment_method"`
	PaymentStatus   string      `gorm:"column:payment_status"`
	DeliveryAddress string      `gorm:"column:delivery_address;type:jsonb"`
	DeliveryDate    *time.Time  `gorm:"column:delivery_date"`
	Notes           string      `gorm:"column:notes"`
	CreatedBy       int         `gorm:"column:created_by"`
	CreatedByUser   *user.User  `gorm:"foreignKey:CreatedBy"`
	CreatedAt       time.Time   `gorm:"autoCreateTime:mili"`
	UpdatedAt       time.Time   `gorm:"autoUpdateTime:mili"`
}

func (Sale) TableName() string { return "sales" }

func (s *Sale) toDomainMapper() *domainCommercial.Sale {
	return &domainCommercial.Sale{
		ID:              s.ID,
		CustomerID:      s.CustomerID,
		SaleDate:        s.SaleDate,
		SaleType:        s.SaleType,
		TotalAmount:     s.TotalAmount,
		PaymentMethod:   s.PaymentMethod,
		PaymentStatus:   s.PaymentStatus,
		DeliveryAddress: s.DeliveryAddress,
		DeliveryDate:    s.DeliveryDate,
		Notes:           s.Notes,
		CreatedBy:       s.CreatedBy,
		CreatedAt:       s.CreatedAt,
		UpdatedAt:       s.UpdatedAt,
	}
}

func fromDomainSaleMapper(s *domainCommercial.Sale) *Sale {
	return &Sale{
		ID:              s.ID,
		CustomerID:      s.CustomerID,
		SaleDate:        s.SaleDate,
		SaleType:        s.SaleType,
		TotalAmount:     s.TotalAmount,
		PaymentMethod:   s.PaymentMethod,
		PaymentStatus:   s.PaymentStatus,
		DeliveryAddress: s.DeliveryAddress,
		DeliveryDate:    s.DeliveryDate,
		Notes:           s.Notes,
		CreatedBy:       s.CreatedBy,
		CreatedAt:       s.CreatedAt,
		UpdatedAt:       s.UpdatedAt,
	}
}

type SaleItem struct {
	ID                   string           `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	SaleID               string           `gorm:"column:sale_id;type:uuid"`
	Sale                 *Sale            `gorm:"foreignKey:SaleID"`
	ProductID            string           `gorm:"column:product_id;type:uuid"`
	Product              *product.Product `gorm:"foreignKey:ProductID"`
	UnitPrice            float64          `gorm:"column:unit_price"`
	Quantity             int              `gorm:"column:quantity"`
	TotalPrice           float64          `gorm:"column:total_price"`
	CommissionPercentage float64          `gorm:"column:commission_percentage"`
	CommissionAmount     float64          `gorm:"column:commission_amount"`
	CommissionPaid       bool             `gorm:"column:commission_paid"`
}

func (SaleItem) TableName() string { return "sale_items" }

func (si *SaleItem) toDomainMapper() *domainCommercial.SaleItem {
	return &domainCommercial.SaleItem{
		ID:                   si.ID,
		SaleID:               si.SaleID,
		ProductID:            si.ProductID,
		UnitPrice:            si.UnitPrice,
		Quantity:             si.Quantity,
		TotalPrice:           si.TotalPrice,
		CommissionPercentage: si.CommissionPercentage,
		CommissionAmount:     si.CommissionAmount,
		CommissionPaid:       si.CommissionPaid,
	}
}

func fromDomainSaleItemMapper(si *domainCommercial.SaleItem) *SaleItem {
	return &SaleItem{
		ID:                   si.ID,
		SaleID:               si.SaleID,
		ProductID:            si.ProductID,
		UnitPrice:            si.UnitPrice,
		Quantity:             si.Quantity,
		TotalPrice:           si.TotalPrice,
		CommissionPercentage: si.CommissionPercentage,
		CommissionAmount:     si.CommissionAmount,
		CommissionPaid:       si.CommissionPaid,
	}
}

type Payment struct {
	ID            string     `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	SaleID        string     `gorm:"column:sale_id;type:uuid"`
	Sale          *Sale      `gorm:"foreignKey:SaleID"`
	Amount        float64    `gorm:"column:amount"`
	PaymentDate   *time.Time `gorm:"column:payment_date"`
	PaymentMethod string     `gorm:"column:payment_method"`
	Status        string     `gorm:"column:status"`
	Reference     string     `gorm:"column:reference"`
	CreatedAt     time.Time  `gorm:"autoCreateTime:mili"`
	UpdatedAt     time.Time  `gorm:"autoUpdateTime:mili"`
}

func (Payment) TableName() string { return "payments" }

func (p *Payment) toDomainMapper() *domainCommercial.Payment {
	return &domainCommercial.Payment{
		ID:            p.ID,
		SaleID:        p.SaleID,
		Amount:        p.Amount,
		PaymentDate:   p.PaymentDate,
		PaymentMethod: p.PaymentMethod,
		Status:        p.Status,
		Reference:     p.Reference,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}
}

func fromDomainPaymentMapper(p *domainCommercial.Payment) *Payment {
	return &Payment{
		ID:            p.ID,
		SaleID:        p.SaleID,
		Amount:        p.Amount,
		PaymentDate:   p.PaymentDate,
		PaymentMethod: p.PaymentMethod,
		Status:        p.Status,
		Reference:     p.Reference,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}
}
