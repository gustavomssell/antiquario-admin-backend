package operation

import "time"

type CommissionRule struct {
	ID                   string
	UserID               *int    // FK user (Vendedor)
	ProductCategoryID    *string // FK category
	CommissionPercentage float64
	IsDefault            bool
	Active               bool
	CreatedBy            int
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

type ConsignmentSetting struct {
	ID                  string
	DefaultPercentage   float64
	DefaultDeadlineDays int
	CreatedBy           int
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type ConsignmentReturn struct {
	ID                string
	ProductID         string // FK Product
	ReturnDate        *time.Time
	ReturnReason      string
	ConditionOnReturn string
	ReturnedBy        int // FK user
	Notes             string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type ReservationSetting struct {
	ID                      string
	DefaultReservationDays  int
	RequiresDeposit         bool
	DepositPercentage       float64
	AutoQualificationEnabled bool
	MinPurchasesForAuto     int
	MinTotalSpent           float64
	GoodReputationThreshold float64
	CreatedBy               int
	UpdatedAt               time.Time
}

type CustomerQualification struct {
	ID                      string
	CustomerID              string // FK Customer
	QualificationType       string // manual/auto/vip/premium
	IsActive                bool
	QualificationDate       *time.Time
	QualifiedBy             *int
	AutoQualificationReason string
	Notes                   string
	CreatedBy               int
	CreatedAt               time.Time
	UpdatedAt               time.Time
}

type CustomerReputation struct {
	ID                         string
	CustomerID                 string
	TotalPurchases             int
	TotalSpent                 float64
	AvgRating                  float64
	LatePaymentsCount          int
	CancelledReservationsCount int
	ReputationScore            float64
	LastCalculated             *time.Time
	CreatedAt                  time.Time
	UpdatedAt                  time.Time
}

type Reservation struct {
	ID                    string
	ProductID             string
	CustomerID            string
	ReservationDate       *time.Time
	ExpiryDate            *time.Time
	CustomReservationDays *int
	RequiresDeposit       bool
	DepositPercentage     float64
	DepositAmount         float64
	DepositPaid           bool
	DepositDate           *time.Time
	Status                string // active/expired/converted/cancelled
	CancellationReason    string
	Notes                 string
	CreatedBy             int
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

type Appraisal struct {
	ID                  string
	ProductID           string
	AppraiserName       string
	AppraisalDate       *time.Time
	EstimatedValue      float64
	ConditionAssessment string
	AuthenticityRating  float64
	Notes               string
	DocumentURL         string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type Certificate struct {
	ID                string
	ProductID         string
	CertificateType   string
	Issuer            string
	IssueDate         *time.Time
	ExpiryDate        *time.Time
	CertificateNumber string
	DocumentURL       string
	Verified          bool
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type Auction struct {
	ID          string
	Title       string
	Description string
	StartDate   *time.Time
	EndDate     *time.Time
	Status      string
	CreatedBy   int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type AuctionItem struct {
	ID           string
	AuctionID    string
	ProductID    string
	StartingBid  float64
	ReservePrice float64
	CurrentBid   float64
}

type Bid struct {
	ID            string
	AuctionItemID string
	BidderID      string // FK Customer
	BidAmount     float64
	BidDate       *time.Time
	Status        string // active/won/lost
}

type SearchResultReservation struct {
	Data       *[]Reservation
	Total      int64
	Page       int
	PageSize   int
	TotalPages int
}

type SearchResultAppraisal struct {
	Data       *[]Appraisal
	Total      int64
	Page       int
	PageSize   int
	TotalPages int
}

type SearchResultAuction struct {
	Data       *[]Auction
	Total      int64
	Page       int
	PageSize   int
	TotalPages int
}
