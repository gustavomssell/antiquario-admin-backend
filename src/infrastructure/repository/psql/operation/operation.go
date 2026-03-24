package operation

import (
	"time"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/catalog"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/commercial"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/product"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/user"
	domainOperation "github.com/gbrayhan/microservices-go/src/domain/operation"
)

type CommissionRule struct {
	ID                   string            `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID               *int              `gorm:"column:user_id"`
	User                 *user.User        `gorm:"foreignKey:UserID"`
	ProductCategoryID    *string           `gorm:"column:product_category_id;type:uuid"`
	ProductCategory      *catalog.Category `gorm:"foreignKey:ProductCategoryID"`
	CommissionPercentage float64           `gorm:"column:commission_percentage"`
	IsDefault            bool              `gorm:"column:is_default"`
	Active               bool              `gorm:"column:active;default:true"`
	CreatedBy            int               `gorm:"column:created_by"`
	CreatedAt            time.Time         `gorm:"autoCreateTime:mili"`
	UpdatedAt            time.Time         `gorm:"autoUpdateTime:mili"`
}
func (CommissionRule) TableName() string { return "commission_rules" }

type ConsignmentSetting struct {
	ID                  string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	DefaultPercentage   float64   `gorm:"column:default_percentage"`
	DefaultDeadlineDays int       `gorm:"column:default_deadline_days"`
	CreatedBy           int       `gorm:"column:created_by"`
	CreatedAt           time.Time `gorm:"autoCreateTime:mili"`
	UpdatedAt           time.Time `gorm:"autoUpdateTime:mili"`
}
func (ConsignmentSetting) TableName() string { return "consignment_settings" }

type ConsignmentReturn struct {
	ID                string           `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ProductID         string           `gorm:"column:product_id;type:uuid"`
	Product           *product.Product `gorm:"foreignKey:ProductID"`
	ReturnDate        *time.Time       `gorm:"column:return_date"`
	ReturnReason      string           `gorm:"column:return_reason"`
	ConditionOnReturn string           `gorm:"column:condition_on_return"`
	ReturnedBy        int              `gorm:"column:returned_by"`
	ReturnedByUser    *user.User       `gorm:"foreignKey:ReturnedBy"`
	Notes             string           `gorm:"column:notes"`
	CreatedAt         time.Time        `gorm:"autoCreateTime:mili"`
	UpdatedAt         time.Time        `gorm:"autoUpdateTime:mili"`
}
func (ConsignmentReturn) TableName() string { return "consignment_returns" }

type ReservationSetting struct {
	ID                       string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	DefaultReservationDays   int       `gorm:"column:default_reservation_days"`
	RequiresDeposit          bool      `gorm:"column:requires_deposit"`
	DepositPercentage        float64   `gorm:"column:deposit_percentage"`
	AutoQualificationEnabled bool      `gorm:"column:auto_qualification_enabled"`
	MinPurchasesForAuto      int       `gorm:"column:min_purchases_for_auto"`
	MinTotalSpent            float64   `gorm:"column:min_total_spent"`
	GoodReputationThreshold  float64   `gorm:"column:good_reputation_threshold"`
	CreatedBy                int       `gorm:"column:created_by"`
	UpdatedAt                time.Time `gorm:"autoUpdateTime:mili"`
}
func (ReservationSetting) TableName() string { return "reservation_settings" }

type CustomerQualification struct {
	ID                      string               `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CustomerID              string               `gorm:"column:customer_id;type:uuid"`
	Customer                *commercial.Customer `gorm:"foreignKey:CustomerID"`
	QualificationType       string               `gorm:"column:qualification_type"`
	IsActive                bool                 `gorm:"column:is_active;default:true"`
	QualificationDate       *time.Time           `gorm:"column:qualification_date"`
	QualifiedBy             *int                 `gorm:"column:qualified_by"`
	AutoQualificationReason string               `gorm:"column:auto_qualification_reason"`
	Notes                   string               `gorm:"column:notes"`
	CreatedBy               int                  `gorm:"column:created_by"`
	CreatedAt               time.Time            `gorm:"autoCreateTime:mili"`
	UpdatedAt               time.Time            `gorm:"autoUpdateTime:mili"`
}
func (CustomerQualification) TableName() string { return "customer_qualifications" }

type CustomerReputation struct {
	ID                         string               `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CustomerID                 string               `gorm:"column:customer_id;type:uuid;uniqueIndex"`
	Customer                   *commercial.Customer `gorm:"foreignKey:CustomerID"`
	TotalPurchases             int                  `gorm:"column:total_purchases"`
	TotalSpent                 float64              `gorm:"column:total_spent"`
	AvgRating                  float64              `gorm:"column:avg_rating"`
	LatePaymentsCount          int                  `gorm:"column:late_payments_count"`
	CancelledReservationsCount int                  `gorm:"column:cancelled_reservations_count"`
	ReputationScore            float64              `gorm:"column:reputation_score"`
	LastCalculated             *time.Time           `gorm:"column:last_calculated"`
	CreatedAt                  time.Time            `gorm:"autoCreateTime:mili"`
	UpdatedAt                  time.Time            `gorm:"autoUpdateTime:mili"`
}
func (CustomerReputation) TableName() string { return "customer_reputations" }

type Reservation struct {
	ID                    string               `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ProductID             string               `gorm:"column:product_id;type:uuid"`
	Product               *product.Product     `gorm:"foreignKey:ProductID"`
	CustomerID            string               `gorm:"column:customer_id;type:uuid"`
	Customer              *commercial.Customer `gorm:"foreignKey:CustomerID"`
	ReservationDate       *time.Time           `gorm:"column:reservation_date"`
	ExpiryDate            *time.Time           `gorm:"column:expiry_date"`
	CustomReservationDays *int                 `gorm:"column:custom_reservation_days"`
	RequiresDeposit       bool                 `gorm:"column:requires_deposit"`
	DepositPercentage     float64              `gorm:"column:deposit_percentage"`
	DepositAmount         float64              `gorm:"column:deposit_amount"`
	DepositPaid           bool                 `gorm:"column:deposit_paid"`
	DepositDate           *time.Time           `gorm:"column:deposit_date"`
	Status                string               `gorm:"column:status"`
	CancellationReason    string               `gorm:"column:cancellation_reason"`
	Notes                 string               `gorm:"column:notes"`
	CreatedBy             int                  `gorm:"column:created_by"`
	CreatedByUser         *user.User           `gorm:"foreignKey:CreatedBy"`
	CreatedAt             time.Time            `gorm:"autoCreateTime:mili"`
	UpdatedAt             time.Time            `gorm:"autoUpdateTime:mili"`
}
func (Reservation) TableName() string { return "reservations" }

type Appraisal struct {
	ID                  string           `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ProductID           string           `gorm:"column:product_id;type:uuid"`
	Product             *product.Product `gorm:"foreignKey:ProductID"`
	AppraiserName       string           `gorm:"column:appraiser_name"`
	AppraisalDate       *time.Time       `gorm:"column:appraisal_date"`
	EstimatedValue      float64          `gorm:"column:estimated_value"`
	ConditionAssessment string           `gorm:"column:condition_assessment"`
	AuthenticityRating  float64          `gorm:"column:authenticity_rating"`
	Notes               string           `gorm:"column:notes"`
	DocumentURL         string           `gorm:"column:document_url"`
	CreatedAt           time.Time        `gorm:"autoCreateTime:mili"`
	UpdatedAt           time.Time        `gorm:"autoUpdateTime:mili"`
}
func (Appraisal) TableName() string { return "appraisals" }

type Certificate struct {
	ID                string           `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ProductID         string           `gorm:"column:product_id;type:uuid"`
	Product           *product.Product `gorm:"foreignKey:ProductID"`
	CertificateType   string           `gorm:"column:certificate_type"`
	Issuer            string           `gorm:"column:issuer"`
	IssueDate         *time.Time       `gorm:"column:issue_date"`
	ExpiryDate        *time.Time       `gorm:"column:expiry_date"`
	CertificateNumber string           `gorm:"column:certificate_number"`
	DocumentURL       string           `gorm:"column:document_url"`
	Verified          bool             `gorm:"column:verified"`
	CreatedAt         time.Time        `gorm:"autoCreateTime:mili"`
	UpdatedAt         time.Time        `gorm:"autoUpdateTime:mili"`
}
func (Certificate) TableName() string { return "certificates" }

type Auction struct {
	ID          string     `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Title       string     `gorm:"column:title"`
	Description string     `gorm:"column:description"`
	StartDate   *time.Time `gorm:"column:start_date"`
	EndDate     *time.Time `gorm:"column:end_date"`
	Status      string     `gorm:"column:status"`
	CreatedBy   int        `gorm:"column:created_by"`
	CreatedAt   time.Time  `gorm:"autoCreateTime:mili"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime:mili"`
}
func (Auction) TableName() string { return "auctions" }

type AuctionItem struct {
	ID           string           `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	AuctionID    string           `gorm:"column:auction_id;type:uuid"`
	Auction      *Auction         `gorm:"foreignKey:AuctionID"`
	ProductID    string           `gorm:"column:product_id;type:uuid"`
	Product      *product.Product `gorm:"foreignKey:ProductID"`
	StartingBid  float64          `gorm:"column:starting_bid"`
	ReservePrice float64          `gorm:"column:reserve_price"`
	CurrentBid   float64          `gorm:"column:current_bid"`
}
func (AuctionItem) TableName() string { return "auction_items" }

type Bid struct {
	ID            string               `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	AuctionItemID string               `gorm:"column:auction_item_id;type:uuid"`
	AuctionItem   *AuctionItem         `gorm:"foreignKey:AuctionItemID"`
	BidderID      string               `gorm:"column:bidder_id;type:uuid"`
	Bidder        *commercial.Customer `gorm:"foreignKey:BidderID"`
	BidAmount     float64              `gorm:"column:bid_amount"`
	BidDate       *time.Time           `gorm:"column:bid_date"`
	Status        string               `gorm:"column:status"`
}
func (Bid) TableName() string { return "bids" }

func (r *Reservation) toDomainMapper() *domainOperation.Reservation {
	return &domainOperation.Reservation{
		ID:                    r.ID,
		ProductID:             r.ProductID,
		CustomerID:            r.CustomerID,
		ReservationDate:       r.ReservationDate,
		ExpiryDate:            r.ExpiryDate,
		CustomReservationDays: r.CustomReservationDays,
		RequiresDeposit:       r.RequiresDeposit,
		DepositPercentage:     r.DepositPercentage,
		DepositAmount:         r.DepositAmount,
		DepositPaid:           r.DepositPaid,
		DepositDate:           r.DepositDate,
		Status:                r.Status,
		CancellationReason:    r.CancellationReason,
		Notes:                 r.Notes,
		CreatedBy:             r.CreatedBy,
		CreatedAt:             r.CreatedAt,
		UpdatedAt:             r.UpdatedAt,
	}
}

func fromDomainReservationMapper(domain *domainOperation.Reservation) *Reservation {
	return &Reservation{
		ID:                    domain.ID,
		ProductID:             domain.ProductID,
		CustomerID:            domain.CustomerID,
		ReservationDate:       domain.ReservationDate,
		ExpiryDate:            domain.ExpiryDate,
		CustomReservationDays: domain.CustomReservationDays,
		RequiresDeposit:       domain.RequiresDeposit,
		DepositPercentage:     domain.DepositPercentage,
		DepositAmount:         domain.DepositAmount,
		DepositPaid:           domain.DepositPaid,
		DepositDate:           domain.DepositDate,
		Status:                domain.Status,
		CancellationReason:    domain.CancellationReason,
		Notes:                 domain.Notes,
		CreatedBy:             domain.CreatedBy,
	}
}

func (a *Appraisal) toDomainMapper() *domainOperation.Appraisal {
	return &domainOperation.Appraisal{
		ID:                  a.ID,
		ProductID:           a.ProductID,
		AppraiserName:       a.AppraiserName,
		AppraisalDate:       a.AppraisalDate,
		EstimatedValue:      a.EstimatedValue,
		ConditionAssessment: a.ConditionAssessment,
		AuthenticityRating:  a.AuthenticityRating,
		Notes:               a.Notes,
		DocumentURL:         a.DocumentURL,
		CreatedAt:           a.CreatedAt,
		UpdatedAt:           a.UpdatedAt,
	}
}

func fromDomainAppraisalMapper(d *domainOperation.Appraisal) *Appraisal {
	return &Appraisal{
		ID:                  d.ID,
		ProductID:           d.ProductID,
		AppraiserName:       d.AppraiserName,
		AppraisalDate:       d.AppraisalDate,
		EstimatedValue:      d.EstimatedValue,
		ConditionAssessment: d.ConditionAssessment,
		AuthenticityRating:  d.AuthenticityRating,
		Notes:               d.Notes,
		DocumentURL:         d.DocumentURL,
	}
}

func (auc *Auction) toDomainMapper() *domainOperation.Auction {
	return &domainOperation.Auction{
		ID:          auc.ID,
		Title:       auc.Title,
		Description: auc.Description,
		StartDate:   auc.StartDate,
		EndDate:     auc.EndDate,
		Status:      auc.Status,
		CreatedBy:   auc.CreatedBy,
		CreatedAt:   auc.CreatedAt,
		UpdatedAt:   auc.UpdatedAt,
	}
}

func fromDomainAuctionMapper(d *domainOperation.Auction) *Auction {
	return &Auction{
		ID:          d.ID,
		Title:       d.Title,
		Description: d.Description,
		StartDate:   d.StartDate,
		EndDate:     d.EndDate,
		Status:      d.Status,
		CreatedBy:   d.CreatedBy,
	}
}
