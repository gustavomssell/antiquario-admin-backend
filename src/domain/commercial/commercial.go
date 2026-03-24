package commercial

import "time"

type Supplier struct {
	ID             string
	Name           string
	Type           string // e.g., pessoa_fisica, pessoa_juridica
	ContactInfo    string // JSON
	Address        string // JSON
	DocumentNumber string
	Notes          string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type Customer struct {
	ID                     string
	UserID                 int    // FK to User
	CustomerType           string // e.g., individual, empresa
	Preferences            string // JSON
	PurchaseHistorySummary string // JSON
	CreditLimit            float64
	Notes                  string
	CreatedAt              time.Time
	UpdatedAt              time.Time
}

type Acquisition struct {
	ID              string
	SupplierID      string // FK to Supplier
	AcquisitionDate *time.Time
	TotalValue      float64
	PaymentMethod   string
	Notes           string
	CreatedBy       int // FK to User
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type AcquisitionItem struct {
	ID            string
	AcquisitionID string // FK to Acquisition
	ProductID     string // FK to Product
	UnitPrice     float64
	Quantity      int
	TotalPrice    float64
}

type Sale struct {
	ID              string
	CustomerID      string // FK to Customer
	SaleDate        *time.Time
	SaleType        string // e.g., direct, auction, set
	TotalAmount     float64
	PaymentMethod   string
	PaymentStatus   string // e.g., Pago, Pendente, Atrasado
	DeliveryAddress string // JSON
	DeliveryDate    *time.Time
	Notes           string
	CreatedBy       int // FK to User (Vendedor)
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type SaleItem struct {
	ID                   string
	SaleID               string // FK to Sale
	ProductID            string // FK to Product
	UnitPrice            float64
	Quantity             int
	TotalPrice           float64
	CommissionPercentage float64
	CommissionAmount     float64
	CommissionPaid       bool
}

type Payment struct {
	ID            string
	SaleID        string // FK to Sale
	Amount        float64
	PaymentDate   *time.Time
	PaymentMethod string
	Status        string // e.g., pago, pendente, atrasado
	Reference     string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type SearchResultSupplier struct {
	Data       *[]Supplier
	Total      int64
	Page       int
	PageSize   int
	TotalPages int
}

type SearchResultCustomer struct {
	Data       *[]Customer
	Total      int64
	Page       int
	PageSize   int
	TotalPages int
}

type SearchResultAcquisition struct {
	Data       *[]Acquisition
	Total      int64
	Page       int
	PageSize   int
	TotalPages int
}

type SearchResultSale struct {
	Data       *[]Sale
	Total      int64
	Page       int
	PageSize   int
	TotalPages int
}

type SearchResultPayment struct {
	Data       *[]Payment
	Total      int64
	Page       int
	PageSize   int
	TotalPages int
}
