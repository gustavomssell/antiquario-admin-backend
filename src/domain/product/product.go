package product

import "time"

type Product struct {
	ID                     string // UUID
	Code                   string // e.g., ANT-2024-0001
	QRCodeURL              string
	ProductSetID           *string // FK to ProductSet (UUID)
	SetPosition            *string // e.g., A, B, C...
	Name                   string
	Description            string
	CategoryID             *string // FK to Category (UUID)
	PeriodID               *string // FK to Period (UUID)
	StyleID                *string // FK to Style (UUID)
	Dimensions             string  // JSON format (altura, largura, profundidade)
	Weight                 float64
	ConditionRating        int     // 1-10
	AcquisitionType        string  // e.g., purchase, consignment
	AcquisitionDate        *time.Time
	AcquisitionPrice       float64
	ConsignmentPercentage  float64
	ConsignmentDeadline    *time.Time
	SupplierID             *string // FK to Supplier (UUID)
	EstimatedValue         float64
	SellingPrice           float64
	CommissionRate         float64
	Status                 string  // ENUM: available, sold, auction, reserved, restoration, evaluation, consignment, damaged, exhibition, returned
	ProvenanceStory        string  // TEXT
	HistoricalNotes        string  // TEXT
	IsSetItem              bool
	AvailableQuantity      int
	CreatedBy              int // FK to User
	UpdatedBy              int // FK to User
	CreatedAt              time.Time
	UpdatedAt              time.Time
}

type SearchResultProduct struct {
	Data       *[]Product
	Total      int64
	Page       int
	PageSize   int
	TotalPages int
}
